# CRH Spider

- station_names 获取全国火车站列表 from [get station names from 12306](https://kyfw.12306.cn/otn/resources/js/framework/station_name.js?station_version=1.9059)
- train_lists 获取车次列表 from [get train lists from 12306](https://kyfw.12306.cn/otn/resources/js/query/train_list.js?scriptVersion=1.2)
- train_group 根据train_lists中的数据获取唯一的列车发行的数据
- train_telecode  根据train_list_group中的数据，将station_names获得的数据进行填补telecode等数据
- get_train_no 根据train_lists中的数据将train_list_group的train_no补充完整