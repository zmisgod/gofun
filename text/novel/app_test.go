package novel

import (
	"github.com/zmisgod/gofun/utils"
	"log"
	"testing"
)

func TestFetchChapter(t *testing.T) {
	var textName = "丹武双绝"
	var url = "https://www.18xs.org/book_25306/"
	_, err := utils.CreateFolder(textName)
	if err != nil {
		log.Fatal(err)
	}
	FetchChapter(textName, url)
}