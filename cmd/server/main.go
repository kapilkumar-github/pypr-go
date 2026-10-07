package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"uuid"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/kapilkumar9395/pypr/internal/database"
	"github.com/kapilkumar9395/pypr/internal/infra/cache"
	emailInfra "github.com/kapilkumar9395/pypr/internal/infra/email"
	"github.com/kapilkumar9395/pypr/internal/modules/auth"
	"github.com/kapilkumar9395/pypr/internal/modules/contact"
	"github.com/kapilkumar9395/pypr/internal/modules/organization"
	"github.com/kapilkumar9395/pypr/internal/modules/sequence"
	"github.com/kapilkumar9395/pypr/internal/modules/variable"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Warning: .env file not found")
	}

	ctx := context.Background()

	app, err := NewApp(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer app.DB.Close()

	server := NewServer(app)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
	log.Printf("Server is running on %s", server.HTTP.Addr)
}

type App struct {
	DB                *pgxpool.Pool
	AuthRepo          *auth.AuthRepository
	OrganizationRepo  *organization.OrganizationRepository
	ContactRepository *contact.ContactRepository
	SequenceRepo      *sequence.SequenceRepo
	VariableRepo      *variable.VariableRepo

	// Infra
	EmailFactory *emailInfra.Factory

	JwtService *auth.JwtService
}

type Server struct {
	Router *gin.Engine
	HTTP   *http.Server
}

func (s *Server) Run() error {
	return s.HTTP.ListenAndServe()
}

func NewApp(ctx context.Context) (*App, error) {
	db, err := database.NewPostgresPool(ctx)
	if err != nil {
		return nil, err
	}

	// Infra - Email
	resendEmailSender := emailInfra.NewResendSender(os.Getenv("RESEND_API_KEY"))
	// Cache
	cache := cache.New[uuid.UUID](5 * time.Minute)

	jwtService := auth.NewJwtService(os.Getenv("JWT_SECRET_KEY"))

	return &App{
		DB:                db,
		AuthRepo:          auth.NewAuthRepository(db),
		OrganizationRepo:  organization.NewOrganizationRepository(db),
		ContactRepository: contact.NewContactRepository(db),
		SequenceRepo:      sequence.NewSequenceRepo(db),
		VariableRepo:      variable.NewVariableRepo(db, cache),

		// Infra
		EmailFactory: emailInfra.NewFactory(resendEmailSender),
		JwtService:   &jwtService,
	}, nil
}

func NewServer(app *App) *Server {
	router := gin.New()
	router.RedirectTrailingSlash = false
	corsOrigins := strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ",")
	for i := range corsOrigins {
		corsOrigins[i] = strings.TrimSpace(corsOrigins[i])
	}
	router.Use(cors.New(cors.Config{
		AllowOrigins:     corsOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	api := router.Group("/api")

	//middleware
	authMiddleware := auth.NewAuthMiddleware(
		os.Getenv("JWT_SECRET_KEY"),
	)

	protected := api.Group("")
	protected.Use(authMiddleware.Authenticate())

	authService := auth.NewAuthService(app.AuthRepo, app.JwtService)
	authHandler := auth.NewAuthHandler(authService)
	authHandler.RegisterRoutes(api)

	organizationService := organization.NewOrganizationService(*app.OrganizationRepo, *app.EmailFactory)
	organizationHandler := organization.NewOrganizationHandler(organizationService)
	organizationHandler.RegisterRoutes(protected)

	contactService := contact.NewContactService(app.ContactRepository)
	contactHandler := contact.NewContactHandler(contactService)
	contactHandler.RegisterRoutes(protected)

	sequenceService := sequence.NewSequenceService(app.SequenceRepo, app.VariableRepo)
	sequenceHandler := sequence.NewSequenceHandler(sequenceService)
	sequenceHandler.RegisterRoutes(protected)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	return &Server{
		Router: router,
		HTTP: &http.Server{
			Addr:    "127.0.0.1:" + port,
			Handler: router,
		},
	}
}
