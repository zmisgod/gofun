package jd_union

import (
	"context"
	"fmt"
	"log"
	"testing"
)

var app = &App{
	ID:     "",
	Key:    "",
	Secret: "",
}

var ctx = context.Background()

func TestApp_JdUnionOpenGoodsMaterialQuery(t *testing.T) {
	res, err := app.JdUnionOpenGoodsMaterialQuery(ctx, map[string]interface{}{
		"eliteId": 1,
	})
	if err != nil {
		log.Fatalln(err)
	}
	if len(res.Data) > 0 {
		for _, v := range res.Data {
			fmt.Println(v)
		}
	}
}

func TestApp_JdUnionOpenSellingPromotionGet(t *testing.T) {
	res, err := app.JdUnionOpenSellingPromotionGet(ctx, map[string]interface{}{
		"materialId": "https://item.jd.com/10044136889444.html",
		"couponUrl":  "https://coupon.m.jd.com/coupons/show.action?linkKey=AAROH_xIpeffAs_-naABEFoeQNXepOFXfqKPPNJ6Bv6t-rAJQluMOCkNpEZEhpiJKaAId0ebBw38GE5siHZJ2od68RBsMA",
		"siteId":     app.ID,
	})
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(res)
}

func TestApp_JdUnionOpenGoodsQuery(t *testing.T) {
	res, err := app.JdUnionOpenGoodsQuery(ctx, map[string]interface{}{
		"keyword":   "手机",
		"skuIds":    []uint64{},
		"pageSize":  20,
		"pageIndex": 1,
	})
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(res)
}
