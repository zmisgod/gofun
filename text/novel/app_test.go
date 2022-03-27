package novel

import (
	"context"
	"fmt"
	"testing"
)

func TestFetchChapter(t *testing.T) {
	var textName = "丹武双绝"
	var url = "https://www.18xs.org/book_25306/"
	//_, err := utils.CreateFolder(textName)
	//if err != nil {
	//	log.Fatal(err)
	//}
	var nil int
	nil = 1
	fmt.Println(nil)
	FetchChapter(context.Background(), textName, url)
}
