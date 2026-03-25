package authdto

type SessionUserResponse struct {
	ID           string  `json:"id"`
	Email        *string `json:"email,omitempty"`
	DisplayName  string  `json:"displayName"`
	Organization *string `json:"organization,omitempty"`
	Occupation   *string `json:"occupation,omitempty"`
	AvatarURL    *string `json:"avatarUrl,omitempty"`
}

type SessionResponse struct {
	Authenticated bool                 `json:"authenticated"`
	User          *SessionUserResponse `json:"user"`
}

type ForgotPasswordResponse struct {
	Submitted bool `json:"submitted"`
}

type ResetPasswordResponse struct {
	Reset bool `json:"reset"`
}

type ResetPasswordValidationResponse struct {
	Valid  bool   `json:"valid"`
	Reason string `json:"reason,omitempty"`
}

type LogoutResponse struct {
	LoggedOut bool `json:"loggedOut"`
}
