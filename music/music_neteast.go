package music

import (
	"context"
	"errors"
	"fmt"
	"regexp"
)

type NetEastMusic struct{}

var errorNetEastFetchError = errors.New("根据网易云获取音乐失败")

func (a *NetEastMusic) Fetch(ctx context.Context, linkInfo *LinkInfo) (string, error) {
	var re = regexp.MustCompile(`[?|&]id=([0-9]*)`)
	list := re.FindAllStringSubmatch(linkInfo.Url, 1000)
	idString := ""
	if len(list) > 0 {
		if len(list[0]) >= 2 {
			idString = list[0][1]
		}
	}
	if idString == "" {
		return "", errorNetEastFetchError
	}
	return fmt.Sprintf("http://music.163.com/song/media/outer/url?id=%s.mp3", idString), nil
}
