package validator

import (
	"testing"
)

func TestValidatePasswordStrength(t *testing.T) {
	tests := []struct {
		name    string
		password string
		wantErr bool
	}{
		{
			name:     "有效密码",
			password: "Test123!@#",
			wantErr:  false,
		},
		{
			name:     "密码太短",
			password: "Test1!",
			wantErr:  true,
		},
		{
			name:     "密码太长",
			password: "Test123!@#Test123!@#Test123!@#Test123!@#Test123!@#",
			wantErr:  true,
		},
		{
			name:     "缺少数字",
			password: "TestPass!@#",
			wantErr:  true,
		},
		{
			name:     "缺少大写字母",
			password: "test123!@#",
			wantErr:  true,
		},
		{
			name:     "缺少小写字母",
			password: "TEST123!@#",
			wantErr:  true,
		},
		{
			name:     "缺少特殊字符",
			password: "Test1234",
			wantErr:  true,
		},
		{
			name:     "边界值-8字符",
			password: "Test1!@#",
			wantErr:  false,
		},
		{
			name:     "边界值-32字符",
			password: "Test123!@#Test123!@#Test123!@#",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePasswordStrength(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePasswordStrength() error = %v, wantErr %v", err, tt.wantErr)
				if err != nil {
					t.Logf("错误信息: %v", err.Error())
				}
			}
		})
	}
}

