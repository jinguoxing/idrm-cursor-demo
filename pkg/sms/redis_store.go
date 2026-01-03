package sms

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/redis/go-redis/v9"
)

// redisStore Redis存储实现
type redisStore struct {
	client *redis.Client
}

// NewRedisStore 创建SMS服务（Redis实现）
func NewRedisStore(client *redis.Client) Service {
	return &redisStore{
		client: client,
	}
}

// generateCode 生成6位数字验证码
func generateCode() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

// getCodeKey 获取验证码存储Key
func getCodeKey(mobile, codeType string) string {
	return fmt.Sprintf("idrm:sms:code:%s:%s", codeType, mobile)
}

// getRateLimitKey 获取频率限制Key
func getRateLimitKey(mobile string) string {
	return fmt.Sprintf("idrm:sms:ratelimit:%s", mobile)
}

// CheckRateLimit 检查发送频率限制（60秒内只能发送一次）
func (s *redisStore) CheckRateLimit(ctx context.Context, mobile string) error {
	key := getRateLimitKey(mobile)

	// 检查是否在60秒内发送过
	exists, err := s.client.Exists(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("检查频率限制失败: %w", err)
	}

	if exists > 0 {
		ttl, err := s.client.TTL(ctx, key).Result()
		if err != nil {
			return fmt.Errorf("获取剩余时间失败: %w", err)
		}
		return fmt.Errorf("发送过于频繁，请%d秒后重试", int(ttl.Seconds())+1)
	}

	return nil
}

// SendCode 发送验证码
func (s *redisStore) SendCode(ctx context.Context, mobile string, codeType string) error {
	// 检查频率限制
	if err := s.CheckRateLimit(ctx, mobile); err != nil {
		return err
	}

	// 生成6位验证码
	code := generateCode()

	// 存储验证码（5分钟有效期）
	codeKey := getCodeKey(mobile, codeType)
	if err := s.client.Set(ctx, codeKey, code, 5*time.Minute).Err(); err != nil {
		return fmt.Errorf("存储验证码失败: %w", err)
	}

	// 设置频率限制（60秒）
	rateLimitKey := getRateLimitKey(mobile)
	if err := s.client.Set(ctx, rateLimitKey, "1", 60*time.Second).Err(); err != nil {
		return fmt.Errorf("设置频率限制失败: %w", err)
	}

	// TODO: 调用短信服务API发送验证码
	// 这里暂时只存储，实际发送需要集成短信服务提供商

	return nil
}

// VerifyCode 验证验证码
func (s *redisStore) VerifyCode(ctx context.Context, mobile string, codeType string, code string) error {
	key := getCodeKey(mobile, codeType)

	// 获取存储的验证码
	storedCode, err := s.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return errors.New("验证码不存在或已过期")
	}
	if err != nil {
		return fmt.Errorf("获取验证码失败: %w", err)
	}

	// 验证验证码
	if storedCode != code {
		return errors.New("验证码错误")
	}

	// 验证成功后删除验证码（一次性使用）
	if err := s.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("删除验证码失败: %w", err)
	}

	return nil
}
