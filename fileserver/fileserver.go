package fileserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zmisgod/goSpider/utils"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var htmlTitle = "ListDir 文件管理系统"

//ListDir struct
type ListDir struct {
	BasePath      string
	Port          int
	Host          string
	Page          int
	PageSize      int
	ShowType      string
	TemplatePath  string
	absPath       string
	ValidDownload []string
	temBasePath   string
}

//Create 创建listDir
func Create(port int, basePath, host, showType, templatePath string, validDownload []string) (*ListDir, error) {
	if host == "" {
		host = "127.0.0.1"
	}
	notExist := utils.CheckPathIsNotExists(basePath)
	if notExist {
		return nil, errors.New(fmt.Sprintf("base path is not exists: %s", basePath))
	}
	return &ListDir{
		BasePath:      basePath,
		Host:          host,
		Port:          port,
		ShowType:      showType,
		TemplatePath:  templatePath,
		ValidDownload: validDownload,
		temBasePath:   basePath,
	}, nil
}

//ShowList 显示列表
func (h *ListDir) ShowList() (map[string]interface{}, error) {
	info, err := os.Open(h.BasePath)
	defer info.Close()
	rows := make(map[string]interface{})
	if err != nil {
		return rows, err
	}
	rest, err := info.Readdir(-1)
	if err != nil {
		return rows, err
	}
	dir := make([]string, 0)
	files := make([]string, 0)
	chunk := ArrayChunk(h.Page, h.PageSize, rest)
	for _, v := range chunk {
		if v.Name() != "" {
			if v.IsDir() {
				dir = append(dir, v.Name())
			} else {
				files = append(files, v.Name())
			}
		}
	}
	rows["ShowNextPage"] = int(len(chunk)) == h.PageSize
	rows["Dir"] = dir
	rows["File"] = files
	return rows, nil
}

//ArrayChunk 将文件列表按照分页形式输出到页面
func ArrayChunk(page, pageSize int, array []os.FileInfo) []os.FileInfo {
	first := (page - 1) * pageSize
	last := pageSize
	var fileLists []os.FileInfo
	if len(array) < first {
		return fileLists
	} else if len(array) < first+last {
		last = len(array)
	}
	for i := first; i < last+first; i++ {
		if i <= len(array)-1 {
			fileLists = append(fileLists, array[i])
		}
	}
	return fileLists
}

//CreateServer 启动一个http server
func (h *ListDir) CreateServer() error {
	http.HandleFunc("/", h.ShowServer)
	http.HandleFunc("/static/", h.ShowStaticFile)
	http.HandleFunc("/404", h.NotFound)
	log.Println(fmt.Sprintf("start file server on %s:%d", h.Host, h.Port))
	return http.ListenAndServe(h.Host+":"+strconv.Itoa(h.Port), nil)
}

//ShowStaticFile 显示静态文件
func (h *ListDir) ShowStaticFile(w http.ResponseWriter, req *http.Request) {
	requestPath := strings.Replace(req.URL.Path, "/static/", "", 10000)
	i := 0
	for _, v := range h.ValidDownload {
		res := strings.Split(requestPath, v)
		if len(res) == 2 {
			i = 1
			break
		}
	}
	if i == 1 {
		http.ServeFile(w, req, h.BasePath+requestPath)
	} else {
		http.Redirect(w, req, "/404", 302)
	}
}

//ShowServer 显示列表
func (h *ListDir) ShowServer(w http.ResponseWriter, req *http.Request) {
	var page, pageSize int
	var err error
	fPage := req.FormValue("page")
	if fPage == "" {
		page = 1
	} else {
		page, err = strconv.Atoi(fPage)
		if err != nil {
			io.WriteString(w, err.Error())
		}
	}
	fPageSize := req.FormValue("pageSize")
	if fPageSize == "" {
		pageSize = 15
	} else {
		pageSize, err = strconv.Atoi(fPageSize)
		if err != nil {
			io.WriteString(w, err.Error())
		}
	}
	requestPath := strings.Replace(req.URL.Path, ".", "", 10000000)
	h.Page = page
	h.PageSize = pageSize
	h.BasePath = h.BasePath + strings.TrimLeft(requestPath, "/")
	if h.ShowType == "json" {
		h.showJSON(w)
	} else {
		h.showTemplate(w, requestPath)
	}
}

//showJSON 谁出json
func (h *ListDir) showJSON(w http.ResponseWriter) {
	res, err := h.ShowList()
	h.BasePath = h.temBasePath
	result := make(map[string]interface{})
	if err != nil {
		result["code"] = 400
		result["data"] = ""
		result["msg"] = err.Error()
	} else {
		result["code"] = 200
		result["data"] = res
		result["msg"] = "ok"
	}
	jsons, _ := json.Marshal(result)
	io.WriteString(w, string(jsons))
}

//showTemplate 显示模板
func (h *ListDir) showTemplate(w http.ResponseWriter, requestPath string) {
	var (
		tmp *template.Template
		err error
		res *template.Template
	)
	if h.TemplatePath == "" {
		tmp = template.New("list")
		res, err = tmp.Parse(listTemplate)
	} else {
		res, err = tmp.ParseFiles(h.TemplatePath)
	}
	if err != nil {
		fmt.Println(err)
	} else {
		data, _ := h.ShowList()
		h.BasePath = h.temBasePath
		if requestPath == "/" {
			data["URI"] = ""
		} else {
			data["URI"] = requestPath
		}
		data["Page"] = h.Page
		data["Npage"] = h.Page + 1
		data["Ppage"] = h.Page - 1
		data["Title"] = htmlTitle
		data["PageSize"] = h.PageSize
		err = res.Execute(w, data)
		if err != nil {
			fmt.Println(err)
		}
	}
}

//NotFound 显示模板
func (h *ListDir) NotFound(w http.ResponseWriter, req *http.Request) {
	tmp := template.New("notfound")
	res, err := tmp.Parse(notfoundTemplate)
	if err != nil {
		fmt.Println(err)
	} else {
		data := make(map[string]string)
		data["Title"] = htmlTitle

		err = res.Execute(w, data)
		if err != nil {
			fmt.Println(err)
		}
	}
}
