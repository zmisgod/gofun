package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

var url = "http://www.alfuli.com/fuliba/65594.html/3"

func main() {
	innerl := `
	<!DOCTYPE HTML>
	<html>
	<head>
	
	<meta charset="UTF-8">
	<meta http-equiv="X-UA-Compatible" content="IE=11,IE=10,IE=9,IE=8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0, user-scalable=0, minimum-scale=1.0, maximum-scale=1.0">
	<meta http-equiv="Cache-Control" content="no-transform">
	<meta http-equiv="Cache-Control" content="no-siteapp">
	<title>极品网红柚木写真原版系列 32期SM捆綁无圣光套图 [41P]-lu福利吧</title>
	<link rel='dns-prefetch' href='//apps.bdimg.com' />
	<link rel='dns-prefetch' href='//s.w.org' />
	<link rel="alternate" type="application/rss+xml" title="lu福利吧 &raquo; 极品网红柚木写真原版系列 32期SM捆綁无圣光套图 [41P]评论Feed" href="http://www.alfuli.com/fuliba/65594.html/feed" />
	<link rel='stylesheet' id='main-css'  href='http://www.alfuli.com/wp-content/themes/xiu/style.css?ver=5.2' type='text/css' media='all' />
	<script type='text/javascript' src='https://apps.bdimg.com/libs/jquery/2.0.0/jquery.min.js?ver=5.2'></script>
	<link rel='https://api.w.org/' href='http://www.alfuli.com/wp-json/' />
	<link rel="EditURI" type="application/rsd+xml" title="RSD" href="http://www.alfuli.com/xmlrpc.php?rsd" />
	<link rel="wlwmanifest" type="application/wlwmanifest+xml" href="http://www.alfuli.com/wp-includes/wlwmanifest.xml" /> 
	<link rel='prev' title='色列漫画：女主角小葵的故事 不知火舞h同人' href='http://www.alfuli.com/fuliba/62092.html' />
	<link rel="canonical" href="http://www.alfuli.com/fuliba/65594.html/3" />
	<link rel='shortlink' href='http://www.alfuli.com/?p=65594' />
	<link rel="alternate" type="application/json+oembed" href="http://www.alfuli.com/wp-json/oembed/1.0/embed?url=http%3A%2F%2Fwww.alfuli.com%2Ffuliba%2F65594.html" />
	<link rel="alternate" type="text/xml+oembed" href="http://www.alfuli.com/wp-json/oembed/1.0/embed?url=http%3A%2F%2Fwww.alfuli.com%2Ffuliba%2F65594.html&#038;format=xml" />
	<style>.tmall_pc_left {
		position: fixed;
		left: 50%;
		margin-left: -620px;
		bottom: 50px;
		z-index: 1000;
	}</style><meta name="keywords" content="lu福利吧, SM, 写真, 套图, 好福利, 捆绑, 无圣光, 极品网红柚木, 福利吧, 美女图片">
	<meta name="description" content="极品网红柚木写真原版系列 32期SM捆綁无圣光套图 [41P]  抱歉阿，老铁们最近的资源都没啥营养阿，我们就发一篇好福利吧，柚木的保存了好久一直没发上来，分享">
	<link rel="shortcut icon" href="http://www.alfuli.com/wp-content/uploads/2016/12/favicon.ico">
	<!--[if lt IE 9]><script src="http://www.alfuli.com/wp-content/themes/xiu/js/html5.js"></script><![endif]-->
	</head>
	<body class="post-template-default single single-post postid-65594 single-format-standard paged-3 single-paged-3 search_not ui-c3">
	<section class="container">
	<header class="header">
		<div class="logo"><a href="http://www.alfuli.com" title="lu福利吧-万能的福利吧"><img src="http://www.alfuli.com/wp-content/uploads/2016/12/logo.png">lu福利吧</a></div>	<ul class="nav"><li class="navmore"></li><li id="menu-item-19" class="menu-item menu-item-type-taxonomy menu-item-object-category current-post-ancestor current-menu-parent current-post-parent menu-item-has-children menu-item-19"><a href="http://www.alfuli.com/fuliba"><span class="glyphicon glyphicon-home"></span>福利吧</a>
	<ul class="sub-menu">
		<li id="menu-item-28" class="menu-item menu-item-type-taxonomy menu-item-object-category menu-item-28"><a href="http://www.alfuli.com/fuliba/ppp">啪啪啪邪恶动态图</a></li>
		<li id="menu-item-514" class="menu-item menu-item-type-taxonomy menu-item-object-category menu-item-514"><a href="http://www.alfuli.com/fuliba/xemh">邪恶漫画少女漫画</a></li>
		<li id="menu-item-31" class="menu-item menu-item-type-taxonomy menu-item-object-category current-post-ancestor current-menu-parent current-post-parent menu-item-31"><a href="http://www.alfuli.com/fuliba/mntp">美女图片</a></li>
		<li id="menu-item-30" class="menu-item menu-item-type-taxonomy menu-item-object-category menu-item-30"><a href="http://www.alfuli.com/fuliba/tbr">汤不热</a></li>
		<li id="menu-item-29" class="menu-item menu-item-type-taxonomy menu-item-object-category menu-item-29"><a href="http://www.alfuli.com/fuliba/wpfl">微拍福利</a></li>
		<li id="menu-item-2430" class="menu-item menu-item-type-taxonomy menu-item-object-category menu-item-2430"><a href="http://www.alfuli.com/fuliba/fhdq">番号大全</a></li>
	</ul>
	</li>
	<li id="menu-item-18" class="menu-item menu-item-type-taxonomy menu-item-object-category menu-item-18"><a href="http://www.alfuli.com/ytt"><span class="glyphicon glyphicon-time"></span>有头条</a></li>
	<li id="menu-item-16" class="menu-item menu-item-type-taxonomy menu-item-object-category menu-item-16"><a href="http://www.alfuli.com/tcd"><span class="glyphicon glyphicon-edit"></span>吐槽点</a></li>
	<li id="menu-item-17" class="menu-item menu-item-type-taxonomy menu-item-object-category menu-item-17"><a href="http://www.alfuli.com/xjs"><span class="glyphicon glyphicon-user"></span>性教授</a></li>
	<li id="menu-item-57258" class="menu-item menu-item-type-taxonomy menu-item-object-category menu-item-57258"><a href="http://www.alfuli.com/zztj"><span class="glyphicon glyphicon-thumbs-up"></span>ZZ推荐</a></li>
	<li id="menu-item-511" class="menu-item menu-item-type-post_type menu-item-object-page menu-item-511"><a href="http://www.alfuli.com/sitemap"><span class="glyphicon glyphicon-list"></span>网站地图</a></li>
	<li id="menu-item-47217" class="menu-item menu-item-type-post_type menu-item-object-page menu-item-47217"><a href="http://www.alfuli.com/fldh"><span class="glyphicon glyphicon-picture"></span>福利导航</a></li>
	</ul>	<form method="get" class="search-form" action="http://www.alfuli.com/" ><input class="form-control" name="s" type="text" placeholder="输入关键字" value=""><input class="btn" type="submit" value="搜索"></form>	<span class="glyphicon glyphicon-search m-search"></span>	<div class="feeds">
						<a class="feed feed-rss" rel="external nofollow" href="http://www.bh-bj.com/feed" target="_blank"><i></i>RSS订阅</a>
				</div>
		<div class="slinks">
			<a href="http://www.alfuli.com/" title="lu福利吧-万能的福利吧">lu福利吧</a> 
	<a href="http://www.alfuli.com/tags" title="福利吧_标签云集">标签云集</a> 
	
		</div>
	
		</header>
	<div class="content-wrap">
		<div class="content">
					<header class="article-header">
							<div class="breadcrumbs"><span class="text-muted">当前位置：</span><a href="http://www.alfuli.com">lu福利吧</a> <small>></small> <a href="http://www.alfuli.com/fuliba">福利吧</a> <small>></small> <span class="text-muted">正文</span></div>
							<h1 class="article-title"><a href="http://www.alfuli.com/fuliba/65594.html">极品网红柚木写真原版系列 32期SM捆綁无圣光套图 [41P]</a></h1>
				<ul class="article-meta">
									<li><a href="http://www.alfuli.com/author/admin">好多福利</a> 发布于 2017-12-27</li>
					<li>分类：<a href="http://www.alfuli.com/fuliba" rel="category tag">福利吧</a> / <a href="http://www.alfuli.com/fuliba/mntp" rel="category tag">美女图片</a></li>
									<li><span class="post-views">阅读(10335)</span></li>
					<li>评论(0)</li>
					<li></li>
				</ul>
			</header>
			<div class="ads ads-content ads-post"><!--<a href="http://www.selanglu.info" title="你想看的" target="_blank"><img src="http://www.bh-bj.com/wp-content/uploads/2017/03/top.gif" alt="你想看的"></a>-->
	<a href="http://bbs.fxfuli.net" target="_blank"><img src="http://www.alfuli.com/wp-content/uploads/2017/12/lt.jpg" alt="vip论坛"></a><br/>
	<a href="https://wvwv.aycxdq.com/fxfuli.com.php" target="_blank"><img src="http://www.alfuli.com/wp-content/uploads/2017/09/720x70.gif"></a>
	
	</div>		<article class="article-content">
				<p> <a href="http://www.alfuli.com/wp-content/uploads/2017/12/DSC_3870.jpg"><a href="http://www.alfuli.com/fuliba/65594.html/4" title="点击图片查看下一张" ><img class="aligncenter size-full wp-image-65597" src="http://www.alfuli.com/wp-content/uploads/2017/12/DSC_3870.jpg" alt="极品网红柚木写真原版系列 32期SM捆綁无圣光套图 [41P]" ></a></a></p>
	<p class="post-copyright">未经允许不得转载：<a href="http://www.alfuli.com">lu福利吧</a> &raquo; <a href="http://www.alfuli.com/fuliba/65594.html">极品网红柚木写真原版系列 32期SM捆綁无圣光套图 [41P]</a></p>		</article>
			<div class="article-paging"> <a href="http://www.alfuli.com/fuliba/65594.html"><span>1</span></a> <a href="http://www.alfuli.com/fuliba/65594.html/2"><span>2</span></a> <span>3</span> <a href="http://www.alfuli.com/fuliba/65594.html/4"><span>4</span></a> <a href="http://www.alfuli.com/fuliba/65594.html/5"><span>5</span></a> <a href="http://www.alfuli.com/fuliba/65594.html/6"><span>6</span></a> <a href="http://www.alfuli.com/fuliba/65594.html/7"><span>7</span></a> <a href="http://www.alfuli.com/fuliba/65594.html/8"><span>8</span></a> <a href="http://www.alfuli.com/fuliba/65594.html/9"><span>9</span></a> <a href="http://www.alfuli.com/fuliba/65594.html/10"><span>10</span></a> <a href="http://www.alfuli.com/fuliba/65594.html/11"><span>11</span></a> <a href="http://www.alfuli.com/fuliba/65594.html/12"><span>12</span></a> <a href="http://www.alfuli.com/fuliba/65594.html/13"><span>13</span></a> <a href="http://www.alfuli.com/fuliba/65594.html/14"><span>14</span></a> <a href="http://www.alfuli.com/fuliba/65594.html/15"><span>15</span></a> <a href="http://www.alfuli.com/fuliba/65594.html/16"><span>16</span></a> <a href="http://www.alfuli.com/fuliba/65594.html/17"><span>17</span></a> <a href="http://www.alfuli.com/fuliba/65594.html/18"><span>18</span></a> <a href="http://www.alfuli.com/fuliba/65594.html/19"><span>19</span></a> <a href="http://www.alfuli.com/fuliba/65594.html/20"><span>20</span></a> <a href="http://www.alfuli.com/fuliba/65594.html/21"><span>21</span></a> <a href="http://www.alfuli.com/fuliba/65594.html/22"><span>22</span></a> <a href="http://www.alfuli.com/fuliba/65594.html/23"><span>23</span></a> <a href="http://www.alfuli.com/fuliba/65594.html/24"><span>24</span></a> <a href="http://www.alfuli.com/fuliba/65594.html/25"><span>25</span></a> <a href="http://www.alfuli.com/fuliba/65594.html/26"><span>26</span></a> <a href="http://www.alfuli.com/fuliba/65594.html/27"><span>27</span></a> <a href="http://www.alfuli.com/fuliba/65594.html/28"><span>28</span></a> <a href="http://www.alfuli.com/fuliba/65594.html/29"><span>29</span></a> <a href="http://www.alfuli.com/fuliba/65594.html/30"><span>30</span></a> <a href="http://www.alfuli.com/fuliba/65594.html/31"><span>31</span></a> <a href="http://www.alfuli.com/fuliba/65594.html/32"><span>32</span></a> <a href="http://www.alfuli.com/fuliba/65594.html/33"><span>33</span></a> <a href="http://www.alfuli.com/fuliba/65594.html/34"><span>34</span></a> <a href="http://www.alfuli.com/fuliba/65594.html/35"><span>35</span></a> <a href="http://www.alfuli.com/fuliba/65594.html/36"><span>36</span></a> <a href="http://www.alfuli.com/fuliba/65594.html/37"><span>37</span></a> <a href="http://www.alfuli.com/fuliba/65594.html/38"><span>38</span></a> <a href="http://www.alfuli.com/fuliba/65594.html/39"><span>39</span></a> <a href="http://www.alfuli.com/fuliba/65594.html/40"><span>40</span></a> <a href="http://www.alfuli.com/fuliba/65594.html/41"><span>41</span></a> <a href="http://www.alfuli.com/fuliba/65594.html/42"><span>42</span></a></div>				<div class="article-social">
				<a href="javascript:;" class="action action-like" data-pid="65594" data-event="like"><i class="glyphicon glyphicon-thumbs-up"></i>赞 (<span>6</span>)</a>			
	
			</div>
	
			<div class="action-share bdsharebuttonbox">
				分享到：<a class="bds_qzone" data-cmd="qzone"></a><a class="bds_tsina" data-cmd="tsina"></a><a class="bds_weixin" data-cmd="weixin"></a><a class="bds_tqq" data-cmd="tqq"></a><a class="bds_sqq" data-cmd="sqq"></a><a class="bds_bdhome" data-cmd="bdhome"></a><a class="bds_tqf" data-cmd="tqf"></a><a class="bds_renren" data-cmd="renren"></a><a class="bds_diandian" data-cmd="diandian"></a><a class="bds_youdao" data-cmd="youdao"></a><a class="bds_ty" data-cmd="ty"></a><a class="bds_kaixin001" data-cmd="kaixin001"></a><a class="bds_taobao" data-cmd="taobao"></a><a class="bds_douban" data-cmd="douban"></a><a class="bds_fbook" data-cmd="fbook"></a><a class="bds_twi" data-cmd="twi"></a><a class="bds_mail" data-cmd="mail"></a><a class="bds_copy" data-cmd="copy"></a><a class="bds_more" data-cmd="more">更多</a> (<a class="bds_count" data-cmd="count"></a>)		</div>
			<div class="article-tags">
				标签：<a href="http://www.alfuli.com/tag/lu%e7%a6%8f%e5%88%a9%e5%90%a7" rel="tag">lu福利吧</a><a href="http://www.alfuli.com/tag/sm" rel="tag">SM</a><a href="http://www.alfuli.com/tag/%e5%86%99%e7%9c%9f" rel="tag">写真</a><a href="http://www.alfuli.com/tag/%e5%a5%97%e5%9b%be" rel="tag">套图</a><a href="http://www.alfuli.com/tag/%e5%a5%bd%e7%a6%8f%e5%88%a9" rel="tag">好福利</a><a href="http://www.alfuli.com/tag/%e6%8d%86%e7%bb%91" rel="tag">捆绑</a><a href="http://www.alfuli.com/tag/%e6%97%a0%e5%9c%a3%e5%85%89" rel="tag">无圣光</a><a href="http://www.alfuli.com/tag/%e6%9e%81%e5%93%81%e7%bd%91%e7%ba%a2%e6%9f%9a%e6%9c%a8" rel="tag">极品网红柚木</a>		</div>
			<nav class="article-nav">
				<span class="article-nav-prev">上一篇<br><a href="http://www.alfuli.com/fuliba/62092.html" rel="prev">色列漫画：女主角小葵的故事 不知火舞h同人</a></span>
				<span class="article-nav-next"></span>
			</nav>
			<div class="ads ads-content ads-related">亲们点击图片即可翻到下一页额~~
	</div>		<div class="relates relates-model-thumb"><h3 class="title"><strong>相关推荐</strong></h3><ul><li><a target="_blank" href="http://www.alfuli.com/fuliba/60977.html"><span><img data-original="http://www.alfuli.com/wp-content/uploads/2017/12/1442545218-60.jpg" class="thumb"/></span>[XiuRen]秀人网第527期女神琳琳ailin写真[62P]</a></li><li><a target="_blank" href="http://www.alfuli.com/fuliba/60976.html"><span><img data-original="http://www.alfuli.com/wp-content/uploads/2017/12/140F52b5-31.jpg" class="thumb"/></span>[XiuRen]秀人网第526期女神妹子谭睿琪Ailsa[40P]</a></li><li><a target="_blank" href="http://www.alfuli.com/fuliba/65528.html"><span><img data-original="http://www.alfuli.com/wp-content/uploads/2017/12/011-4.jpg" class="thumb"/></span>微博网红软妹子@一只肉酱阿之湿身死库水福利图包 [21P]</a></li><li><a target="_blank" href="http://www.alfuli.com/fuliba/60975.html"><span><img data-original="http://www.alfuli.com/wp-content/uploads/2017/12/141Q25X7-0.jpg" class="thumb"/></span>[XiuRen]秀人网第525期谢芷馨Sindy高清写真[50P]</a></li><li><a target="_blank" href="http://www.alfuli.com/fuliba/61113.html"><span><img data-original="http://www.alfuli.com/wp-content/uploads/2017/12/0033.jpg" class="thumb"/></span>PR社@原来是茜公举殿下 &#8211; 户外温泉87福利图包 [34P4V]</a></li><li><a target="_blank" href="http://www.alfuli.com/fuliba/60973.html"><span><img data-original="http://www.alfuli.com/wp-content/uploads/2017/12/1450034300-0.jpg" class="thumb"/></span>[XiuRen]秀人网车模梦心玥私拍无圣光套图[46P]</a></li><li><a target="_blank" href="http://www.alfuli.com/fuliba/60972.html"><span><img data-original="http://www.alfuli.com/wp-content/uploads/2017/12/14512142R-49.jpg" class="thumb"/></span>[XiuRen]秀人网第524期Elise谭晓彤内衣秀[50P]</a></li><li><a target="_blank" href="http://www.alfuli.com/fuliba/61034.html"><span><img data-original="http://www.alfuli.com/wp-content/uploads/2017/12/1-39-4.jpg" class="thumb"/></span>极品网红少女私人玩物—蓝白条纹袜福利图包 [63P6V]</a></li></ul></div>		<div class="sticky"><h3 class="title"><strong>热门推荐</strong></h3><ul><li class="item"><a target="_blank" href="http://www.alfuli.com/zztj/61024.html"><span><img data-original="http://www.alfuli.com/wp-content/uploads/2017/12/a15b4afegy1fipg71s2mmj209y0hoq3b.jpg" class="thumb"/></span>很多女神小姐姐在羞羞</a></li></ul></div>		<div class="ads ads-content ads-comment"></div>		<h3 class="title" id="comments">
		<div class="text-muted pull-right">
			</div>
		<strong>评论 <b> 0 </b></strong>
	</h3>
	<div id="respond" class="no_webshot">
			
		<form action="http://www.alfuli.com/wp-comments-post.php" method="post" id="commentform">
			
			<div class="comt-title">
				<div class="comt-avatar">
					<img alt='' data-original='http://2.gravatar.com/avatar/?s=50&#038;d=http%3A%2F%2Fwww.alfuli.com%2Fwp-content%2Fthemes%2Fxiu%2Fimages%2Favatar-default.png&#038;r=g' srcset='http://0.gravatar.com/avatar/?s=100&#038;d=http%3A%2F%2Fwww.alfuli.com%2Fwp-content%2Fthemes%2Fxiu%2Fimages%2Favatar-default.png&#038;r=g 2x' class='avatar avatar-50 photo avatar-default' height='50' width='50' />			</div>
				<div class="comt-author">
							</div>
				<a id="cancel-comment-reply-link" href="javascript:;">取消</a>
			</div>
			
			<div class="comt">
				<div class="comt-box">
					<textarea placeholder="听说评论的亲们，都会变成雕大的人呀！" class="input-block-level comt-area" name="comment" id="comment" cols="100%" rows="3" tabindex="1" onkeydown="if(event.ctrlKey&amp;&amp;event.keyCode==13){document.getElementById('submit').click();return false};"></textarea>
					<div class="comt-ctrl">
						<div class="comt-tips"><input type='hidden' name='comment_post_ID' value='65594' id='comment_post_ID' />
	<input type='hidden' name='comment_parent' id='comment_parent' value='0' />
	<label for="comment_mail_notify" class="checkbox inline hide" style="padding-top:0"><input type="checkbox" name="comment_mail_notify" id="comment_mail_notify" value="comment_mail_notify" checked="checked"/>有人回复时邮件通知我</label></div>
						<button type="submit" name="submit" id="submit" tabindex="5"><i class="icon-ok-circle icon-white icon12"></i> 来一发</button>
						<!-- <span data-type="comment-insert-smilie" class="muted comt-smilie"><i class="icon-thumbs-up icon12"></i> 表情</span> -->
					</div>
				</div>
	
													<div class="comt-comterinfo" id="comment-author-info" >
							<ul>
								<li class="form-inline"><label class="hide" for="author">昵称</label><input class="ipt" type="text" name="author" id="author" value="" tabindex="2" placeholder="昵称"><span class="text-muted">昵称 (必填)</span></li>
								<li class="form-inline"><label class="hide" for="email">邮箱</label><input class="ipt" type="text" name="email" id="email" value="" tabindex="3" placeholder="邮箱"><span class="text-muted">邮箱 (必填)</span></li>
								<!--<li class="form-inline"><label class="hide" for="url">网址</label><input class="ipt" type="text" name="url" id="url" value="" tabindex="4" placeholder="网址"><span class="text-muted">网址</span></li>-->
							</ul>
						</div>
										</div>
	
		</form>
		</div>
		</div>
	</div>
	<aside class="sidebar">	
	<div class="widget widget_ads"><div class="widget_ads_inner"><a href="mailto:sgg99com@126.com"><img src="http://www.alfuli.com/wp-content/uploads/2017/09/ad.jpg" alt="福利吧_ads"></a></div></div><div class="widget widget_ads"><div class="widget_ads_inner"><a href="http://www.puaqs.com" target="_blank"><img src="http://www.alfuli.com/wp-content/uploads/2017/10/pcad.jpg"></a><br/>
	<a href="http://www.alfuli.com/zztj/59535.html" target="_blank"><img src="http://www.alfuli.com/wp-content/uploads/2017/11/alipay.jpg" alt="支付宝福利"></a>
	</div></div><div class="widget widget_searchbox"><h3 class="title"><strong>搜索福利吧</strong></h3><form method="get" class="search-form" action="http://www.alfuli.com/" ><input class="form-control" name="s" type="text" placeholder="听说搜索出来的都是大尺度福利额！" value=""><input class="btn" type="submit" value="搜索"></form></div><div class="widget widget_postlist"><h3 class="title"><strong>一波随机福利</strong></h3><ul class="items-01">		<li><a target="_blank" href="http://www.alfuli.com/ytt/53861.html"><span class="thumbnail"><span><img data-original="http://www.alfuli.com/wp-content/uploads/2017/10/0655.jpg" class="thumb"/></span></span><span class="text">15岁小女孩带闺蜜同行跨省见网友，还好没赶到地哒~</span><span class="text-muted post-views">阅读(10945)</span></a></li>
					<li><a target="_blank" href="http://www.alfuli.com/xjs/50130.html"><span class="thumbnail"><span><img data-original="http://www.alfuli.com/wp-content/uploads/2017/10/10007.jpg" class="thumb"/></span></span><span class="text">嘤嘤嘤~好污额女生第一次之后身体有什么变化？</span><span class="text-muted post-views">阅读(41179)</span></a></li>
					<li><a target="_blank" href="http://www.alfuli.com/zztj/58409.html"><span class="thumbnail"><span><img data-original="http://www.alfuli.com/wp-content/uploads/2017/11/0060.jpg" class="thumb"/></span></span><span class="text">精选Cosplay美女福利图集[128P]</span><span class="text-muted post-views">阅读(100648)</span></a></li>
					<li><a target="_blank" href="http://www.alfuli.com/tcd/31050.html"><span class="thumbnail"><span><img data-original="http://www.alfuli.com/wp-content/uploads/2017/06/c1419h82uwh.m701_20170625214634.jpg" class="thumb"/></span></span><span class="text">社会套路很深!女孩子借贷要谨慎~送给那些懵懂的孩子们</span><span class="text-muted post-views">阅读(28070)</span></a></li>
					<li><a target="_blank" href="http://www.alfuli.com/ytt/54121.html"><span class="thumbnail"><span><img data-original="http://www.alfuli.com/wp-content/uploads/2017/10/FqmrE0xJWFXehjm1RuUgaJsiTEvs.jpg" class="thumb"/></span></span><span class="text">难以置信的高校恋爱须知：恋爱应避免同居~</span><span class="text-muted post-views">阅读(3322)</span></a></li>
					<li><a target="_blank" href="http://www.alfuli.com/tcd/31265.html"><span class="thumbnail"><span><img data-original="http://www.alfuli.com/wp-content/uploads/2017/06/tjl.jpg" class="thumb"/></span></span><span class="text">中国人体艺术摄影@汤加丽,男人看了会脸红心跳的裸模的自述</span><span class="text-muted post-views">阅读(22327)</span></a></li>
					<li><a target="_blank" href="http://www.alfuli.com/tcd/16754.html"><span class="thumbnail"><span><img data-original="http://www.alfuli.com/wp-content/uploads/2017/04/0O97PBYETBVR9686NEJ2.png" class="thumb"/></span></span><span class="text">17岁被母亲骗去拍裸照、陪睡，20岁结婚再离婚，29岁获封日本影后，她的人生才叫开挂！</span><span class="text-muted post-views">阅读(49882)</span></a></li>
					<li><a target="_blank" href="http://www.alfuli.com/zztj/61187.html"><span class="thumbnail"><span><img data-original="http://www.alfuli.com/wp-content/uploads/2017/12/0626_wztp_1.jpg" class="thumb"/></span></span><span class="text">超感人大话西游爱情物语：我们老了，我在来生等你</span><span class="text-muted post-views">阅读(484)</span></a></li>
					<li><a target="_blank" href="http://www.alfuli.com/tcd/28182.html"><span class="thumbnail"><span><img data-original="http://www.alfuli.com/wp-content/uploads/2017/06/xq_20170613065528.jpg" class="thumb"/></span></span><span class="text">当美女相亲被提及是否处女,美女的回答简直是太残暴了~</span><span class="text-muted post-views">阅读(10412)</span></a></li>
					<li><a target="_blank" href="http://www.alfuli.com/zztj/57269.html"><span class="thumbnail"><span><img data-original="http://www.alfuli.com/wp-content/uploads/2017/11/1510107636_tZfgZWxs.jpg" class="thumb"/></span></span><span class="text">一波仪态万千的高颜值美女</span><span class="text-muted post-views">阅读(14164)</span></a></li>
			</ul></div><div class="widget widget_tag_cloud"><h3 class="title"><strong>热门标签</strong></h3><div class="tagcloud"><a href="http://www.alfuli.com/tag/95" class="tag-cloud-link tag-link-550 tag-link-position-1" style="font-size: 10.545454545455pt;" aria-label="95 (44个项目)">95</a>
	<a href="http://www.alfuli.com/tag/2017" class="tag-cloud-link tag-link-1048 tag-link-position-2" style="font-size: 13.727272727273pt;" aria-label="2017 (89个项目)">2017</a>
	<a href="http://www.alfuli.com/tag/graphis" class="tag-cloud-link tag-link-216 tag-link-position-3" style="font-size: 8.8484848484848pt;" aria-label="Graphis (30个项目)">Graphis</a>
	<a href="http://www.alfuli.com/tag/lu%e7%a6%8f%e5%88%a9" class="tag-cloud-link tag-link-14 tag-link-position-4" style="font-size: 11.5pt;" aria-label="lu福利 (54个项目)">lu福利</a>
	<a href="http://www.alfuli.com/tag/lu%e7%a6%8f%e5%88%a9%e5%90%a7" class="tag-cloud-link tag-link-308 tag-link-position-5" style="font-size: 21.681818181818pt;" aria-label="lu福利吧 (499个项目)">lu福利吧</a>
	<a href="http://www.alfuli.com/tag/mistar" class="tag-cloud-link tag-link-581 tag-link-position-6" style="font-size: 9.0606060606061pt;" aria-label="MiStar (31个项目)">MiStar</a>
	<a href="http://www.alfuli.com/tag/tuigirl" class="tag-cloud-link tag-link-136 tag-link-position-7" style="font-size: 12.348484848485pt;" aria-label="TuiGirl (65个项目)">TuiGirl</a>
	<a href="http://www.alfuli.com/tag/vip" class="tag-cloud-link tag-link-178 tag-link-position-8" style="font-size: 9.5909090909091pt;" aria-label="VIP (35个项目)">VIP</a>
	<a href="http://www.alfuli.com/tag/xiuren" class="tag-cloud-link tag-link-272 tag-link-position-9" style="font-size: 15.424242424242pt;" aria-label="xiuren (128个项目)">xiuren</a>
	<a href="http://www.alfuli.com/tag/%e4%ba%ba%e6%b0%94" class="tag-cloud-link tag-link-998 tag-link-position-10" style="font-size: 9.1666666666667pt;" aria-label="人气 (32个项目)">人气</a>
	<a href="http://www.alfuli.com/tag/%e5%86%99%e7%9c%9f" class="tag-cloud-link tag-link-68 tag-link-position-11" style="font-size: 16.166666666667pt;" aria-label="写真 (152个项目)">写真</a>
	<a href="http://www.alfuli.com/tag/%e5%88%b6%e6%9c%8d" class="tag-cloud-link tag-link-251 tag-link-position-12" style="font-size: 9.9090909090909pt;" aria-label="制服 (38个项目)">制服</a>
	<a href="http://www.alfuli.com/tag/%e5%8f%af%e7%88%b1" class="tag-cloud-link tag-link-275 tag-link-position-13" style="font-size: 8.4242424242424pt;" aria-label="可爱 (27个项目)">可爱</a>
	<a href="http://www.alfuli.com/tag/%e5%9b%bd%e6%a8%a1" class="tag-cloud-link tag-link-384 tag-link-position-14" style="font-size: 8.530303030303pt;" aria-label="国模 (28个项目)">国模</a>
	<a href="http://www.alfuli.com/tag/%e5%a4%a7%e5%b0%ba%e5%ba%a6" class="tag-cloud-link tag-link-47 tag-link-position-15" style="font-size: 16.909090909091pt;" aria-label="大尺度 (175个项目)">大尺度</a>
	<a href="http://www.alfuli.com/tag/%e5%a5%97%e5%9b%be" class="tag-cloud-link tag-link-56 tag-link-position-16" style="font-size: 22pt;" aria-label="套图 (541个项目)">套图</a>
	<a href="http://www.alfuli.com/tag/%e5%a5%b3%e7%a5%9e" class="tag-cloud-link tag-link-123 tag-link-position-17" style="font-size: 12.136363636364pt;" aria-label="女神 (62个项目)">女神</a>
	<a href="http://www.alfuli.com/tag/%e5%a5%b3%e9%83%8e" class="tag-cloud-link tag-link-537 tag-link-position-18" style="font-size: 12.348484848485pt;" aria-label="女郎 (65个项目)">女郎</a>
	<a href="http://www.alfuli.com/tag/%e5%a5%bd%e7%a6%8f%e5%88%a9" class="tag-cloud-link tag-link-778 tag-link-position-19" style="font-size: 9.0606060606061pt;" aria-label="好福利 (31个项目)">好福利</a>
	<a href="http://www.alfuli.com/tag/%e5%a6%b9%e5%ad%90" class="tag-cloud-link tag-link-18 tag-link-position-20" style="font-size: 12.772727272727pt;" aria-label="妹子 (72个项目)">妹子</a>
	<a href="http://www.alfuli.com/tag/%e5%ab%a9%e6%a8%a1" class="tag-cloud-link tag-link-182 tag-link-position-21" style="font-size: 8pt;" aria-label="嫩模 (25个项目)">嫩模</a>
	<a href="http://www.alfuli.com/tag/%e5%b0%91%e5%a5%b3" class="tag-cloud-link tag-link-191 tag-link-position-22" style="font-size: 9.5909090909091pt;" aria-label="少女 (35个项目)">少女</a>
	<a href="http://www.alfuli.com/tag/%e5%b7%a8%e4%b9%b3" class="tag-cloud-link tag-link-134 tag-link-position-23" style="font-size: 11.818181818182pt;" aria-label="巨乳 (58个项目)">巨乳</a>
	<a href="http://www.alfuli.com/tag/%e5%be%ae%e5%8d%9a" class="tag-cloud-link tag-link-562 tag-link-position-24" style="font-size: 11.606060606061pt;" aria-label="微博 (55个项目)">微博</a>
	<a href="http://www.alfuli.com/tag/%e6%80%a7%e6%84%9f" class="tag-cloud-link tag-link-59 tag-link-position-25" style="font-size: 17.545454545455pt;" aria-label="性感 (204个项目)">性感</a>
	<a href="http://www.alfuli.com/tag/%e6%8e%a8%e5%a5%b3%e9%83%8e" class="tag-cloud-link tag-link-137 tag-link-position-26" style="font-size: 12.454545454545pt;" aria-label="推女郎 (66个项目)">推女郎</a>
	<a href="http://www.alfuli.com/tag/%e6%92%b8%e4%b8%80%e6%92%b8" class="tag-cloud-link tag-link-218 tag-link-position-27" style="font-size: 9.1666666666667pt;" aria-label="撸一撸 (32个项目)">撸一撸</a>
	<a href="http://www.alfuli.com/tag/%e6%92%b8%e7%82%b9" class="tag-cloud-link tag-link-460 tag-link-position-28" style="font-size: 8.2121212121212pt;" aria-label="撸点 (26个项目)">撸点</a>
	<a href="http://www.alfuli.com/tag/%e6%97%85%e6%8b%8d" class="tag-cloud-link tag-link-645 tag-link-position-29" style="font-size: 8.4242424242424pt;" aria-label="旅拍 (27个项目)">旅拍</a>
	<a href="http://www.alfuli.com/tag/%e6%97%a0%e5%9c%a3%e5%85%89" class="tag-cloud-link tag-link-140 tag-link-position-30" style="font-size: 21.681818181818pt;" aria-label="无圣光 (503个项目)">无圣光</a>
	<a href="http://www.alfuli.com/tag/%e6%97%a0%e5%9c%a3%e5%85%89%e5%a5%97%e5%9b%be" class="tag-cloud-link tag-link-856 tag-link-position-31" style="font-size: 12.560606060606pt;" aria-label="无圣光套图 (68个项目)">无圣光套图</a>
	<a href="http://www.alfuli.com/tag/%e6%9e%81%e5%93%81" class="tag-cloud-link tag-link-157 tag-link-position-32" style="font-size: 10.757575757576pt;" aria-label="极品 (46个项目)">极品</a>
	<a href="http://www.alfuli.com/tag/%e6%a8%a1%e7%89%b9" class="tag-cloud-link tag-link-90 tag-link-position-33" style="font-size: 17.227272727273pt;" aria-label="模特 (191个项目)">模特</a>
	<a href="http://www.alfuli.com/tag/%e7%a6%8f%e5%88%a9" class="tag-cloud-link tag-link-12 tag-link-position-34" style="font-size: 13.515151515152pt;" aria-label="福利 (85个项目)">福利</a>
	<a href="http://www.alfuli.com/tag/%e7%a6%8f%e5%88%a9%e5%90%a7" class="tag-cloud-link tag-link-2 tag-link-position-35" style="font-size: 15.530303030303pt;" aria-label="福利吧 (132个项目)">福利吧</a>
	<a href="http://www.alfuli.com/tag/%e7%a6%8f%e5%88%a9%e5%9b%be" class="tag-cloud-link tag-link-71 tag-link-position-36" style="font-size: 11.075757575758pt;" aria-label="福利图 (49个项目)">福利图</a>
	<a href="http://www.alfuli.com/tag/%e7%a6%8f%e5%88%a9%e8%a7%86%e9%a2%91" class="tag-cloud-link tag-link-26 tag-link-position-37" style="font-size: 10.545454545455pt;" aria-label="福利视频 (44个项目)">福利视频</a>
	<a href="http://www.alfuli.com/tag/%e7%a7%80%e4%ba%ba%e7%bd%91" class="tag-cloud-link tag-link-273 tag-link-position-38" style="font-size: 15.954545454545pt;" aria-label="秀人网 (144个项目)">秀人网</a>
	<a href="http://www.alfuli.com/tag/%e7%a7%81%e6%8b%8d" class="tag-cloud-link tag-link-181 tag-link-position-39" style="font-size: 15.424242424242pt;" aria-label="私拍 (127个项目)">私拍</a>
	<a href="http://www.alfuli.com/tag/%e7%b2%89%e5%ab%a9" class="tag-cloud-link tag-link-156 tag-link-position-40" style="font-size: 10.863636363636pt;" aria-label="粉嫩 (47个项目)">粉嫩</a>
	<a href="http://www.alfuli.com/tag/%e7%be%8e%e5%a5%b3" class="tag-cloud-link tag-link-159 tag-link-position-41" style="font-size: 14.469696969697pt;" aria-label="美女 (103个项目)">美女</a>
	<a href="http://www.alfuli.com/tag/%e8%90%9d%e8%8e%89" class="tag-cloud-link tag-link-54 tag-link-position-42" style="font-size: 9.6969696969697pt;" aria-label="萝莉 (36个项目)">萝莉</a>
	<a href="http://www.alfuli.com/tag/%e8%af%b1%e6%83%91" class="tag-cloud-link tag-link-77 tag-link-position-43" style="font-size: 12.030303030303pt;" aria-label="诱惑 (60个项目)">诱惑</a>
	<a href="http://www.alfuli.com/tag/%e9%82%aa%e6%81%b6%e5%8a%a8%e6%80%81%e5%9b%be" class="tag-cloud-link tag-link-21 tag-link-position-44" style="font-size: 9.5909090909091pt;" aria-label="邪恶动态图 (35个项目)">邪恶动态图</a>
	<a href="http://www.alfuli.com/tag/%e9%ad%85%e5%a6%8d%e7%a4%be" class="tag-cloud-link tag-link-582 tag-link-position-45" style="font-size: 9.2727272727273pt;" aria-label="魅妍社 (33个项目)">魅妍社</a></div>
	</div><div class="widget widget_text"><h3 class="title"><strong>热门网站</strong></h3>			<div class="textwidget"><p><a href="http://znwz.net" title="宅男网 - 宅男网站,宅男吧,宅男福利吧,宅男神器" target="_blank">宅男网</a>  <a href="http://www.5ajw.com" title="爱邪恶_我们都爱邪恶漫画_内涵漫画_邪恶少女漫画_日本、韩国、色小组邪恶漫画图" target="_blank">爱邪恶漫画</a> <a href="http://www.mmxy.net" title="mm福利吧" target="_blank">mm福利吧</a><br />
	<br /><a href="http://www.qqxyz.net" title="好多福利吧" target="_blank">好多福利吧</a>  <a href="http://www.zxflyy.com" title="在线福利影院" target="_blank">在线福利影院</a><br />
	<br />
	<img src="http://www.bh-bj.com/wp-content/uploads/2017/01/fuliba.png" alt="lu福利吧_万能的福利吧"></p>
	</div>
			</div></aside><footer class="footer">
		&copy; 2017 <a href="http://www.alfuli.com">lu福利吧</a> &nbsp;  &nbsp; 京ICP备14014386号-2|<b><a href="http://www.alfuli.com/ad" title="lu福利吧_广告合作">广告合作</a></b>|<strong><a href="http://www.alfuli.com/sitemap.xml" title="lu福利吧_网站地图">百度网站地图</a></strong>
		<div style="display:none"><script type="text/javascript">var cnzz_protocol = (("https:" == document.location.protocol) ? " https://" : " http://");document.write(unescape("%3Cspan id='cnzz_stat_icon_1258494869'%3E%3C/span%3E%3Cscript src='" + cnzz_protocol + "s4.cnzz.com/stat.php%3Fid%3D1258494869%26show%3Dpic1' type='text/javascript'%3E%3C/script%3E"));</script></div></footer>
	</section>
	
	<script>
	window.jui = {
		uri: 'http://www.alfuli.com/wp-content/themes/xiu',
		roll: '1 2 3',
		ajaxpager: '0'
	}
	</script>
	<script type='text/javascript' src='https://apps.bdimg.com/libs/bootstrap/3.2.0/js/bootstrap.min.js?ver=5.2'></script>
	<script type='text/javascript' src='http://www.alfuli.com/wp-content/themes/xiu/js/custom.js?ver=5.2'></script>
	<script type='text/javascript' src='http://www.alfuli.com/wp-includes/js/wp-embed.min.js?ver=4.9'></script>
	
	<script>
	(function(){
		var bp = document.createElement('script');
		var curProtocol = window.location.protocol.split(':')[0];
		if (curProtocol === 'https') {
			bp.src = 'https://zz.bdstatic.com/linksubmit/push.js';        
		}
		else {
			bp.src = 'http://push.zhanzhang.baidu.com/push.js';
		}
		var s = document.getElementsByTagName("script")[0];
		s.parentNode.insertBefore(bp, s);
	})();
	</script>
	
	
	</body>
	</html>`
	rex := regexp.MustCompile(`<div class="article-paging">(.*?)<\/div>`)
	ul := rex.FindString(innerl)
	ulrex := regexp.MustCompile(`([a-zA-z]+:[^\s]*")`)
	lists := ulrex.FindAllString(ul, -1)
	for _, v := range lists {
		fmt.Println(strings.Replace(v, "\"", "", -1))
	}
}

func fetch(url string) {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	innerHTML := string(respBody)

	rexp := regexp.MustCompile(`<img class="aligncenter size-full\s(.*?)" src="(.*?)"`)
	result := rexp.FindString(innerHTML)
	//<img class="aligncenter size-full wp-image-65595" src="http://www.alfuli.com/wp-content/uploads/2017/12/DSC_3816.jpg"

	rexpr := regexp.MustCompile(`([a-zA-z]+:[^\s]*)`)
	img := strings.Replace(string(rexpr.FindString(result)), "\"", "", -1)
	println(img)
}
