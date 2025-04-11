package postgres

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/unxai/geonames-service/logger"
	"github.com/unxai/geonames-service/models"
	"go.uber.org/zap"
)

// PostgresStorage 实现了 Storage 接口的 PostgreSQL 存储
type PostgresStorage struct {
	db *sql.DB
}

// NewPostgresStorage 创建一个新的 PostgreSQL 存储实例
func NewPostgresStorage(db *sql.DB) *PostgresStorage {
	return &PostgresStorage{db: db}
}

const batchSize = 1000 // 每批处理的数据量

// SaveLocations 批量保存位置数据
func (s *PostgresStorage) SaveLocations(locations []models.Location) error {
	if len(locations) == 0 {
		return nil
	}

	// 分批处理数据
	for i := 0; i < len(locations); i += batchSize {
		end := i + batchSize
		if end > len(locations) {
			end = len(locations)
		}
		batch := locations[i:end]

		// 开启事务
		tx, err := s.db.Begin()
		if err != nil {
			logger.Logger.Error("开启事务失败", zap.Error(err))
			return fmt.Errorf("开启事务失败: %w", err)
		}
		defer tx.Rollback()

		// 构建批量插入的值占位符
		valueStrings := make([]string, 0, len(batch))
		valueArgs := make([]interface{}, 0, len(batch)*15)
		for j, loc := range batch {
			valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
				j*15+1, j*15+2, j*15+3, j*15+4, j*15+5, j*15+6, j*15+7, j*15+8, j*15+9, j*15+10, j*15+11, j*15+12, j*15+13, j*15+14, j*15+15))
			valueArgs = append(valueArgs,
				loc.GeonameID, loc.Name, loc.ASCII_Name, loc.AlternateNames, loc.Latitude, loc.Longitude,
				loc.FeatureClass, loc.FeatureCode, loc.CountryCode, loc.Admin1Code, loc.Admin2Code,
				loc.Population, loc.Elevation, loc.TimeZone, loc.ModificationDate)
		}

		// 构建完整的SQL语句
		sql := fmt.Sprintf(`
		INSERT INTO locations (
			geoname_id, name, ascii_name, alternate_names, latitude, longitude,
			feature_class, feature_code, country_code, admin1_code, admin2_code,
			population, elevation, timezone, modification_date
		) VALUES %s
		ON CONFLICT (geoname_id) DO UPDATE SET
			name = EXCLUDED.name,
			ascii_name = EXCLUDED.ascii_name,
			alternate_names = EXCLUDED.alternate_names,
			latitude = EXCLUDED.latitude,
			longitude = EXCLUDED.longitude,
			feature_class = EXCLUDED.feature_class,
			feature_code = EXCLUDED.feature_code,
			country_code = EXCLUDED.country_code,
			admin1_code = EXCLUDED.admin1_code,
			admin2_code = EXCLUDED.admin2_code,
			population = EXCLUDED.population,
			elevation = EXCLUDED.elevation,
			timezone = EXCLUDED.timezone,
			modification_date = EXCLUDED.modification_date
		`, strings.Join(valueStrings, ","))

		// 执行批量插入
		_, err = tx.Exec(sql, valueArgs...)
		if err != nil {
			logger.Logger.Error("执行批量插入失败",
				zap.Int("batch_start", i),
				zap.Int("batch_size", len(batch)),
				zap.Error(err))
			return fmt.Errorf("执行批量插入失败: %w", err)
		}

		// 提交事务
		if err := tx.Commit(); err != nil {
			logger.Logger.Error("提交事务失败",
				zap.Int("batch_start", i),
				zap.Int("batch_size", len(batch)),
				zap.Error(err))
			return fmt.Errorf("提交事务失败: %w", err)
		}

		logger.Logger.Info("成功处理一批数据",
			zap.Int("batch_start", i),
			zap.Int("batch_size", len(batch)))
	}

	return nil
}
