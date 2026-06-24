package dto

// Request DTOs

type SendCodeRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type RegisterRequest struct {
	FirstName   string  `json:"firstName" binding:"required"`
	LastName    *string `json:"lastName"`
	PhoneNumber *string `json:"phoneNumber"`
	Email       string  `json:"email" binding:"required,email"`
	Code        string  `json:"code" binding:"required"`
}

type LoginRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required"`
}

type UpdateStatusRequest struct {
	UserID uint64 `json:"userId" binding:"required"`
	Status int8   `json:"status" binding:"min=0,max=4"`
}

// Response DTOs

type MessageResponse struct {
	Message string `json:"message"`
}

type UserResponse struct {
	ID          uint64  `json:"id"`
	FirstName   string  `json:"firstName"`
	LastName    *string `json:"lastName"`
	PhoneNumber *string `json:"phoneNumber,omitempty"`
	Email       string  `json:"email"`
	Status      int8    `json:"status"`
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type UserListItem struct {
	FirstName string  `json:"firstName"`
	LastName  *string `json:"lastName"`
	Email     string  `json:"email"`
	Status    int8    `json:"status"`
}

type StatisticsResponse struct {
	New        float64 `json:"new"`
	Reviewed   float64 `json:"reviewed"`
	Interviewed float64 `json:"interviewed"`
	OfferSent  float64 `json:"offerSent"`
	Rejected   float64 `json:"rejected"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
