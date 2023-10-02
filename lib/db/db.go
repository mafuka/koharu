// db 包提供了数据库操作的相关方法。
package db

import (
	"fmt"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

// Init 初始化数据库。
func Init(dbPath string) error {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return err
	}

	return db.AutoMigrate()
}

// Close 关闭数据库。
func Close() error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

// PrivateMsg 表示私聊会话的消息表模型。
type PrivateMsg struct {
	gorm.Model
	UserID  int64  // 用户 ID
	Message string // 消息内容
	Time    time.Time
}

// GroupMsg 表示群聊会话的消息表模型。
type GroupMsg struct {
	gorm.Model
	GroupID int64  // 群组 ID
	UserID  int64  // 用户 ID
	Message string // 消息内容
	Time    time.Time
}

// GetDB 返回数据库实例。
func GetDB() *gorm.DB {
	return db
}

// CheckIfTableExists 检查表是否存在。
func CheckIfTableExists(tableName string) bool {
	return db.Migrator().HasTable(tableName)
}

// CreateTable 创建表。其中，tableName 为表名，model 为表模型。
func CreateTable(tableName string, model interface{}) error {
	if CheckIfTableExists(tableName) {
		return fmt.Errorf("table %s already exists", tableName)
	}

	return db.Table(tableName).AutoMigrate(model)
}

// DropTable 删除表。其中，tableName 为表名。
func DropTable(tableName string) error {
	if !CheckIfTableExists(tableName) {
		return fmt.Errorf("table %s does not exist", tableName)
	}

	return db.Migrator().DropTable(tableName)
}

// GetTable 获取表。其中，tableName 为表名。
func GetTable(tableName string) *gorm.DB {
	return db.Table(tableName)
}
