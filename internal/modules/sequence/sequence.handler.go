package sequence

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kapilkumar9395/pypr/internal/common"
	"github.com/kapilkumar9395/pypr/internal/modules/auth"
)

type SequenceHandler struct {
	service *SequenceService
}

func NewSequenceHandler(service *SequenceService) *SequenceHandler {
	return &SequenceHandler{
		service: service,
	}
}

func (h *SequenceHandler) RegisterRoutes(router *gin.RouterGroup) {
	seq := router.Group("/sequence")

	seq.POST("/create", h.HandleCreateSequence)
}

func (h *SequenceHandler) HandleCreateSequence(c *gin.Context) {
	var req CreateSequenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()

	organizationId := common.GetContextValue(c, auth.ContextOrganizationId)
	userId := common.GetContextValue(c, auth.ContextUserId)

	err := h.service.CreateSequence(ctx, organizationId, userId, req)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, common.GenericResponse{
		Message: "Sequence created successfully",
		Data:    nil,
	})
}
