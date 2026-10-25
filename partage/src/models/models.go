package models

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

type LoginRequest struct {
	Role     string `json:"role"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token     string `json:"token"`
	Role      string `json:"role"`
	ExpiresAt string `json:"expiresAt"`
}

type SetupRequest struct {
	AdminPassword string `json:"adminPassword"`
	UserPassword  string `json:"userPassword"`
}

type ChangePasswordRequest struct {
	CurrentAdminPassword string `json:"currentAdminPassword"`
	AdminPassword        string `json:"adminPassword"`
	UserPassword         string `json:"userPassword"`
}

type ServerStatus struct {
	Status    string `json:"status"`
	Version   string `json:"version"`
	Database  string `json:"database"`
	Timestamp string `json:"timestamp"`
}
