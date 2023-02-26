package png

import (
	"context"
	"errors"
	"fmt"
	"github.com/zmisgod/gofun/image_research"
	"io/ioutil"
	"os"
)

type PData struct {
	index      uint          `json:"index"`
	dataChunks []interface{} `json:"dataChunks"`
	length     uint          `json:"length"`
	file       *os.File
	data       []byte
	width      uint
	height     uint
	bitDepth   uint
}

var pndValidHeader = []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}

func NewPng(_file string) (*PData, error) {
	fd, err := os.OpenFile(_file, os.O_RDONLY, 0777)
	if err != nil {
		return nil, err
	}
	return &PData{
		file: fd,
	}, nil
}

func (a *PData) decode(ctx context.Context) ([]uint, error) {
	resp := make([]uint, 0)
	data, err := ioutil.ReadAll(a.file)
	if err != nil {
		return resp, err
	}
	a.data = data
	if err := a.decodeHeader(ctx); err != nil {
		return resp, err
	}
	_type, err := a.decodeChunk(ctx)
	if err != nil {
		return resp, err
	}
	fmt.Println(_type)
	return a.getData(ctx)
}

func (a *PData) read(_n int64) ([]byte, error) {
	_bytes := make([]byte, _n)
	_, err := a.file.ReadAt(_bytes, int64(a.index))
	if err != nil {
		return _bytes, err
	}
	a.index += uint(_n)
	return _bytes, nil
}

//解析头部信息
func (a *PData) decodeHeader(ctx context.Context) error {
	data, err := a.read(8)
	if err != nil {
		return err
	}
	if len(data) != len(pndValidHeader) {
		return errors.New("参数错误")
	}
	for k, v := range data {
		if pndValidHeader[k] != v {
			return errors.New("非png")
		}
	}
	return nil
}

//解析IHDR数据块
func (a *PData) decodeChunk(ctx context.Context) (string, error) {
	_data, _ := a.read(4)
	length := image_research.ReadInt32(_data, 0)
	_data, _ = a.read(4)
	_type := image_research.BufferToString(_data)
	chunkData, _ := a.read(int64(length))
	//crc, _ := a.read(4) //crc 冗余校验码
	switch _type {
	case "IHDR":
		a.decodeIHDR(chunkData)
	case "PLTE":
		a.decodePLTE(chunkData)
	case "IDAT":
		a.decodeIDAT(chunkData)
	case "IEND":
		a.decodeIEND(chunkData)
	case "tRNS":
		a.decodetRNS(chunkData)
	}
	return _type, nil
}

func (a *PData) decodeIHDR(data []byte) error {
	a.width = image_research.ReadInt32(data, 0)
	a.height = image_research.ReadInt32(data, 4)
	fmt.Println(a.width, a.height)
	return nil
}

func (a *PData) decodePLTE(data []byte) error {
	return nil
}

func (a *PData) decodeIDAT(data []byte) error {
	return nil
}

func (a *PData) decodeIEND(data []byte) error {
	return nil
}

func (a *PData) decodetRNS(data []byte) error {
	return nil
}

func (a *PData) getData(ctx context.Context) ([]uint, error) {
	list := make([]uint, 0)

	return list, nil
}
