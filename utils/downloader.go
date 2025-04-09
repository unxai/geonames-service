package utils

import (
	"archive/zip"
	"bufio"
	"bytes"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/unxai/geonames-service/models"
)

// DownloadGeoData 修改为返回数据而不是保存文件
func DownloadGeoData() ([]models.Location, error) {
	// 发起 HTTP 请求获取数据
	resp, err := http.Get("http://download.geonames.org/export/dump/allCountries.zip")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 创建一个临时的内存缓冲区来存储 zip 数据
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 从内存中读取 zip 文件
	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return nil, err
	}

	var locations []models.Location
	// 解析 zip 文件中的数据
	for _, file := range zipReader.File {
		if file.Name == "allCountries.txt" {
			rc, err := file.Open()
			if err != nil {
				return nil, err
			}
			defer rc.Close()

			scanner := bufio.NewScanner(rc)
			for scanner.Scan() {
				fields := strings.Split(scanner.Text(), "\t")
				if len(fields) > 10 {
					geonameID, _ := strconv.Atoi(fields[0])
					lat, _ := strconv.ParseFloat(fields[4], 64)
					lon, _ := strconv.ParseFloat(fields[5], 64)
					pop, _ := strconv.Atoi(fields[14])

					location := models.Location{
						GeonameID:    geonameID,
						Name:         fields[1],
						ASCII_Name:   fields[2],
						Latitude:     lat,
						Longitude:    lon,
						CountryCode:  fields[8],
						Population:   pop,
						FeatureClass: fields[6],
						FeatureCode:  fields[7],
					}
					locations = append(locations, location)
				}
			}
		}
	}
	return locations, nil
}
