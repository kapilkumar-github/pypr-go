package contact

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kapilkumar9395/pypr/internal/common"
	"github.com/kapilkumar9395/pypr/internal/modules/auth"
)

type ContactHandler struct {
	service *ContactService
}

func NewContactHandler(service *ContactService) *ContactHandler {
	return &ContactHandler{service}
}

func (h *ContactHandler) RegisterRoutes(router *gin.RouterGroup) {
	contacts := router.Group("/contacts")

	contacts.GET("/all", h.HandleGetContacts)
	contacts.POST("/create", h.HandleCreateContact)
}

func (h *ContactHandler) HandleCreateContact(c *gin.Context) {
	var req CreateContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	organizationId := common.GetContextValue(c, auth.ContextOrganizationId)
	userId := common.GetContextValue(c, auth.ContextUserId)

	createContactRes := h.service.CreateContact(c.Request.Context(), req, organizationId, userId)
	if createContactRes.Error != nil {
		log.Printf("%s", createContactRes.Error)
		createContactRes.Error = nil
		c.JSON(400, createContactRes)
		return
	}

	c.JSON(http.StatusCreated, createContactRes)
}

func (h *ContactHandler) HandleGetContacts(c *gin.Context) {
	organizationId := common.GetContextValue(c, auth.ContextOrganizationId)
	contactsRes := h.service.GetOrgAllContacts(c.Request.Context(), organizationId)
	if contactsRes.Error != nil {
		log.Printf("%s", contactsRes.Error)
		contactsRes.Error = nil
		c.JSON(400, contactsRes)
		return
	}

	c.JSON(http.StatusOK, contactsRes)
}
