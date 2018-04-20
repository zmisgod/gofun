# Picture Spider

target url: http://www.wnlfl.net/fuliba

expact: save images

status: finish ✅

## usage

```
git clone
cd /your/path/goSpider/pictureSpider/
./pictureSpider -start 1 -length 2
```
just wait and do not turn your terminal off

you also can filter the keywords which you dislike

change the dislike keywords in `.env` ,you must use `,` to splite the different keywords in `.env`
```
dislike=秀人网,动漫
```

you can also add `favourite` to download only contain these keywords in `.env`, it will ignore `dislike` keywords
```
favourite=你好
```

## screenshot

![screenshot](https://github.com/zmisgod/goSpider/blob/master/demo/run.png)