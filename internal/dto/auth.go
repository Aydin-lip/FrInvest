package dto

// Request DTOs

type SendVerificationRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type RegisterRequest struct {
	FirstName   string  `json:"firstName" binding:"required"`
	LastName    *string `json:"lastName"`
	PhoneNumber *string `json:"phoneNumber"`
	Email       string  `json:"email" binding:"required,email"`
}

type UpdateStatusRequest struct {
	UserID uint64 `json:"userId" binding:"required"`
	Status int8   `json:"status" binding:"min=0,max=4"`
}

// Response DTOs

type MessageResponse struct {
	Message string `json:"message"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type RateLimitResponse struct {
	Message          string `json:"message"`
	RemainingSeconds int    `json:"remainingSeconds"`
}

type UserListItem struct {
	FirstName string  `json:"firstName"`
	LastName  *string `json:"lastName"`
	Email     string  `json:"email"`
	Status    int8    `json:"status"`
}

type StatusPercentagesResponse struct {
	New         float64 `json:"New"`
	Reviewed    float64 `json:"Reviewed"`
	Interviewed float64 `json:"Interviewed"`
	OfferSent   float64 `json:"OfferSent"`
	Rejected    float64 `json:"Rejected"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}
