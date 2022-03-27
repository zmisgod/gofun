package jd_union

import (
	"context"
	"encoding/json"
	"errors"
	"log"
)

func (app *App) JdUnionOpenGoodsMaterialQuery(ctx context.Context, params map[string]interface{}) (result *JdUnionOpenGoodsMaterialQueryResult, err error) {
	body, err := app.Request(ctx, JdUnionOpenGoodsMaterialQueryPath, map[string]interface{}{"goodsReq": params})

	resp := &JdUnionOpenGoodsMaterialQueryTopLevel{}
	if err != nil {
		log.Println(string(body))
		return
	}
	if err = json.Unmarshal(body, resp); err != nil {
		return
	}
	if resp.JdUnionOpenGoodsMaterialQueryResponse.Result != "" {
		result = &JdUnionOpenGoodsMaterialQueryResult{}
		if err = json.Unmarshal([]byte(resp.JdUnionOpenGoodsMaterialQueryResponse.Result), result); err != nil {
			return
		}
		if result.Code != 200 {
			err = errors.New(result.Message)
		}
	} else {
		err = ResultIsNullError
	}
	return
}

//JdUnionOpenSellingPromotionGet 转链获取，支持工具商
func (app App) JdUnionOpenSellingPromotionGet(ctx context.Context, params map[string]interface{}) (result *JdUnionOpenSellingPromotionGetResult, err error) {
	body, err := app.Request(ctx, JdUnionOpenSellingPromotionGetPath, map[string]interface{}{"req": params})

	resp := &JdUnionOpenSellingPromotionGetTopLevel{}
	if err != nil {
		log.Println(string(body))
		return
	}
	if err = json.Unmarshal(body, resp); err != nil {
		return
	}
	if resp.JdUnionOpenSellingPromotionGetResponse.Result != "" {
		result = &JdUnionOpenSellingPromotionGetResult{}
		if err = json.Unmarshal([]byte(resp.JdUnionOpenSellingPromotionGetResponse.Result), result); err != nil {
			return
		}
		if result.Code != 200 {
			err = errors.New(result.Message)
		}
	} else {
		err = ResultIsNullError
	}
	return
}
