package device

import (
	"testing"
)

func TestParseDeviceInfo(t *testing.T) {
	tests := []struct {
		name      string
		userAgent string
		wantType  string
		wantIDLen int
	}{
		{
			name:      "Web浏览器",
			userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			wantType:  "Web",
			wantIDLen: 16,
		},
		{
			name:      "Android设备",
			userAgent: "Mozilla/5.0 (Linux; Android 10; Mobile) AppleWebKit/537.36",
			wantType:  "Android",
			wantIDLen: 16,
		},
		{
			name:      "iOS设备",
			userAgent: "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/605.1.15",
			wantType:  "iOS",
			wantIDLen: 16,
		},
		{
			name:      "iPad设备",
			userAgent: "Mozilla/5.0 (iPad; CPU OS 14_0 like Mac OS X) AppleWebKit/605.1.15",
			wantType:  "iOS",
			wantIDLen: 16,
		},
		{
			name:      "空User-Agent",
			userAgent: "",
			wantType:  "Web",
			wantIDLen: 16,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotType, gotID := ParseDeviceInfo(tt.userAgent)
			if gotType != tt.wantType {
				t.Errorf("ParseDeviceInfo() deviceType = %v, want %v", gotType, tt.wantType)
			}
			if len(gotID) != tt.wantIDLen {
				t.Errorf("ParseDeviceInfo() deviceID长度 = %v, want %v", len(gotID), tt.wantIDLen)
			}
		})
	}
}

func TestParseDeviceInfo_Consistent(t *testing.T) {
	// 测试相同User-Agent生成相同的设备ID
	ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
	id1, _ := ParseDeviceInfo(ua)
	id2, _ := ParseDeviceInfo(ua)
	
	if id1 != id2 {
		t.Errorf("ParseDeviceInfo() 相同User-Agent应生成相同设备ID, id1=%v, id2=%v", id1, id2)
	}
}

