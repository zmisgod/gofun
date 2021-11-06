package lunar

import (
	"fmt"
	"testing"
)

func TestNewLunarToSolar(t *testing.T) {
	solarDate := "1994-09-20"
	fmt.Println(SolarToLunar(solarDate))
	fmt.Println(SolarToSimpleLunar(solarDate))

	lunarDate := "1994-08-15"
	fmt.Println(LunarToSolar(lunarDate, false))
}