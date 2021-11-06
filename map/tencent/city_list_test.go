package tencent

import (
	"log"
	"os"
	"testing"
)

func TestNewCityList(t *testing.T) {
	_file, err := os.Open("./city_list.json")
	if err != nil {
		t.Error(err)
	}
	re, err := NewCity(_file)
	if err != nil {
		t.Error(err)
	}
	province := re.GetAllProvince()
	log.Println(province)

	city := re.GetCitiesByProvince("澳门特别行政区")
	log.Println(city)

	districts := re.GetDistrictByCityName("澳门特别行政区", "澳门特别行政区")
	log.Println(districts)
}