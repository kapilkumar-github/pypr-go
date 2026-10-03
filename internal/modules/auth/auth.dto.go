package auth

type CreateUserParams struct {
	Email        string `json:"email" binding:"required,email"`
	FirstName    string `json:"firstName" binding:"required"`
	LastName     string `json:"lastName"`
	PasswordHash string `json:"passwordHash" binding:"required"`
	Timezone     string `json:"timezone"`
}

type RegisterUserRequest struct {
	EmailId         string  `json:"emailId" binding:"required,email"`
	FirstName       string  `json:"firstName" binding:"required"`
	LastName        string  `json:"lastName"`
	Password        string  `json:"password" binding:"required"`
	InvitationToken *string `json:"invitationToken"`
	Timezone        *string `json:"timezone"`
}

type LoginRequest struct {
	EmailId  string `json:"emailId"`
	Password string `json:"password"`
}

type OrganizationByTokenProjection struct {
	OrganizationId string `json:"organizationId" db:"organization_id"`
	Role           string `json:"role" db:"role"`
	Name           string `json:"name" db:"name"`
}

type UserLoginProjection struct {
	UserId         string `json:"userId" db:"id"`
	Email          string `json:"email" db:"email"`
	FirstName      string `json:"firstName" db:"first_name"`
	LastName       string `json:"lastName" db:"last_name"`
	Timezone       string `json:"timezone" db:"timezone"`
	PasswordHash   string `json:"passwordHash" db:"password_hash"`
	OrganizationId string `json:"organizationId" db:"organization_id"`
}

type LoginResponse struct {
	AccessToken string `json:"accessToken"`
}
