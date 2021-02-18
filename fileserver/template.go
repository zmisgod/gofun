package fileserver

var listTemplate = `<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<meta http-equiv="X-UA-Compatible" content="ie=edge">
		<title>{{$.Title}}</title>
		<script src="https://cdn.bootcss.com/jquery/3.3.1/jquery.js"></script>
		<style>
			a:visited {
				color: black
			}
	
			a {
				color: black
			}
	
			a.file {
				text-decoration: none;
			}
	
			a.dir {
				text-decoration: underline
			}
		</style>
	</head>
	<body>
		<a href="/">home</a>
		<a href="javascript:" id="goback">goback</a>
		<ul>
			{{range .File}}
			<a class="file" href="/static/{{$.URI}}/{{.}}">
				<li>{{.}}</li>
			</a>
			{{ end }} 
			{{range .Dir}}
			<a class="dir" href="{{$.URI}}/{{.}}">
				<li>{{.}}</li>
			</a>
			{{ end }}
		</ul>
		<div>
			{{$b := 2}} 
			{{if ge $.Page $b }}
			<a href="?page={{$.Ppage}}&pageSize={{$.PageSize}}">上一页</a>
			{{ end }}

			{{if $.ShowNextPage}}
			<a href="?page={{$.Npage}}&pageSize={{$.PageSize}}">下一页</a>
			{{ end }}
		</div>
	</body>
	<script>
		$(document).ready(function () {
			$("#goback").click(function () {
				history.back()
			})
		})
	</script>
	</html>`

var notfoundTemplate = `<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<meta http-equiv="X-UA-Compatible" content="ie=edge">
		<title>{{$.Title}}</title>
	</head>
	<body>
		<h1>This File Can't Be Opened</h1>
		<p>Please check your valid extension</p>
	</body>
	</html>`
