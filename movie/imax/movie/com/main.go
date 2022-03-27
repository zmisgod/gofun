package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

func main() {
	text := "<div class=\"view view-archived-movies view-id-archived_movies view-display-id-page view-dom-id-3badf5455744dd4f293f12e238bc6b1c wrapper-container-box\">\n  <div class=\"header-wrapper-box\">\n    <div class=\"archived-movies-header\">\n            <h1>Archived Movies</h1>\n                    <div class=\"view-header\">\n          <p>IMAX's library of documentaries and major Hollywood features is an amazing testament to the art, science, and passion of filmmakers and film lovers around the globe. Here you can find information on every single IMAX release from over the decades. Do you remember your first IMAX film?</p>\n        </div>\n          </div>\n\n    <hr>\n  </div>\n      <div class=\"view-filters\">\n      <form class=\"ctools-auto-submit-full-form\" action=\"/archived-movies\" method=\"get\" id=\"views-exposed-form-archived-movies-page\" accept-charset=\"UTF-8\"><div><div class=\"views-exposed-form\">\n  <div class=\"views-exposed-widgets clearfix\">\n          <div id=\"edit-field-us-release-date-value-wrapper\" class=\"views-exposed-widget views-widget-filter-field_us_release_date_value\">\n                  <label for=\"edit-field-us-release-date-value\">\n            Release Date          </label>\n                        <div class=\"views-widget\">\n          <div id=\"edit-field-us-release-date-value-value-wrapper\"><div id=\"edit-field-us-release-date-value-value-inside-wrapper\"><div  class=\"container-inline-date\"><div class=\"form-type-date-select form-item-field-us-release-date-value-value form-item form-group\">\n  <div id=\"edit-field-us-release-date-value-value\"  class=\"date-padding clearfix\"><div class=\"form-type-select form-item-field-us-release-date-value-value-year form-item form-group\">\n  <label class=\"element-invisible\" for=\"edit-field-us-release-date-value-value-year\">Year </label>\n <div class=\"date-year\"><select class=\"date-year form-control form-select\" id=\"edit-field-us-release-date-value-value-year\" name=\"field_us_release_date_value[value][year]\"><option value=\"\">-Year</option><option value=\"2019\" selected=\"selected\">2019</option><option value=\"2018\">2018</option><option value=\"2017\">2017</option><option value=\"2016\">2016</option><option value=\"2015\">2015</option><option value=\"2014\">2014</option><option value=\"2013\">2013</option><option value=\"2012\">2012</option><option value=\"2011\">2011</option><option value=\"2010\">2010</option><option value=\"2009\">2009</option><option value=\"2008\">2008</option><option value=\"2007\">2007</option><option value=\"2006\">2006</option><option value=\"2005\">2005</option><option value=\"2004\">2004</option><option value=\"2003\">2003</option><option value=\"2002\">2002</option><option value=\"2001\">2001</option><option value=\"2000\">2000</option><option value=\"1999\">1999</option><option value=\"1998\">1998</option><option value=\"1997\">1997</option><option value=\"1996\">1996</option><option value=\"1995\">1995</option><option value=\"1994\">1994</option><option value=\"1993\">1993</option><option value=\"1992\">1992</option><option value=\"1991\">1991</option><option value=\"1990\">1990</option><option value=\"1989\">1989</option><option value=\"1988\">1988</option><option value=\"1987\">1987</option><option value=\"1986\">1986</option><option value=\"1985\">1985</option><option value=\"1984\">1984</option><option value=\"1983\">1983</option><option value=\"1982\">1982</option><option value=\"1981\">1981</option><option value=\"1980\">1980</option><option value=\"1979\">1979</option><option value=\"1978\">1978</option><option value=\"1977\">1977</option><option value=\"1976\">1976</option><option value=\"1975\">1975</option><option value=\"1974\">1974</option><option value=\"1973\">1973</option><option value=\"1972\">1972</option><option value=\"1971\">1971</option><option value=\"1970\">1970</option></select></div>\n</div>\n</div>\n</div>\n</div></div></div>        </div>\n              </div>\n                    <div class=\"views-exposed-widget views-submit-button\">\n      <button class=\"ctools-use-ajax ctools-auto-submit-click js-hide btn btn-info form-submit\" id=\"edit-submit-archived-movies\" name=\"\" value=\"Apply\" type=\"submit\">Apply</button>\n    </div>\n      </div>\n</div>\n</div></form>    </div>\n  \n  \n      <div class=\"view-content\">\n        <div class=\"views-row views-row-1 views-row-odd views-row-first views-row-last\">\n      \n  <div class=\"views-field field views-field-field-block-image field-name-field-block-image\">        <div class=\"field-content\"><img typeof=\"foaf:Image\" src=\"https://6a25bbd04bd33b8a843e-9626a8b6c7858057941524bfdad5f5b0.ssl.cf5.rackcdn.com/styles/square/rcf/movies/block_image/5fe5233a03c4eebb111f8a3e5ad86161_0.jpg?itok=l5Rq5Wox\" width=\"400\" height=\"400\" alt=\"\" /></div>  </div>  \n          \n<div class=\"group-description views-fieldset\" data-module=\"views_fieldsets\">\n\n      <div class=\"views-field field views-field-title-field field-name-title-field\"><div class=\"field-content\"><a href=\"/movies/blade-runner-final-cut\">Blade Runner: The Final Cut</a></div></div>  \n</div>\n\n    </div>\n    </div>\n  \n  \n  \n  \n  \n  \n</div>"
	rd := strings.NewReader(text)
	doc, err := goquery.NewDocumentFromReader(rd)
	if err != nil {
		fmt.Println(err)
	} else {
		content := doc.Find(".view-archived-movies .view-content .views-row")
		content.EachWithBreak(func(i int, selection *goquery.Selection) bool {
			res := selection.Find(".views-field")
			res.EachWithBreak(func(j int, imgSelection *goquery.Selection) bool {
				if j == 0 {
					rest, exist := imgSelection.Find(".field-content img").Attr("src")
					if !exist {
						return false
					}
					fmt.Println(rest)
					return true
				} else if j == 1 {
					obj := imgSelection.Find(".field-content a")
					url, exist := obj.Attr("href")
					if exist {
						fmt.Println(url)
					}
					rest, err := obj.Html()
					if err != nil {
						return false
					}
					fmt.Println(rest)
					return true
				} else {
					return false
				}
			})
			return true
		})
	}
}
