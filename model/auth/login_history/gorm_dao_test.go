package login_history

import (
	"context"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB 创建测试数据库（每个测试使用独立的内存数据库）
func setupTestDB(t *testing.T) *gorm.DB {
	// 使用 :memory: 而不是 file::memory:?cache=shared，确保每个测试都有独立的数据库
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("创建测试数据库失败: %v", err)
	}

	// 自动迁移
	if err := db.AutoMigrate(&LoginHistory{}); err != nil {
		t.Fatalf("数据库迁移失败: %v", err)
	}

	return db
}

func TestGormDao_Insert(t *testing.T) {
	db := setupTestDB(t)
	model := newGormDao(db)
	ctx := context.Background()

	history := &LoginHistory{
		Id:        "01234567-89ab-7def-0123-456789abc000",
		UserID:    "01234567-89ab-7def-0123-456789abcde0",
		IP:        "127.0.0.1",
		DeviceType: "Web",
		DeviceID:  "device123",
		UserAgent: "Mozilla/5.0",
		LoginAt:   time.Now(),
	}

	got, err := model.Insert(ctx, history)
	if err != nil {
		t.Errorf("Insert() error = %v", err)
		return
	}
	if got == nil {
		t.Errorf("Insert() 返回结果为空")
	}
	if got.Id != history.Id {
		t.Errorf("Insert() Id = %v, want %v", got.Id, history.Id)
	}
}

func TestGormDao_FindOne(t *testing.T) {
	db := setupTestDB(t)
	model := newGormDao(db)
	ctx := context.Background()

	// 先创建一条记录
	historyID := "01234567-89ab-7def-0123-456789abc001"
	history := &LoginHistory{
		Id:        historyID,
		UserID:    "01234567-89ab-7def-0123-456789abcde0",
		IP:        "127.0.0.1",
		DeviceType: "Web",
		DeviceID:  "device123",
		UserAgent: "Mozilla/5.0",
		LoginAt:   time.Now(),
	}
	_, err := model.Insert(ctx, history)
	if err != nil {
		t.Fatalf("准备测试数据失败: %v", err)
	}

	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{
			name:    "查询存在的记录",
			id:      historyID,
			wantErr: false,
		},
		{
			name:    "查询不存在的记录",
			id:      "00000000-0000-0000-0000-000000000000",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := model.FindOne(ctx, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindOne() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("FindOne() 返回结果为空")
			}
			if !tt.wantErr && got.Id != tt.id {
				t.Errorf("FindOne() Id = %v, want %v", got.Id, tt.id)
			}
		})
	}
}

func TestGormDao_FindByUserID(t *testing.T) {
	db := setupTestDB(t)
	model := newGormDao(db)
	ctx := context.Background()

	// 准备测试数据
	userID := "01234567-89ab-7def-0123-456789abc100"
	now := time.Now()
	for i := 0; i < 5; i++ {
		history := &LoginHistory{
			Id:        "01234567-89ab-7def-0123-456789abc" + string(rune('1'+i)) + "0",
			UserID:    userID,
			IP:        "127.0.0.1",
			DeviceType: "Web",
			DeviceID:  "device123",
			UserAgent: "Mozilla/5.0",
			LoginAt:   now.Add(time.Duration(i) * time.Minute),
		}
		_, err := model.Insert(ctx, history)
		if err != nil {
			t.Fatalf("准备测试数据失败: %v", err)
		}
	}

	tests := []struct {
		name     string
		userID   string
		page     int
		pageSize int
		wantLen  int
		wantErr  bool
	}{
		{
			name:     "分页查询-第1页",
			userID:   userID,
			page:     1,
			pageSize: 2,
			wantLen:  2,
			wantErr:  false,
		},
		{
			name:     "分页查询-第2页",
			userID:   userID,
			page:     2,
			pageSize: 2,
			wantLen:  2,
			wantErr:  false,
		},
		{
			name:     "分页查询-最后一页",
			userID:   userID,
			page:     3,
			pageSize: 2,
			wantLen:  1,
			wantErr:  false,
		},
		{
			name:     "查询不存在的用户",
			userID:   "00000000-0000-0000-0000-000000000000",
			page:     1,
			pageSize: 10,
			wantLen:  0,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, total, err := model.FindByUserID(ctx, tt.userID, tt.page, tt.pageSize)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindByUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(got) != tt.wantLen {
				t.Errorf("FindByUserID() len = %v, want %v", len(got), tt.wantLen)
			}
			// 查询不存在的用户时，total 应该是 0，而不是 5
			expectedTotal := int64(5)
			if tt.name == "查询不存在的用户" {
				expectedTotal = 0
			}
			if !tt.wantErr && total != expectedTotal {
				t.Errorf("FindByUserID() total = %v, want %v", total, expectedTotal)
			}
		})
	}
}

func TestGormDao_CountByUserID(t *testing.T) {
	db := setupTestDB(t)
	model := newGormDao(db)
	ctx := context.Background()

	// 准备测试数据
	userID := "01234567-89ab-7def-0123-456789abc200"
	for i := 0; i < 3; i++ {
		history := &LoginHistory{
			Id:        "01234567-89ab-7def-0123-456789abc" + string(rune('2'+i)) + "00",
			UserID:    userID,
			IP:        "127.0.0.1",
			DeviceType: "Web",
			DeviceID:  "device123",
			UserAgent: "Mozilla/5.0",
			LoginAt:   time.Now(),
		}
		_, err := model.Insert(ctx, history)
		if err != nil {
			t.Fatalf("准备测试数据失败: %v", err)
		}
	}

	count, err := model.CountByUserID(ctx, userID)
	if err != nil {
		t.Errorf("CountByUserID() error = %v", err)
		return
	}
	if count != 3 {
		t.Errorf("CountByUserID() = %v, want 3", count)
	}
}

func TestGormDao_DeleteOldRecords(t *testing.T) {
	db := setupTestDB(t)
	model := newGormDao(db)
	ctx := context.Background()

	// 准备测试数据
	now := time.Now()
	oldTime := now.Add(-100 * 24 * time.Hour) // 100天前
	newTime := now.Add(-10 * 24 * time.Hour)   // 10天前

	// 创建旧记录
	oldHistory := &LoginHistory{
		Id:        "01234567-89ab-7def-0123-456789abc300",
		UserID:    "01234567-89ab-7def-0123-456789abcde0",
		IP:        "127.0.0.1",
		DeviceType: "Web",
		DeviceID:  "device123",
		UserAgent: "Mozilla/5.0",
		LoginAt:   oldTime,
	}
	_, err := model.Insert(ctx, oldHistory)
	if err != nil {
		t.Fatalf("准备测试数据失败: %v", err)
	}

	// 创建新记录
	newHistory := &LoginHistory{
		Id:        "01234567-89ab-7def-0123-456789abc301",
		UserID:    "01234567-89ab-7def-0123-456789abcde0",
		IP:        "127.0.0.1",
		DeviceType: "Web",
		DeviceID:  "device123",
		UserAgent: "Mozilla/5.0",
		LoginAt:   newTime,
	}
	_, err = model.Insert(ctx, newHistory)
	if err != nil {
		t.Fatalf("准备测试数据失败: %v", err)
	}

	// 删除90天前的记录
	beforeTime := now.Add(-90 * 24 * time.Hour)
	if err := model.DeleteOldRecords(ctx, beforeTime); err != nil {
		t.Errorf("DeleteOldRecords() error = %v", err)
	}

	// 验证旧记录已删除
	_, err = model.FindOne(ctx, "01234567-89ab-7def-0123-456789abc300")
	if err == nil {
		t.Errorf("DeleteOldRecords() 旧记录应该被删除")
	}

	// 验证新记录未删除
	_, err = model.FindOne(ctx, "01234567-89ab-7def-0123-456789abc301")
	if err != nil {
		t.Errorf("DeleteOldRecords() 新记录不应该被删除: %v", err)
	}
}

func TestGormDao_Trans(t *testing.T) {
	db := setupTestDB(t)
	model := newGormDao(db)
	ctx := context.Background()

	// 测试事务成功
	err := model.Trans(ctx, func(ctx context.Context, m Model) error {
		history := &LoginHistory{
			Id:        "01234567-89ab-7def-0123-456789abc400",
			UserID:    "01234567-89ab-7def-0123-456789abcde0",
			IP:        "127.0.0.1",
			DeviceType: "Web",
			DeviceID:  "device123",
			UserAgent: "Mozilla/5.0",
			LoginAt:   time.Now(),
		}
		_, err := m.Insert(ctx, history)
		return err
	})
	if err != nil {
		t.Errorf("Trans() 事务执行失败: %v", err)
	}

	// 验证数据已提交
	_, err = model.FindOne(ctx, "01234567-89ab-7def-0123-456789abc400")
	if err != nil {
		t.Errorf("Trans() 事务提交后数据不存在: %v", err)
	}
}

