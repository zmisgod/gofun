# goSpider

learn to write a spider with golang

<hr />

## picture spider

target url: http://www.alfuli.com/fuliba

expact: save images

status: finish ✅

### useage

```
git clone
cd /your/path/pictureSpider
./pictureSpider http://www.alfuli.com/fuliba 10
```

./pictureSpider [targetURL] [total pages]
just wait and do not turn your terminal off

you also can filter the keywords which you dislike

then find `pictureSpider.go` replace which keywords you dislike
```
var expactKeyWord = []string{"漫画", "美女", "xxx"}
```
then run ```go run pictureSpider.go http://www.alfuli.com/fuliba 10```

<hr />

## movie spider

target url:https://movie.douban.com/j/search_subjects?type=movie&tag=%E8%B1%86%E7%93%A3%E9%AB%98%E5%88%86&sort=rank&page_limit=20&page_start=20

expact: save these data into database

status: finish ✅

### useage

create a mysql table
```
CREATE TABLE `douban_movie` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `rate` varchar(20) DEFAULT NULL COMMENT '评分',
  `cover` varchar(200) DEFAULT NULL,
  `title` varchar(200) DEFAULT NULL COMMENT '名称',
  `url` varchar(255) DEFAULT NULL COMMENT '豆瓣地址',
  `playable` tinyint(1) unsigned DEFAULT NULL COMMENT '是否可播放',
  `cover_x` varchar(20) DEFAULT NULL COMMENT '?',
  `cover_y` varchar(20) DEFAULT NULL COMMENT '?',
  `is_new` tinyint(1) unsigned DEFAULT NULL COMMENT '?',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=27658397 DEFAULT CHARSET=utf8;
```
then find `doubanSpider.go` replace database config to yours

```
var dbUser = "root"
var dbPass = "111111"
var dbTable = "mytest"
```

then run ```go run doubanSpider.go```

## screenshot

<img src="https://github.com/zmisgod/goSpider/blob/master/demo/run.png">

<img src="https://github.com/zmisgod/goSpider/blob/master/demo/douban_movie.png">

## Postscript

<a href="https://zmis.me/detail_1291">写go爬虫后记有感</a>

## ChangeLog

- 2018-01-17 pictureSpider.go support dislike keywords
- 2018-01-17 doubanSpider.go all done
- 2018-01-01 start douban spider task
- 2017-12-29 pcitureSpider.go all done
- 2017-12-24 start this project

## contact

<a href="https://zmis.me">zmis.me新博客</a>

<a href="https://weibo.com/zmisgod">@zmisgod</a>