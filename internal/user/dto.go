package user

type SignUpRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=1"`
}

type SignUpResponse struct {
	Message string `json:"message"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"passsword" binding:"required,min=1"`
}

type LoginResponse struct {
	AccessToken string `json:"accessToken"`
}
