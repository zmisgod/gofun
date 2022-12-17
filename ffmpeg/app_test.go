package ffmpeg

import (
	"log"
	"testing"
)

func TestDoTranscode(t *testing.T) {
	_files := []string{
		"https://cdn.poizon.com/du_app/2020/video/222341803_byte5570027_dur0_04e0fa415de1bd39e16dfe3b7085ddb8_1608103378948_du_android_w1088h1920.mp4",
	}
	if err := NewTranscode(_files,
		SetAutoSize(true),
		SetInputFolder("ossIn/"),
		SetTranscodeType(TranscodeType720P265),
		SetOutputFolder("ossOut/"),
		SetOutputFilePrefix("h265_"),
	).DoTranscode(); err != nil {
		log.Fatalln(err)
	}
}
