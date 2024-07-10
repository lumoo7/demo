package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/xuri/excelize/v2"
)

const (
	CheckInBaseTime  = "8:30"
	CheckOutBaseTime = "18:30"
)

type Data struct {
	CheckInDayAt string
	CheckInAt    string
	CheckOutAt   string
	Count        string
}

func main() {
	// 1.从参数获取文件地址
	// fmt.Println("address: ", os.Args[1])
	xlsx, err := excelize.OpenFile("/Users/huang/Downloads/上下班打卡_日报_20240701-20240711.xlsx")
	if err != nil {
		fmt.Println("read file error. ", err)
	}
	// 2.读取数据
	rows, err := xlsx.GetRows(xlsx.GetSheetName(xlsx.GetActiveSheetIndex()))
	if err != nil {
		fmt.Println("read file error")
	}
	var dataList []Data
	for i, row := range rows {
		if i < 4 {
			continue
		}
		var data Data
		for i, v := range row {
			switch i {
			case 0:
				data.CheckInDayAt = v
			case 7:
				data.CheckInAt = v
			case 8:
				data.CheckOutAt = v
			case 9:
				data.Count = v
			default:
				continue
			}
		}
		dataList = append(dataList, data)
	}
	// 3.解析数据
	var totalMinutes float64
	morningBaseAt, err := timeFormat(CheckInBaseTime)
	if err != nil {
		fmt.Println(err)
	}
	nightBaseAt, err := timeFormat(CheckOutBaseTime)
	if err != nil {
		fmt.Println(err)
	}
	for _, v := range dataList {
		if v.CheckInAt == "未打卡" && v.CheckOutAt == "未打卡" {
			continue
		}
		// 上班打卡
		t1, _ := timeFormat(v.CheckInAt)
		count, err := strconv.Atoi(v.Count)
		if err != nil {
			fmt.Println(err)
		}
		if t1.Before(morningBaseAt) && count == 2 {
			sub := morningBaseAt.Sub(t1).Minutes()
			totalMinutes += sub
		}
		// 弹性打卡
		if t1.After(morningBaseAt) && count == 2 {
			sub := morningBaseAt.Sub(t1).Minutes()
			totalMinutes += sub
		}
		// 下班打卡
		t2, _ := timeFormat(v.CheckOutAt)
		if t2.After(nightBaseAt) && count == 2 {
			sub := t2.Sub(nightBaseAt).Minutes()
			totalMinutes += sub
		}
		fmt.Printf("%s【最早】%s\t【最晚】%s\t【累计】%f\n", v.CheckInDayAt, v.CheckInAt, v.CheckOutAt, totalMinutes)
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
