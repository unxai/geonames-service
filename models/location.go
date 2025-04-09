package models

type Location struct {
    GeonameID     int     `json:"geoname_id" db:"geoname_id"`
    Name          string  `json:"name" db:"name"`
    ASCII_Name    string  `json:"ascii_name" db:"ascii_name"`
    Latitude      float64 `json:"latitude" db:"latitude"`
    Longitude     float64 `json:"longitude" db:"longitude"`
    CountryCode   string  `json:"country_code" db:"country_code"`
    Population    int     `json:"population" db:"population"`
    FeatureClass  string  `json:"feature_class" db:"feature_class"`
    FeatureCode   string  `json:"feature_code" db:"feature_code"`
}