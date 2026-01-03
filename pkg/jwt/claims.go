package jwt

// Claims JWT Token Claims结构体
type Claims struct {
	UserID string `json:"user_id"` // 用户ID (UUID v7)
}
