package apis

type PowerTableDoc struct {
	ID        uint    `json:"id" example:"1"`
	CreatedAt string  `json:"createdAt" example:"2026-04-23T10:00:00Z"`
	UpdatedAt string  `json:"updatedAt" example:"2026-04-23T10:00:00Z"`
	Price     float64 `json:"price" example:"99.99"`
	Company   string  `json:"company" example:"Acme Corp"`
	UserID    string  `json:"userId" example:"user-abc-123"`
}

type RevokedTokenDoc struct {
	ID        uint   `json:"id" example:"1"`
	JTI       string `json:"jti" example:"550e8400-e29b-41d4-a716-446655440000"`
	RevokedAt string `json:"revokedAt" example:"2026-04-23T10:00:00Z"`
	ExpiresAt string `json:"expiresAt" example:"2026-04-24T10:00:00Z"`
}
