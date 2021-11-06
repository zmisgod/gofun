package tencent

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
)

type Resp struct {
	Status int          `json:"status"`
	Result [][]*OneInfo `json:"result"`
}

type OneInfo struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"fullname"`
}

type FetchCity struct {
	OriginalData io.Reader   `json:"original_data"`
	ParseData    *Resp       `json:"parse_data"`
	CityList     []*CityTree `json:"city_list"`
}

type CityTree struct {
	ID    string      `json:"id"`
	Name  string      `json:"name"`
	List  []*CityTree `json:"list"`
	Level int         `json:"level"`
}

func NewCity(fileObj io.Reader) (*FetchCity, error) {
	re, err := ioutil.ReadAll(fileObj)
	if err != nil {
		return nil, err
	}
	res := FetchCity{}
	var resp Resp
	err = json.Unmarshal(re, &resp)
	if err != nil {
		return nil, err
	}
	res.ParseData = &resp
	res.handleData()
	return &res, nil
}

func getStringSplit(a string) (pre, mid, sub string) {
	sub = a[4:]
	mid = a[2:4]
	pre = a[0:2]
	return
}

func (a *FetchCity) handleData() {
	var levelOne []*CityTree
	levelTwo := make(map[string][]CityTree)
	levelThree := make(map[string][]CityTree)
	directCity := make(map[string]map[string]CityTree)
	if len(a.ParseData.Result) > 0 {
		for _, v := range a.ParseData.Result {
			for _, j := range v {
				_, mid, sub := getStringSplit(j.ID)
				if mid == "00" && sub == "00" {
					levelOne = append(levelOne, &CityTree{
						ID:    j.ID,
						Name:  j.FullName,
						Level: 1,
					})
				}
			}
		}
	}

	for _, v := range levelOne {
		oneLevelPre, _, _ := getStringSplit(v.ID)
		for _, j := range a.ParseData.Result {
			for _, x := range j {
				nowPre, nowMid, nowSub := getStringSplit(x.ID)
				if oneLevelPre == nowPre && nowSub == "00" && nowMid != "00" {
					levelTwo[v.ID] = append(levelTwo[v.ID], CityTree{
						ID:    x.ID,
						Name:  x.FullName,
						Level: 2,
					})
				}
				//北上广深、澳门、香港 需要特殊处理
				if oneLevelPre == nowPre && (oneLevelPre == "31" || oneLevelPre == "50" || oneLevelPre == "11" ||
					oneLevelPre == "12" || oneLevelPre == "81" || oneLevelPre == "82") {
					if nowSub != "00" {
						ex, ok := directCity[v.ID]
						if ok {
							ex[x.FullName] = CityTree{
								ID:    x.ID,
								Name:  x.FullName,
								Level: 2,
							}
							directCity[v.ID] = ex
						} else {
							re := make(map[string]CityTree)
							re[v.ID] = CityTree{
								ID:    x.ID,
								Name:  x.FullName,
								Level: 2,
							}
							directCity[v.ID] = re
						}
					}
				}
			}
		}
	}

	for _, v := range levelTwo {
		//城市列表
		for _, j := range v {
			twoLevelPre, twoLevelMid, _ := getStringSplit(j.ID)
			//获取城市下面的区信息
			for _, x := range a.ParseData.Result {
				for _, s := range x {
					nowPre, nowMid, nowSub := getStringSplit(s.ID)
					if twoLevelPre == nowPre && twoLevelMid == nowMid && nowSub != "00" {
						levelThree[j.ID] = append(levelThree[j.ID], CityTree{
							ID:    s.ID,
							Name:  s.FullName,
							Level: 3,
						})
					}
				}
			}
		}
	}
	for _, v := range levelOne {
		ex, ok := levelTwo[v.ID]
		if ok {
			twoList := make([]*CityTree, 0)
			for _, j := range ex {
				threeList := make([]*CityTree, 0)
				city, _ok := levelThree[j.ID]
				if _ok {
					for _, s := range city {
						threeList = append(threeList, &CityTree{
							ID:    s.ID,
							Name:  s.Name,
							List:  nil,
							Level: 3,
						})
					}
				}
				twoList = append(twoList, &CityTree{
					ID:    j.ID,
					Name:  j.Name,
					List:  threeList,
					Level: 2,
				})
			}
			v.List = twoList
		} else {
			one, _, _ := getStringSplit(v.ID)
			threeList := make([]*CityTree, 0)
			ext, _ok := directCity[v.ID]
			if _ok {
				for _, j := range ext {
					threeList = append(threeList, &CityTree{
						ID:    j.ID,
						Name:  j.Name,
						List:  nil,
						Level: 3,
					})
				}
			}
			v.List = append(v.List, &CityTree{
				ID:    fmt.Sprintf("%s0100", one),
				Name:  v.Name,
				List:  threeList,
				Level: 2,
			})
		}
	}
	a.CityList = levelOne
}

func (a *FetchCity) GetAllProvince() []string {
	list := make([]string, 0)
	for _, v := range a.CityList {
		list = append(list, v.Name)
	}
	return list
}

func (a *FetchCity) GetCitiesByProvince(province string) []string {
	list := make([]string, 0)
	for _, v := range a.CityList {
		if v.Name == province {
			for _, j := range v.List {
				list = append(list, j.Name)
			}
		}
	}
	return list
}

func (a *FetchCity) GetDistrictByCityName(province string, cityName string) []string {
	list := make([]string, 0)
	for _, v := range a.CityList {
		if v.Name == province {
			for _, j := range v.List {
				if j.Name == cityName {
					for _, x := range j.List {
						list = append(list, x.Name)
					}
				}
			}
		}
	}
	return list
}
