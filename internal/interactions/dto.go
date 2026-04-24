package interactions

import "locallyn-be/internal/common/constants"

type ConfirmIncidentRequest struct {
	IncidentId string `uri:"incidentId" binding:"required,uuid"`
}

type ConfirmIncidentResponse struct {
	Message string `json:"message"`
}

type PostFeedbackRequest struct {
	PostId   string                 `uri:"postId" binding:"required,uuid"`
	Feedback constants.FeedbackType `json:"feedback" binding:"required,oneof=HELPFUL MISLEADING"`
}

type PostFeedbackURIRequest struct {
	PostId string `uri:"postId" binding:"required,uuid"`
}

type PostFeedbackBodyRequest struct {
	Feedback constants.FeedbackType `json:"feedback" binding:"required,oneof=HELPFUL MISLEADING"`
}

type PostFeedbackResponse struct {
	Message string `json:"message"`
}
