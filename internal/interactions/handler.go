package interactions

import (
	"errors"
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

func (h *Handler) ConfirmIncidentHandler(c *gin.Context) {
	var req ConfirmIncidentRequest

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	claims, err := auth.GetClaimsFromContext(c.Request.Context())

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = h.service.ConfirmIncidentService(c.Request.Context(), req, claims)

	if err != nil {
		status := http.StatusInternalServerError
		switch {
		case errors.Is(err, ErrIncidentNotFound):
			status = http.StatusNotFound
		case errors.Is(err, ErrAlreadyConfirmed):
			status = http.StatusConflict
		}

		c.JSON(status, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, ConfirmIncidentResponse{
		Message: "Incident confirmed",
	})

}

func (h *Handler) PostFeedbackHandler(c *gin.Context) {
	var uriReq PostFeedbackURIRequest

	if err := c.ShouldBindUri(&uriReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var bodyReq PostFeedbackBodyRequest

	if err := c.ShouldBindJSON(&bodyReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	req := PostFeedbackRequest{
		PostId:   uriReq.PostId,
		Feedback: bodyReq.Feedback,
	}

	claims, err := auth.GetClaimsFromContext(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = h.service.PostFeedbackService(c.Request.Context(), req, claims)
	if err != nil {
		status := http.StatusInternalServerError
		switch {
		case errors.Is(err, ErrOwnPostFeedback):
			status = http.StatusBadRequest
		case errors.Is(err, ErrPostNotFound):
			status = http.StatusNotFound
		case errors.Is(err, ErrAlreadyGaveFeedback):
			status = http.StatusConflict
		}

		c.JSON(status, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, PostFeedbackResponse{
		Message: "Feedback submitted",
	})
}

func (h *Handler) ReportPost(c *gin.Context) {
	var req PostReportRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	claims, err := auth.GetClaimsFromContext(c.Request.Context())

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	postReport, err := h.service.PostReportService(c.Request.Context(), req, claims)

	if err != nil {
		status := http.StatusInternalServerError
		switch {
		case errors.Is(err, ErrOwnPostReport):
			status = http.StatusBadRequest
		case errors.Is(err, ErrPostNotFound):
			status = http.StatusNotFound
		}

		c.JSON(status, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, PostReportResponse{
		Report: *postReport,
	})
}
