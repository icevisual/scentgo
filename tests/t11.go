package tests

import "fmt"


type Any interface {}

// 字典结构
type Dictionary struct {
	data map[Any]Any // 键值都为interface{}类型
}
// 根据键获取值
func (d *Dictionary) Get(key Any) Any {
	return d.data[key]
}
// 设置键值
func (d *Dictionary) Set(key Any, value Any) {
	d.data[key] = value
}
// 遍历所有的键值，如果回调返回值为false，停止遍历
func (d *Dictionary) Visit(callback func(k, v Any) bool) {
	if callback == nil {
		return
	}
	for k, v := range d.data {
		if !callback(k, v) {
			return
		}
	}
}
// 清空所有的数据
func (d *Dictionary) Clear() {
	d.data = make(map[Any]Any)
}
// 创建一个字典
func NewDictionary() *Dictionary {
	d := &Dictionary{}
	// 初始化map
	d.Clear()
	return d
}
func T11Run() {
	// 创建字典实例
	dict := NewDictionary()
	// 添加游戏数据
	dict.Set("My Factory", "ad")
	dict.Set("My Factory2", 66)
	dict.Set("Terra Craft", 36)
	dict.Set("Don't Hungry", 24)
	// 获取值及打印值
	favorite := dict.Get("Terra Craft")
	fmt.Println("favorite:", favorite)
	// 遍历所有的字典元素
	dict.Visit(func(key, value Any) bool {
		// 将值转为int类型，并判断是否大于40
		ival,isInt :=  value.(int)
		if isInt && ival > 40 {
			// 输出很贵
			fmt.Println(key, "is expensive")
			return true
		}
		// 默认都是输出很便宜
		fmt.Println(key, "is cheap")
		return true
	})
}

