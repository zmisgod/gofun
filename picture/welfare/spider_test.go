package welfare

import (
	"github.com/joho/godotenv"
	"log"
	"testing"
)

func TestNewSpider(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	NewSpider(1, 2)
}