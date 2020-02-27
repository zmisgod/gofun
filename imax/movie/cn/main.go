package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

func main (){
	text := `<div class="views-row views-row-1 views-row-odd views-row-first"> 
            <div class="views-field field views-field-field-block-image field-name-field-block-image"> 
             <div class="field-content">
                <img typeof="foaf:Image" src="http://www.imax.cn/public/attachment/201808/16/17/5b7541de6fd71.jpg?itok=6v0zSXju" width="400" height="400" alt="">
             </div> 
            </div> 
            <div class="group-description views-fieldset" data-module="views_fieldsets"> 
             <div class="views-field field views-field-field-us-status field-name-field-us-status">
              <div class="field-content">
                已经下映              </div>
             </div> 
             <div class="views-field field views-field-title-field field-name-title-field">
              <div class="field-content">
               
               <a href="/moviedetailed/id-309">海王</a>
              </div>
             </div> 
             <div class="views-field field views-field-view-node field-name-view-node">
              <span class="field-content"><a href="/moviedetailed/id-309" class="btn btn-primary">更多详情</a><!-- <a href="/index.php?ctl=moviedetailed&act=index&id=309" class="btn btn-primary">更多信息</a> --></span>
             </div> 
            </div> 
           </div>`
	rd := strings.NewReader(text)
	doc, err := goquery.NewDocumentFromReader(rd)
	if err != nil {
		fmt.Println(err)
	}else{
		content := doc.Find(".views-row")
		content.EachWithBreak(func(i int, selection *goquery.Selection) bool {
			image := selection.Find(".views-field")
			image.EachWithBreak(func(j int, imgSelection *goquery.Selection) bool {
				rest, exist := imgSelection.Find(".field-content img").Attr("src")
				if !exist {
					return false
				}
				fmt.Println(rest)
				return true
			})
			intro := selection.Find(".group-description .views-field")
			intro.EachWithBreak(func(j int, imgSelection *goquery.Selection) bool {
				if j == 1 {
					obj := imgSelection.Find(".field-content a")
					url, exist := obj.Attr("href")
					if exist {
						fmt.Println(url)
					}
					rest,err := obj.Html()
					if err != nil {
						return false
					}
					fmt.Println(rest)
					return true
				}
				return true
			})
			return true
		})
	}
}