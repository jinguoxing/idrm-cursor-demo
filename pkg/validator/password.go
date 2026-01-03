package validator

import (
	"errors"
	"regexp"
)

// ValidatePasswordStrength 验证密码强度
// 要求：
// - 长度：8-32个字符
// - 必须包含：数字、大写字母、小写字母、特殊字符
func ValidatePasswordStrength(password string) error {
	// 长度检查：8-32字符
	if len(password) < 8 || len(password) > 32 {
		return errors.New("密码长度必须在8-32个字符之间")
	}

	// 复杂度检查：必须包含数字、大写字母、小写字母、特殊字符
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password)

	if !hasDigit || !hasUpper || !hasLower || !hasSpecial {
		return errors.New("密码必须包含数字、大写字母、小写字母和特殊字符")
	}

	return nil
}
