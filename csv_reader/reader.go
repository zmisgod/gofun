package csv_reader

import (
	"encoding/csv"
	"github.com/axgle/mahonia"
	"io/ioutil"
	"os"
	"strings"
)

var DefaultCharset = "gbk"
var CharsetUTF8 = "utf-8"

type CsvReader struct {
	FilePath string
	Charset string
}

func NewCsvReader(filePath string) (*CsvReader, error) {
	_, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}
	return &CsvReader{
		FilePath: filePath,
		Charset:  DefaultCharset,
	}, nil
}

func (a *CsvReader) SetCharset(charset string) {
	if charset != "" {
		a.Charset = charset
	}
}

func (a *CsvReader) GetValues() ([][]string, error ) {
	resultValueArray := make([][]string, 0)
	csvByteString , err := ioutil.ReadFile(a.FilePath)
	if err != nil {
		return resultValueArray, err
	}
	decoder := strings.NewReader(string(csvByteString))
	r := csv.NewReader(decoder)
	if a.Charset != CharsetUTF8 {
		decoder := mahonia.NewDecoder(a.Charset)
		r = csv.NewReader(decoder.NewReader(strings.NewReader(string(csvByteString))))
	}
	reads, err := r.ReadAll()
	if err != nil {
		return resultValueArray, err
	}
	resultValueArray = reads
	return resultValueArray, nil
}