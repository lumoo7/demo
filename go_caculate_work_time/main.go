package main

import (
	"fmt"
	"os"
	"time"

	"github.com/xuri/excelize/v2"
)

const (
	CheckInBaseTime  = "8:30"
	CheckOutBaseTime = "18:30"
)

type Data struct {
	CheckInAt  string
	CheckOutAt string
}

func main() {
	// 1.从参数获取文件地址
	fmt.Println("address: ", os.Args[1])
	xlsx, err := excelize.OpenFile(os.Args[1])
	if err != nil {
		fmt.Println("read file error. ", err)
	}
	// 2.读取数据
	rows, err := xlsx.GetRows(xlsx.GetSheetName(xlsx.GetActiveSheetIndex()))
	if err != nil {
		fmt.Println("read file error")
	}
	var dataList []Data
	// 列
	for i, row := range rows {
		if i < 4 {
			continue
		}
		// 行
		var data Data
		for i, v := range row {
			switch i {
			case 7:
				data.CheckInAt = v
			case 8:
				data.CheckOutAt = v
			default:
				continue
			}
		}
		dataList = append(dataList, data)
	}
	// 3.解析数据
	var totalMinutes float64
	for _, v := range dataList {
		if v.CheckInAt == "未打卡" && v.CheckOutAt == "未打卡" {
			continue
		}
		morningBaseAt, err := timeFormat(CheckInBaseTime)
		if err != nil {
			fmt.Println(err)
		}
		nightBaseAt, err := timeFormat(CheckOutBaseTime)
		if err != nil {
			fmt.Println(err)
		}
		// morning
		t1, _ := timeFormat(v.CheckInAt)
		if t1.Before(morningBaseAt) {
			sub := morningBaseAt.Sub(t1).Minutes()
			totalMinutes += sub
		}
		// 弹性打卡
		if t1.After(morningBaseAt) {
			sub := t1.Sub(morningBaseAt).Minutes()
			totalMinutes -= sub
		}

		// night
		t2, _ := timeFormat(v.CheckOutAt)
		if t2.After(nightBaseAt) {
			sub := t2.Sub(nightBaseAt).Minutes()
			totalMinutes += sub
		}
		fmt.Printf("\n最早:%s\t最晚:%s\n", t1, t2)
	}
	fmt.Println("Total Time: ", totalMinutes/60)
}

func timeFormat(t string) (time.Time, error) {
	t1, err := time.Parse("15:04", t)
	if err != nil {
		return time.Time{}, err
	}
	return t1, nil
}
