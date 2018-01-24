# Movie Spider

target url:https://movie.douban.com/j/search_subjects?type=movie&tag=%E8%B1%86%E7%93%A3%E9%AB%98%E5%88%86&sort=rank&page_limit=20&page_start=20

expact: save these data into database

status: finish ✅

## usage

```
git clone
cd /your/path/goSpider/movieSpider/
```

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

<img src="https://github.com/zmisgod/goSpider/blob/master/demo/douban_movie.png">