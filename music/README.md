## music
a qq-music and net-east music fetch tool

### music fetch usage
1. you need start a node server to generate a qq-music parameters signature string
```
cd music/node
npm install
node srver.js
```

2. then run test file
```
go test -v music/music_test.go music/music.go music/music_neteast.go music/music_qq.go
```