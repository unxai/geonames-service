package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"

	_ "github.com/lib/pq"
	"github.com/unxai/geonames-service/config"
	"github.com/unxai/geonames-service/storage/postgres"
)

var (
	db   *sql.DB
	once sync.Once
)

// GetDB 返回数据库连接池的单例实例
func GetDB() *sql.DB {
	once.Do(func() {
		var err error
		_, err = config.LoadConfig()
		if err != nil {
			panic(fmt.Sprintf("加载配置失败: %v", err))
		}

		// 初始化数据库连接
		db, err = sql.Open("postgres", config.GetDSN())
		if err != nil {
			panic(fmt.Sprintf("连接数据库失败: %v", err))
		}

		// 设置连接池参数
		db.SetMaxOpenConns(50)         // 最大打开连接数
		db.SetMaxIdleConns(10)         // 最大空闲连接数
		db.SetConnMaxLifetime(30 * 60) // 连接最大生命周期（30分钟）

		// 测试连接
		err = db.Ping()
		if err != nil {
			panic(fmt.Sprintf("数据库连接测试失败: %v", err))
		}
	})

	return db
}

// GetStorage 返回存储实例
func GetStorage() *postgres.PostgresStorage {
	return postgres.NewPostgresStorage(GetDB())
}

// CloseDB 关闭数据库连接池
func CloseDB() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

// Migrate 执行数据库迁移
func Migrate() error {
	db := GetDB()

	// 读取migrations目录下的SQL文件
	files, err := os.ReadDir("db/migrations")
	if err != nil {
		return fmt.Errorf("读取迁移文件失败: %w", err)
	}

	// 按文件名排序执行迁移
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		content, err := os.ReadFile(filepath.Join("db/migrations", file.Name()))
		if err != nil {
			return fmt.Errorf("读取迁移文件内容失败: %w", err)
		}

		if _, err := db.Exec(string(content)); err != nil {
			return fmt.Errorf("执行迁移脚本失败(%s): %w", file.Name(), err)
		}
	}

	return nil
}
