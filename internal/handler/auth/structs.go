package auth

type signUpRequest struct {
	Username string `json:"username" binding:"required,min=5,max=13"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=15"`
}

type signInRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type logoutRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type refreshRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type activateRequest struct {
	ActivationCode string `json:"activationCode" binding:"required"`
}
