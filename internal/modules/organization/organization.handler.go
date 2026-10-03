package organization

import (
	"github.com/gin-gonic/gin"
	"github.com/kapilkumar9395/pypr/internal/common"
	"github.com/kapilkumar9395/pypr/internal/modules/auth"
)

type OrganizationHandler struct {
	OrganizationService *OrganizationService
}

func NewOrganizationHandler(organizationService *OrganizationService) *OrganizationHandler {
	return &OrganizationHandler{
		OrganizationService: organizationService,
	}
}

func (h *OrganizationHandler) RegisterRoutes(router *gin.RouterGroup) {
	organization := router.Group("/org")

	organization.POST("/invite", h.InviteUser)
}

func (h *OrganizationHandler) InviteUser(c *gin.Context) {
	var req InviteUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()

	organizationId := common.GetContextValue(c, auth.ContextOrganizationId)
	userId := common.GetContextValue(c, auth.ContextUserId)

	message, err := h.OrganizationService.InviteUser(ctx, organizationId, userId, req.EmailId)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": message})
}
