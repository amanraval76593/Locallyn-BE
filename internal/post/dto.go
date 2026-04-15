package post

import "locallyn-be/internal/common/constants"

type CreatePostRequest struct {
	Content   string                 `json:"content" binding:"required,min=4"`
	Latitude  float64                `json:"lat"  binding:"required,min=-90,max=90"`
	Longitude float64                `json:"long" binding:"required,min=-180,max=180"`
	Type      constants.PostType     `json:"type" binding:"required,oneof=INCIDENT BROADCAST"`
	Category  string                 `json:"category" binding:"required,min=4"`
	Identity  constants.IdentityType `json:"identity" binding:"required,oneof=PUBLIC PSEUDONYMOUS ANONYMOUS"`
	MediaURLs []string               `json:"media_urls" binding:"required"`
}

type CreatePostResponse struct {
	Post Post `json:"post"`
}

type FetchPostByIdRequest struct {
	PostId string `form:"postId" binding:"required,uuid"`
}

type FetchPostByIdResponse struct {
	Post Post `json:"post"`
}
