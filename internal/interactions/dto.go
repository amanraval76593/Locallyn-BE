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

type PostReportRequest struct {
	PostId string `json:"postId" binding:"required,uuid"`
	Reason string `json:"reason" binding:"required,min=4"`
}

type PostReportResponse struct {
	Report PostReport `json:"postReport"`
}
