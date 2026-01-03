package device

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// ParseDeviceInfo 解析设备信息
// 从User-Agent中提取设备类型和设备标识
func ParseDeviceInfo(userAgent string) (deviceType, deviceID string) {
	// 解析设备类型
	deviceType = "Web" // 默认Web

	ua := strings.ToLower(userAgent)
	if strings.Contains(ua, "iphone") || strings.Contains(ua, "ipad") {
		deviceType = "iOS"
	} else if strings.Contains(ua, "android") {
		deviceType = "Android"
	} else if strings.Contains(ua, "mobile") {
		deviceType = "Mobile"
	}

	// 生成设备ID（基于User-Agent的哈希）
	h := sha256.Sum256([]byte(userAgent))
	deviceID = hex.EncodeToString(h[:])[:16] // 取前16个字符

	return deviceType, deviceID
}
