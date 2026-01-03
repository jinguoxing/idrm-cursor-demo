package users

import (
	"context"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB 创建测试数据库
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("创建测试数据库失败: %v", err)
	}

	// 自动迁移
	if err := db.AutoMigrate(&User{}); err != nil {
		t.Fatalf("数据库迁移失败: %v", err)
	}

	return db
}

func TestGormDao_Insert(t *testing.T) {
	db := setupTestDB(t)
	model := newGormDao(db)
	ctx := context.Background()

	tests := []struct {
		name    string
		user    *User
		wantErr bool
	}{
		{
			name: "创建用户成功",
			user: &User{
				Id:          "01234567-89ab-7def-0123-456789abc000",
				Mobile:      "13800138000",
				PasswordHash: "hashed_password",
				Status:      1,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			wantErr: false,
		},
		{
			name: "手机号重复",
			user: &User{
				Id:          "01234567-89ab-7def-0123-456789abc000",
				Mobile:      "13800138000", // 重复手机号
				PasswordHash: "hashed_password2",
				Status:      1,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := model.Insert(ctx, tt.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("Insert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("Insert() 返回结果为空")
			}
		})
	}
}

func TestGormDao_FindOne(t *testing.T) {
	db := setupTestDB(t)
	model := newGormDao(db)
	ctx := context.Background()

	// 先创建一个用户
	userID := "01234567-89ab-7def-0123-456789abc001"
	user := &User{
		Id:          userID,
		Mobile:      "13800138001",
		PasswordHash: "hashed_password",
		Status:      1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	_, err := model.Insert(ctx, user)
	if err != nil {
		t.Fatalf("准备测试数据失败: %v", err)
	}

	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{
			name:    "查询存在的用户",
			id:      userID,
			wantErr: false,
		},
		{
			name:    "查询不存在的用户",
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

func TestGormDao_FindByMobile(t *testing.T) {
	db := setupTestDB(t)
	model := newGormDao(db)
	ctx := context.Background()

	// 先创建一个用户
	mobile := "13800138002"
	user := &User{
		Id:          "01234567-89ab-7def-0123-456789abc002",
		Mobile:      mobile,
		PasswordHash: "hashed_password",
		Status:      1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	_, err := model.Insert(ctx, user)
	if err != nil {
		t.Fatalf("准备测试数据失败: %v", err)
	}

	tests := []struct {
		name    string
		mobile  string
		wantErr bool
	}{
		{
			name:    "查询存在的手机号",
			mobile:  mobile,
			wantErr: false,
		},
		{
			name:    "查询不存在的手机号",
			mobile:  "13999999999",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := model.FindByMobile(ctx, tt.mobile)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindByMobile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("FindByMobile() 返回结果为空")
			}
			if !tt.wantErr && got.Mobile != tt.mobile {
				t.Errorf("FindByMobile() Mobile = %v, want %v", got.Mobile, tt.mobile)
			}
		})
	}
}

func TestGormDao_Update(t *testing.T) {
	db := setupTestDB(t)
	model := newGormDao(db)
	ctx := context.Background()

	// 先创建一个用户
	user := &User{
		Id:          "01234567-89ab-7def-0123-456789abc003",
		Mobile:      "13800138003",
		PasswordHash: "hashed_password",
		Status:      1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	_, err := model.Insert(ctx, user)
	if err != nil {
		t.Fatalf("准备测试数据失败: %v", err)
	}

	// 更新用户信息
	user.Status = 2
	if err := model.Update(ctx, user); err != nil {
		t.Errorf("Update() error = %v", err)
	}

	// 验证更新
	updated, err := model.FindOne(ctx, user.Id)
	if err != nil {
		t.Fatalf("查询更新后的用户失败: %v", err)
	}
	if updated.Status != 2 {
		t.Errorf("Update() Status = %v, want 2", updated.Status)
	}
}

func TestGormDao_UpdatePassword(t *testing.T) {
	db := setupTestDB(t)
	model := newGormDao(db)
	ctx := context.Background()

	// 先创建一个用户
	userID := "01234567-89ab-7def-0123-456789abc004"
	user := &User{
		Id:          userID,
		Mobile:      "13800138004",
		PasswordHash: "old_password",
		Status:      1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	_, err := model.Insert(ctx, user)
	if err != nil {
		t.Fatalf("准备测试数据失败: %v", err)
	}

	// 更新密码
	newPassword := "new_password"
	if err := model.UpdatePassword(ctx, userID, newPassword); err != nil {
		t.Errorf("UpdatePassword() error = %v", err)
	}

	// 验证更新
	updated, err := model.FindOne(ctx, userID)
	if err != nil {
		t.Fatalf("查询更新后的用户失败: %v", err)
	}
	if updated.PasswordHash != newPassword {
		t.Errorf("UpdatePassword() PasswordHash = %v, want %v", updated.PasswordHash, newPassword)
	}
}

func TestGormDao_UpdateStatus(t *testing.T) {
	db := setupTestDB(t)
	model := newGormDao(db)
	ctx := context.Background()

	// 先创建一个用户
	userID := "01234567-89ab-7def-0123-456789abc005"
	user := &User{
		Id:          userID,
		Mobile:      "13800138005",
		PasswordHash: "hashed_password",
		Status:      1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	_, err := model.Insert(ctx, user)
	if err != nil {
		t.Fatalf("准备测试数据失败: %v", err)
	}

	// 更新状态
	if err := model.UpdateStatus(ctx, userID, 3); err != nil {
		t.Errorf("UpdateStatus() error = %v", err)
	}

	// 验证更新
	updated, err := model.FindOne(ctx, userID)
	if err != nil {
		t.Fatalf("查询更新后的用户失败: %v", err)
	}
	if updated.Status != 3 {
		t.Errorf("UpdateStatus() Status = %v, want 3", updated.Status)
	}
}

func TestGormDao_Trans(t *testing.T) {
	db := setupTestDB(t)
	model := newGormDao(db)
	ctx := context.Background()

	// 测试事务成功
	err := model.Trans(ctx, func(ctx context.Context, m Model) error {
		user1 := &User{
			Id:          "01234567-89ab-7def-0123-456789abc006",
			Mobile:      "13800138006",
			PasswordHash: "hashed_password",
			Status:      1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		_, err := m.Insert(ctx, user1)
		return err
	})
	if err != nil {
		t.Errorf("Trans() 事务执行失败: %v", err)
	}

	// 验证数据已提交
	_, err = model.FindOne(ctx, "01234567-89ab-7def-0123-456789abc006")
	if err != nil {
		t.Errorf("Trans() 事务提交后数据不存在: %v", err)
	}

	// 测试事务回滚
	err = model.Trans(ctx, func(ctx context.Context, m Model) error {
		user2 := &User{
			Id:          "01234567-89ab-7def-0123-456789abc007",
			Mobile:      "13800138006", // 重复手机号，触发错误
			PasswordHash: "hashed_password",
			Status:      1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		_, err := m.Insert(ctx, user2)
		return err
	})
	if err == nil {
		t.Errorf("Trans() 应该返回错误")
	}

	// 验证数据未提交（回滚）
	_, err = model.FindOne(ctx, "01234567-89ab-7def-0123-456789abc007")
	if err == nil {
		t.Errorf("Trans() 事务回滚后数据不应该存在")
	}
}

