# Movie Spider

target url : douban's recommand movies API ([选电影](https://movie.douban.com/explore))

expact: save these data into database

status: finish ✅ 

## usage

```
git clone
cd /your/path/goSpider/movieSpider/
```

create a mysql table
```
CREATE TABLE `douban_year_best_movie` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `rate` varchar(20) DEFAULT NULL COMMENT '评分',
  `cover` varchar(200) DEFAULT NULL,
  `title` varchar(200) DEFAULT NULL COMMENT '名称',
  `url` varchar(255) DEFAULT NULL COMMENT '豆瓣地址',
  `playable` tinyint(1) unsigned DEFAULT NULL COMMENT '是否可播放',
  `cover_x` varchar(20) DEFAULT NULL COMMENT '?',
  `cover_y` varchar(20) DEFAULT NULL COMMENT '?',
  `is_new` tinyint(1) unsigned DEFAULT NULL COMMENT '?',
  `year` smallint(4) unsigned DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `year` (`year`),
  KEY `rate` (`rate`,`id`)
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

![screenshot](https://github.com/zmisgod/gofun/blob/master/demo/douban_movie.png)