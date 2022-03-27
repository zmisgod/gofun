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
		"materialId": "https://item.jd.com/11144230.html",
		"siteId":     app.ID,
	})
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(res)
}
