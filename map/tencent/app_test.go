package tencent

import (
	"log"
	"os"
	"testing"
)

func TestNewCity(t *testing.T) {
	_file, err := os.Open("./res.json")
	if err != nil {
		t.Error(err)
	}
	re, err := NewCity(_file)
	if err != nil {
		t.Error(err)
	}
	province := re.GetAllProvince()
	log.Println(province)

	city := re.GetCitiesByProvince("台湾省")
	log.Println(city)

	districts := re.GetDistrictByCityName("台湾省", "台北市")
	log.Println(districts)
}