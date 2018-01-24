# Picture Spider

target url: http://www.alfuli.com/fuliba

expact: save images

status: finish ✅

## usage

```
git clone
cd /your/path/goSpider/pictureSpider/
./pictureSpider http://www.alfuli.com/fuliba 10
```

./pictureSpider [targetURL] [total pages]
just wait and do not turn your terminal off

you also can filter the keywords which you dislike

then find `pictureSpider.go` replace which keywords you dislike
```
var dislikeKeyWord = []string{"漫画", "美女", "xxx"}
```
then run ```go run pictureSpider.go http://www.alfuli.com/fuliba 10```

<hr />

## screenshot

<img src="https://github.com/zmisgod/goSpider/blob/master/demo/run.png">