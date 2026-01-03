package svc

import (
	"database/sql"
	"fmt"

	"github.com/jinguoxing/idrm-cursor-demo/api/internal/config"
	"github.com/jinguoxing/idrm-cursor-demo/model/auth/login_history"
	"github.com/jinguoxing/idrm-cursor-demo/model/auth/users"
	"github.com/jinguoxing/idrm-cursor-demo/model/resource_catalog/category"
	"github.com/jinguoxing/idrm-cursor-demo/pkg/db"
	"github.com/jinguoxing/idrm-cursor-demo/pkg/device"
	"github.com/jinguoxing/idrm-cursor-demo/pkg/jwt"
	"github.com/jinguoxing/idrm-cursor-demo/pkg/sms"

	_ "github.com/go-sql-driver/mysql"
	redisv9 "github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config

	// Model层（使用接口类型，支持自动ORM选择）
	CategoryModel    category.Model
	UserModel        users.Model
	LoginHistoryModel login_history.Model

	// 服务层
	SMSService   sms.Service
	JWTService   jwt.Service
	DeviceParser func(string) (string, string) // deviceType, deviceID
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 使用 Auth 数据库（用户认证功能）
	dbCfg := c.DB.Auth
	
	// 1. 初始化 sqlx 连接（作为备用）
	var sqlConn *sql.DB
	var sqlxErr error
	dsn := buildDSN(dbCfg)
	logx.Infof("尝试连接数据库(SQLx): %s:%d/%s", dbCfg.Host, dbCfg.Port, dbCfg.Database)
	conn := sqlx.NewMysql(dsn)
	sqlConn, sqlxErr = conn.RawDB()
	if sqlxErr != nil {
		logx.Errorf("SQLx RawDB获取失败: %v, DSN: %s", sqlxErr, dsn)
		sqlConn = nil
	} else {
		logx.Info("SQLx 连接成功")
	}

	// 2. 初始化 gorm 连接（优先）
	var gormDB *gorm.DB
	var gormErr error
	logx.Infof("尝试连接数据库(GORM): %s:%d/%s", dbCfg.Host, dbCfg.Port, dbCfg.Database)
	gormDB, gormErr = db.InitGorm(dbCfg)
	if gormErr != nil {
		logx.Errorf("GORM初始化失败: %v", gormErr)
		gormDB = nil
	} else {
		logx.Info("GORM 连接成功")
	}

	// 如果两个都失败，提前panic
	if sqlConn == nil && gormDB == nil {
		panic(fmt.Sprintf("数据库连接失败！SQLx错误: %v, GORM错误: %v", sqlxErr, gormErr))
	}

	// 3. 初始化 Redis go-redis 客户端（用于SMS和JWT服务）
	var goRedisClient *redisv9.Client
	redisAddr := fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port)
	goRedisClient = redisv9.NewClient(&redisv9.Options{
		Addr:     redisAddr,
		Password: c.Redis.Password,
		DB:       c.Redis.DB,
	})
	logx.Infof("Redis客户端初始化完成: %s (DB: %d)", redisAddr, c.Redis.DB)

	// 5. 使用工厂自动选择ORM（gorm优先，sqlx降级）
	categoryModel := category.NewModel(sqlConn, gormDB)
	userModel := users.NewModel(sqlConn, gormDB)
	loginHistoryModel := login_history.NewModel(sqlConn, gormDB)

	// 6. 初始化服务
	var smsService sms.Service
	var jwtService jwt.Service
	if goRedisClient != nil {
		smsService = sms.NewRedisStore(goRedisClient)
		jwtService = jwt.NewJWTService(c.Auth.AccessSecret, goRedisClient)
	} else {
		logx.Warn("Redis未配置，SMS和JWT服务将无法正常工作")
	}

	return &ServiceContext{
		Config:           c,
		CategoryModel:    categoryModel,
		UserModel:        userModel,
		LoginHistoryModel: loginHistoryModel,
		SMSService:       smsService,
		JWTService:       jwtService,
		DeviceParser:     device.ParseDeviceInfo,
	}
}

// buildDSN 构建 sqlx 的 DSN
func buildDSN(cfg db.Config) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.Charset,
	)
}
