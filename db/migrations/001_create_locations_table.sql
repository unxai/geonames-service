-- 创建locations表
CREATE TABLE IF NOT EXISTS locations (
    geoname_id BIGINT PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    ascii_name VARCHAR(200),
    latitude DECIMAL(10, 7) NOT NULL,
    longitude DECIMAL(10, 7) NOT NULL,
    country_code CHAR(2),
    population BIGINT,
    feature_class CHAR(1),
    feature_code VARCHAR(10)
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_locations_country_code ON locations(country_code);
CREATE INDEX IF NOT EXISTS idx_locations_feature_class ON locations(feature_class);