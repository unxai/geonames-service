-- 创建locations表
CREATE TABLE IF NOT EXISTS locations (
    geoname_id BIGINT PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    ascii_name VARCHAR(200),
    alternate_names TEXT,
    latitude DECIMAL(10, 7) NOT NULL,
    longitude DECIMAL(10, 7) NOT NULL,
    feature_class CHAR(1),
    feature_code VARCHAR(10),
    country_code CHAR(2),
    admin1_code VARCHAR(20),
    admin2_code VARCHAR(80),
    population BIGINT,
    elevation INTEGER,
    timezone VARCHAR(40),
    modification_date DATE
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_locations_country_code ON locations(country_code);
CREATE INDEX IF NOT EXISTS idx_locations_feature_class ON locations(feature_class);
CREATE INDEX IF NOT EXISTS idx_locations_admin1_code ON locations(admin1_code);
CREATE INDEX IF NOT EXISTS idx_locations_admin2_code ON locations(admin2_code);
CREATE INDEX IF NOT EXISTS idx_locations_timezone ON locations(timezone);