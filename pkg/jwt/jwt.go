package jwt

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

// Service JWT服务接口
type Service interface {
	// GenerateToken 生成JWT Token
	GenerateToken(ctx context.Context, userID string) (string, error)

	// VerifyToken 验证JWT Token
	VerifyToken(ctx context.Context, tokenString string) (*Claims, error)

	// InvalidateUserTokens 使指定用户的所有Token失效（密码重置后）
	InvalidateUserTokens(ctx context.Context, userID string) error
}

// jwtService JWT服务实现
type jwtService struct {
	secretKey []byte
	redis     *redis.Client
}

// NewJWTService 创建JWT服务
func NewJWTService(secretKey string, redisClient *redis.Client) Service {
	return &jwtService{
		secretKey: []byte(secretKey),
		redis:     redisClient,
	}
}

// GenerateToken 生成JWT Token（7天有效期）
func (s *jwtService) GenerateToken(ctx context.Context, userID string) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     now.Add(7 * 24 * time.Hour).Unix(),
		"iat":     now.Unix(),
		"iss":     "idrm-api",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", fmt.Errorf("生成Token失败: %w", err)
	}

	return tokenString, nil
}

// VerifyToken 验证JWT Token
func (s *jwtService) VerifyToken(ctx context.Context, tokenString string) (*Claims, error) {
	// 解析Token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("意外的签名方法: %v", token.Header["alg"])
		}
		return s.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("Token解析失败: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("Token无效")
	}

	// 提取Claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("Token Claims格式错误")
	}

	// 检查Token是否已失效
	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, errors.New("Token中缺少user_id")
	}

	// 检查Redis中的失效列表
	if s.isTokenInvalid(ctx, userID) {
		return nil, errors.New("Token已失效，请重新登录")
	}

	return &Claims{
		UserID: userID,
	}, nil
}

// InvalidateUserTokens 使指定用户的所有Token失效
func (s *jwtService) InvalidateUserTokens(ctx context.Context, userID string) error {
	key := fmt.Sprintf("idrm:jwt:invalid:%s", userID)
	// 设置标记，有效期7天（与Token有效期一致）
	return s.redis.Set(ctx, key, "1", 7*24*time.Hour).Err()
}

// isTokenInvalid 检查Token是否已失效
func (s *jwtService) isTokenInvalid(ctx context.Context, userID string) bool {
	key := fmt.Sprintf("idrm:jwt:invalid:%s", userID)
	exists, err := s.redis.Exists(ctx, key).Result()
	if err != nil {
		return false // 如果Redis查询失败，不阻止验证
	}
	return exists > 0
}
