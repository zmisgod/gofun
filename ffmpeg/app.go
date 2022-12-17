package ffmpeg

import (
	"bytes"
	"encoding/json"
	"fmt"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

var h2641080Args = ffmpeg.KwArgs{
	"c:v":     "libx264",
	"codec:a": "aac",
	"b:a":     "125k",
	"r":       "30",
	"preset":  "slow",
	"crf":     "26",
}

var h264720Args = ffmpeg.KwArgs{
	"c:v":     "libx264",
	"codec:a": "aac",
	"b:a":     "125k",
	"r":       "30",
	"preset":  "slow",
	"crf":     "26",
}

var h265720Args = ffmpeg.KwArgs{
	"c:v":     "libx265",
	"codec:a": "aac",
	"b:a":     "125k",
	"r":       "30",
	"crf":     "26",
}

var h2651080Args = ffmpeg.KwArgs{
	"c:v":     "libx265",
	"codec:a": "aac",
	"b:a":     "125k",
	"r":       "30",
	"crf":     "26",
}

func (a Transcode) mkAllDir() error {
	if err := _mkdir(a.InputFolder); err != nil {
		return err
	}
	if err := _mkdir(a.OutputFolder); err != nil {
		return err
	}
	return nil
}

func _mkdir(dirName string) error {
	_, err := os.ReadDir(dirName)
	if err != nil {
		if err := os.Mkdir(dirName, 0777); err != nil {
			return err
		}
	}
	return nil
}

func checkFileExists(saveFile string) bool {
	_, err := os.Stat(saveFile)
	if err != nil {
		return false
	}
	return true
}

func logs(format string, v ...interface{}) {
	log.Printf(format+"\n", v)
}

func downloadFile(_link string, saveFile string) error {
	if checkFileExists(saveFile) {
		logs("file: %s exists", saveFile)
		return nil
	}
	resp, err := http.Get(_link)
	if err != nil {
		logs("Get file %s error", saveFile)
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	_file, err := os.Create(saveFile)
	if err != nil {
		logs("Create file %s error", saveFile)
		return err
	}
	defer func() {
		_ = _file.Close()
	}()
	_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs("ReadAll file %s error", saveFile)
		return err
	}
	_, err = _file.Write(_body)
	if err != nil {
		logs("Write file %s error", saveFile)
		return err
	}
	logs("Download file %s successful", saveFile)
	return nil
}

type TranscodeType uint8

const (
	TranscodeTypeUnknown  TranscodeType = 0 //未知
	TranscodeType720P264  TranscodeType = 1 //720p264
	TranscodeType720P265  TranscodeType = 2 //720p265
	TranscodeType1080P264 TranscodeType = 3 //1080p264
	TranscodeType1080P265 TranscodeType = 4 //1080p265
)

func (a TranscodeType) ToArgs() ffmpeg.KwArgs {
	switch a {
	case TranscodeType720P264:
		return h264720Args
	case TranscodeType720P265:
		return h265720Args
	case TranscodeType1080P264:
		return h2641080Args
	case TranscodeType1080P265:
		return h2651080Args
	}
	return h264720Args //default 720p
}

type Transcode struct {
	OutputFolder     string        `json:"output_folder"`
	InputFolder      string        `json:"input_folder"`
	ParallelNumber   int           `json:"parallel_number"`
	OutputFilePrefix string        `json:"output_file_prefix"`
	Files            []string      `json:"files"`
	TranscodeArgs    ffmpeg.KwArgs `json:"transcode_args"` //转码参数
	TranscodeType    TranscodeType `json:"transcode_type"`
	AutoSize         bool          `json:"auto_size"` //自动宽高比
}

func (a *Transcode) defaultParams() {
	if a.ParallelNumber == 0 {
		a.ParallelNumber = 1
	}
	if a.OutputFolder == "" {
		a.OutputFolder = "output"
	}
	if a.InputFolder == "" {
		a.InputFolder = "input"
	}
	if a.TranscodeType == TranscodeTypeUnknown {
		a.TranscodeType = TranscodeType720P264
	}
	a.TranscodeArgs = a.TranscodeType.ToArgs()
}

type SetTranscodeParams func(o *Transcode)

func SetAutoSize(b bool) SetTranscodeParams {
	return func(o *Transcode) {
		o.AutoSize = b
	}
}

func SetTranscodeType(s TranscodeType) SetTranscodeParams {
	return func(o *Transcode) {
		o.TranscodeType = s
	}
}

func SetOutputFolder(s string) SetTranscodeParams {
	return func(o *Transcode) {
		o.OutputFolder = strings.TrimSuffix(s, "/")
	}
}

func SetInputFolder(s string) SetTranscodeParams {
	return func(o *Transcode) {
		o.InputFolder = strings.TrimSuffix(s, "/")
	}
}

func SetParallel(s string) SetTranscodeParams {
	return func(o *Transcode) {
		o.InputFolder = s
	}
}

func SetOutputFilePrefix(s string) SetTranscodeParams {
	return func(o *Transcode) {
		o.OutputFilePrefix = s
	}
}

func NewTranscode(files []string, options ...SetTranscodeParams) *Transcode {
	var s Transcode
	for _, v := range options {
		v(&s)
	}
	s.Files = files
	s.defaultParams()
	return &s
}

func ChunkStrArray(list []string, chunkSize int) [][]string {
	length := len(list)
	chunks := int(math.Ceil(float64(length) / float64(chunkSize))) //each count
	count := 0
	result := make([][]string, 0)
	for i := 0; i < chunkSize; i++ {
		rows := make([]string, 0)
		for g := 0; g < chunks; g++ {
			if count < length {
				rows = append(rows, list[count])
				count++
			}
		}
		if len(rows) > 0 {
			result = append(result, rows)
		}
	}
	return result
}

func (a Transcode) DoTranscode() error {
	if err := a.mkAllDir(); err != nil {
		log.Fatalln(err)
		return err
	}
	chunkSize := ChunkStrArray(a.Files, a.ParallelNumber)
	var wg sync.WaitGroup
	for _, v := range chunkSize {
		wg.Add(1)
		go func(list []string) {
			defer wg.Done()
			for _, j := range list {
				a.handleFile(j)
			}
		}(v)
	}
	wg.Wait()
	return nil
}

func (a Transcode) handleFile(fileName string) {
	originalName, targetName := a.getName(fileName)
	originalFileLocation := a.InputFolder + "/" + originalName
	targetFileLocation := a.OutputFolder + "/" + targetName
	if err := downloadFile(fileName, originalFileLocation); err != nil {
		logs("download file error %s", fileName)
		return
	}
	videoInfo, err := getVideoInfo(originalFileLocation)
	if err != nil {
		logs("getVideoInfoError %+v", err)
		return
	}
	_, err = os.Stat(targetFileLocation)
	if err == nil {
		logs("output file exists %s", targetFileLocation)
		return
	}
	vf := videoInfo.get720pVfParams()
	if a.TranscodeType == TranscodeType1080P264 || a.TranscodeType == TranscodeType1080P265 {
		vf = videoInfo.get1080pVfParams()
	}
	obj := ffmpeg.Input(originalFileLocation, nil)
	if a.AutoSize {
		obj = obj.Filter("scale", ffmpeg.Args{vf})
	}
	if err := obj.Output(targetFileLocation, a.TranscodeArgs).Run(); err != nil {
		logs("run error %+v", err)
	} else {
		logs("successful %s", fileName)
	}
}

func (a Transcode) getName(originalName string) (_originalName, targetName string) {
	_url, err := url.Parse(originalName)
	if err != nil {
		_originalName = fmt.Sprintf("%d_xxxxxx.mp4", time.Now().Unix())
	} else {
		exp := strings.Split(_url.Path, "/")
		_originalName = exp[len(exp)-1]
	}
	targetName = a.OutputFilePrefix + _originalName
	return
}

type FFProbeFormat struct {
	Streams []*FFProbeFormatStream `json:"streams"`
	Format  *FFProbeFormatFormat   `json:"format"`
}

func (a FFProbeFormat) getWidthAndHeight() (width, height int) {
	if len(a.Streams) > 0 {
		width = a.Streams[0].Width
		height = a.Streams[0].Height
	}
	return
}

func (a FFProbeFormat) get720pVfParams() string {
	width, height := a.getWidthAndHeight()
	if width > height {
		return "trunc(oh*a/2)*2:720"
	}
	return "720:trunc(ow/a/2)*2"
}

func (a FFProbeFormat) get1080pVfParams() string {
	width, height := a.getWidthAndHeight()
	if width > height {
		return "trunc(oh*a/2)*2:1080"
	}
	return "1920:trunc(ow/a/2)*2"
}

type FFProbeFormatStreamTags struct {
	Language    string `json:"language"`
	HandlerName string `json:"handler_name"`
	VendorID    string `json:"vendor_id"`
}

type FFProbeFormatStreamDisposition struct {
	Default         int `json:"default"`
	Dub             int `json:"dub"`
	Original        int `json:"original"`
	Comment         int `json:"comment"`
	Lyrics          int `json:"lyrics"`
	Karaoke         int `json:"karaoke"`
	Forced          int `json:"forced"`
	HearingImpaired int `json:"hearing_impaired"`
	VisualImpaired  int `json:"visual_impaired"`
	CleanEffects    int `json:"clean_effects"`
	AttachedPic     int `json:"attached_pic"`
	TimedThumbnails int `json:"timed_thumbnails"`
	Captions        int `json:"captions"`
	Descriptions    int `json:"descriptions"`
	Metadata        int `json:"metadata"`
	Dependent       int `json:"dependent"`
	StillImage      int `json:"still_image"`
}

type FFProbeFormatFormat struct {
	Filename       string                   `json:"filename"`
	NbStreams      int                      `json:"nb_streams"`
	NbPrograms     int                      `json:"nb_programs"`
	FormatName     string                   `json:"format_name"`
	FormatLongName string                   `json:"format_long_name"`
	StartTime      string                   `json:"start_time"`
	Duration       string                   `json:"duration"`
	Size           string                   `json:"size"`
	BitRate        string                   `json:"bit_rate"`
	ProbeScore     int                      `json:"probe_score"`
	Tags           *FFProbeFormatFormatTags `json:"tags"`
}

type FFProbeFormatFormatTags struct {
	MajorBrand       string `json:"major_brand"`
	MinorVersion     string `json:"minor_version"`
	CompatibleBrands string `json:"compatible_brands"`
	Hw               string `json:"Hw"`
	BitRate          string `json:"bitRate"`
	Maxrate          string `json:"maxrate"`
	TeIsReencode     string `json:"te_is_reencode"`
	Encoder          string `json:"encoder"`
}

type FFProbeFormatStream struct {
	Index            int                             `json:"index"`
	CodecName        string                          `json:"codec_name"`
	CodecLongName    string                          `json:"codec_long_name"`
	Profile          string                          `json:"profile"`
	CodecType        string                          `json:"codec_type"`
	CodecTagString   string                          `json:"codec_tag_string"`
	CodecTag         string                          `json:"codec_tag"`
	Width            int                             `json:"width,omitempty"`
	Height           int                             `json:"height,omitempty"`
	CodedWidth       int                             `json:"coded_width,omitempty"`
	CodedHeight      int                             `json:"coded_height,omitempty"`
	ClosedCaptions   int                             `json:"closed_captions,omitempty"`
	FilmGrain        int                             `json:"film_grain,omitempty"`
	HasBFrames       int                             `json:"has_b_frames,omitempty"`
	PixFmt           string                          `json:"pix_fmt,omitempty"`
	Level            int                             `json:"level,omitempty"`
	ChromaLocation   string                          `json:"chroma_location,omitempty"`
	FieldOrder       string                          `json:"field_order,omitempty"`
	Refs             int                             `json:"refs,omitempty"`
	IsAvc            string                          `json:"is_avc,omitempty"`
	NalLengthSize    string                          `json:"nal_length_size,omitempty"`
	ID               string                          `json:"id"`
	RFrameRate       string                          `json:"r_frame_rate"`
	AvgFrameRate     string                          `json:"avg_frame_rate"`
	TimeBase         string                          `json:"time_base"`
	StartPts         int                             `json:"start_pts"`
	StartTime        string                          `json:"start_time"`
	DurationTs       int                             `json:"duration_ts"`
	Duration         string                          `json:"duration"`
	BitRate          string                          `json:"bit_rate"`
	BitsPerRawSample string                          `json:"bits_per_raw_sample,omitempty"`
	NbFrames         string                          `json:"nb_frames"`
	ExtradataSize    int                             `json:"extradata_size"`
	Disposition      *FFProbeFormatStreamDisposition `json:"disposition"`
	Tags             *FFProbeFormatStreamTags        `json:"tags"`
	SampleFmt        string                          `json:"sample_fmt,omitempty"`
	SampleRate       string                          `json:"sample_rate,omitempty"`
	Channels         int                             `json:"channels,omitempty"`
	ChannelLayout    string                          `json:"channel_layout,omitempty"`
	BitsPerSample    int                             `json:"bits_per_sample,omitempty"`
}

func getVideoInfo(originalFile string) (*FFProbeFormat, error) {
	var outBuf bytes.Buffer
	cmd := exec.Command("ffprobe", "-v", "quiet", "-print_format", "json", "-show_format", "-show_streams", originalFile)
	cmd.Stdout = &outBuf
	err := cmd.Run()
	if err != nil {
		logs("getVideoWidthAndHeight command error")
		return nil, err
	}
	var data FFProbeFormat
	if err := json.Unmarshal(outBuf.Bytes(), &data); err != nil {
		logs("getVideoWidthAndHeight json unmashal error")
		return nil, err
	}
	return &data, nil
}
