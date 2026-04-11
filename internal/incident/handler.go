package incident

import (
	"locallyn-be/internal/common/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetNearbyIncidentHandler(c *gin.Context) {
	var req GetNearbyIncidentRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	_, err := auth.GetClaimsFromContext(c.Request.Context())

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	incidents, err := h.service.GetNearbyIncidentService(c.Request.Context(), req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, GetNearbyIncidentResponse{
		NearbyIncidents: incidents,
	})

}

func (h *Handler) GetIncidentHandler(c *gin.Context) {
	var req GetIncidentRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	_, err := auth.GetClaimsFromContext(c.Request.Context())

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
	}

	incident, err := h.service.GetIncidentService(c.Request.Context(), req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(200, GetIncidentResponse{
		Incident: *incident,
	})

}
