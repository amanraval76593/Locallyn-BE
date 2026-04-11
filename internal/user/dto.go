package user

type SignUpRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=1"`
}

type SignUpResponse struct {
	Message string `json:"message"`
}

type VerifyUserRequest struct {
	Token string `form:"token" binding:"required"`
}

type verifyUserResponse struct {
	Message string `json:"message"`
}
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=1"`
}

type LoginResponse struct {
	AccessToken string `json:"accessToken"`
}

type CreateProfileRequest struct {
	Username    string `json:"userName" binding:"required,min=4"`
	DisplayName string `json:"displayName" binding:"required,min=4"`
}

type CreateProfileResponse struct {
	Message string `json:"message"`
}

type GetProfileRequest struct {
	UserName string `form:"userName" binding:"required,min=4"`
}

type GetProfileResponse struct {
	UserProfileData UserProfile `json:"userProfileData"`
}

type UpdateProfileRequest struct {
	Username    *string `json:"userName,omitempty"`
	DisplayName *string `json:"displayName,omitempty"`
	AvatarURL   *string `json:"avatarUrl,omitempty"`
}

type UpdateProfileResponse struct {
	UserProfileData UserProfile `json:"userProfileData"`
}
