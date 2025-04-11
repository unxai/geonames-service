package models

type Location struct {
	GeonameID        int     `json:"geoname_id" db:"geoname_id"`
	Name             string  `json:"name" db:"name"`
	ASCII_Name       string  `json:"ascii_name" db:"ascii_name"`
	AlternateNames   string  `json:"alternate_names" db:"alternate_names"`
	Latitude         float64 `json:"latitude" db:"latitude"`
	Longitude        float64 `json:"longitude" db:"longitude"`
	FeatureClass     string  `json:"feature_class" db:"feature_class"`
	FeatureCode      string  `json:"feature_code" db:"feature_code"`
	CountryCode      string  `json:"country_code" db:"country_code"`
	Admin1Code       string  `json:"admin1_code" db:"admin1_code"`
	Admin2Code       string  `json:"admin2_code" db:"admin2_code"`
	Population       int     `json:"population" db:"population"`
	Elevation        int     `json:"elevation" db:"elevation"`
	TimeZone         string  `json:"timezone" db:"timezone"`
	ModificationDate string  `json:"modification_date" db:"modification_date"`
}
