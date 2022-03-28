package jd_union

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/zmisgod/gofun/utils"
	"io/ioutil"
	"log"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

type App struct {
	ID     string
	Name   string
	Key    string
	Secret string
}

type JdUnionErrResp struct {
	ErrorResponse ErrorResponse `json:"errorResponse"`
}

type ErrorResponse struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

var (
	ResultIsNullError error = errors.New("result is null")
)

const (
	JdUnionOpenGoodsMaterialQueryPath  string = "jd.union.open.goods.material.query"
	JdUnionOpenSellingPromotionGetPath string = "jd.union.open.promotion.common.get"
)

const RouterURL = "https://api.jd.com/routerjson"

func (app *App) Request(ctx context.Context, method string, paramJSON map[string]interface{}) ([]byte, error) {
	// common params
	params := map[string]interface{}{}
	params["method"] = method
	params["app_key"] = app.Key
	params["format"] = "json"
	params["v"] = "1.0"
	params["sign_method"] = "md5"
	params["timestamp"] = time.Now().Format("2006-01-02 15:04:05")

	// api params
	paramJSONStr, _ := json.Marshal(paramJSON)
	params["360buy_param_json"] = string(paramJSONStr)
	params["sign"] = GetSign(app.Secret, params)
	log.Printf("Request: %s, %v", RouterURL, params)

	resp, err := utils.HttpClient(ctx, utils.HttpClientMethodPost, RouterURL, 3*time.Second, "",
		ioutil.NopCloser(strings.NewReader(app.Values(params).Encode())),
		utils.ContentTypeFormUrlEncode,
	)
	log.Printf("Responce:%v %v", resp, err)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	body, err := ioutil.ReadAll(resp.Body)
	log.Printf("Responce Body:%v ", string(body))
	if err != nil {
		return nil, err
	}

	jdErr := JdUnionErrResp{}
	if err := json.Unmarshal(body, &jdErr); err != nil {
		return nil, err
	}

	if jdErr.ErrorResponse.Code != "" {
		return nil, errors.New(jdErr.ErrorResponse.Msg)
	}
	return body, nil
}

func (app *App) Values(params map[string]interface{}) url.Values {
	_val := url.Values{}
	for key, val := range params {
		_val.Add(key, GetString(val))
	}
	return _val
}

func GetSign(clientSecret string, p map[string]interface{}) string {
	var keys []string
	for k := range p {
		if k != "sign" && k != "access_token" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	signStr := clientSecret
	for _, key := range keys {
		signStr += key + GetString(p[key])
	}
	signStr += clientSecret
	return md5Hash(signStr)
}

func GetString(i interface{}) string {
	switch v := i.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case bool:
		return strconv.FormatBool(v)
	default:
		bytes, _ := json.Marshal(v)
		return string(bytes)
	}
}

func md5Hash(signStr string) string {
	h := md5.New()
	h.Write([]byte(signStr))
	cipherStr := h.Sum(nil)
	return strings.ToUpper(hex.EncodeToString(cipherStr))
}
