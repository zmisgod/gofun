package easy_http_client

import (
	"context"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

var ctx context.Context
var sessionId string
var postFromDataUrl string

func init() {
	ctx = context.Background()

	err := godotenv.Load("./../.env")
	if err != nil {
		log.Fatal(err)
	}

	sessionId = os.Getenv("easy_http_client_sessionId")
	postFromDataUrl = os.Getenv("easy_http_client_post_form_url")
}

func TestHttpClient_GET(t *testing.T) {
	resp, err := NewHttpClient("https://api.zmis.me/", HttpClientMethodGet, nil, 1, "").HttpClient(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	rest, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(rest))
}

func TestHttpClient_POSTJson(t *testing.T) {
	jsonString := `{"email":"111@111.com","password":"123456"}`
	nC := NewHttpClient("https://api.zmis.me/", HttpClientMethodPost, nil, 1, "")
	nC.SetJsonData(ctx, jsonString)
	resp, err := nC.HttpClient(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	rest, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(rest))
}

func TestHttpClient_POSTForm(t *testing.T) {
	formData := map[string]string{
		"latitude":"12.0",
		"longitude":"13.0",
		"city":"安庆",
		"movieId":"1298367",
	}
	header := map[string]string{"sessionId":sessionId}
	nC := NewHttpClient(postFromDataUrl, HttpClientMethodPost, header, 1, "")
	nC.SetPostFormData(ctx, formData)
	resp, err := nC.HttpClient(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	rest, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(rest))
}

func TestHttpClient_HEAD(t *testing.T) {
	resp, err := NewHttpClient("https://static.zmis.me/public/images/2021/4/28/407aaa03d57c02bebc5901a27e814bc3.jpg", HttpClientMethodHead, nil, 1, "").HttpClient(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	for _, v := range resp.Header {
		t.Log(v)
	}
}
