package main

import (
	"fmt"
	"os"
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

func readData(r [][]string,d *[]Data) {
	for i, row := range r {
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
		*d = append(*d, data)
	}
}

func parseData(d *[]Data){
	morningBaseAt, err := timeFormat(CheckInBaseTime)
	if err != nil {
		fmt.Println(err)
	}
	nightBaseAt, err := timeFormat(CheckOutBaseTime)
	if err != nil {
		fmt.Println(err)
	}
	var bounceCount int
	var lateCount int
	var totalMinutes float64
	var todayMinutes float64
	for _, v := range *d {
		var morningSurplus float64
		var nightSurplus float64
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
			morningSurplus = morningBaseAt.Sub(t1).Minutes()
		}else if t1.After(morningBaseAt) && count == 2 {
			bounceCount+=1
			morningSurplus = morningBaseAt.Sub(t1).Minutes()
		}
		// 下班打卡
		t2, _ := timeFormat(v.CheckOutAt)
		if t2.Before(nightBaseAt){
			continue
		}else if t2.After(nightBaseAt) && count == 2 {
			nightSurplus = t2.Sub(nightBaseAt).Minutes()
		}
		todayMinutes = morningSurplus + nightSurplus
		totalMinutes +=todayMinutes
		//fmt.Printf("%s\n【最早】%s【剩余工时】%.0f\n【最晚】%s【剩余工时】%.0f\n【今日累计】%.0f mins\n【当月累计】%.0f\n\n",
		//	v.CheckInDayAt, v.CheckInAt,morningSurplus, v.CheckOutAt,nightSurplus,todayMinutes,totalMinutes)
	}
	if bounceCount > 8{
		lateCount = bounceCount-8
		bounceCount = 8
	}
	fmt.Printf("剩余工时: %.2f\n弹性打卡次数: %d次\n迟到次数: %d", totalMinutes/60,bounceCount,lateCount)
}

func main() {
	fmt.Println(`
注意：
（1）此脚本只适用于每日打卡两次的情况（包括周末加班）
（2）计算方式为：最早打卡时间减去基准时间【8:30】为早上剩余时间，最晚打卡时间减去基准时间【18:30】为晚上剩余时间
               在计算弹性打卡和18:30之前打卡的情况后得到的统计数据为总剩余时间
（3）经过测试后，发现和HR统计结果存在一定的误差，此脚本结果仅作为参考`)
	time.Sleep(time.Second * 2)
	fmt.Println("\n按任意键继续")
	b:=make([]byte,1)
	read, _ := os.Stdin.Read(b)
	if read>0{
		// pass
	}
	xlsx, err := excelize.OpenFile(os.Args[1])
	if err != nil {
		fmt.Println("read file error. ", err)
	}
	rows, err := xlsx.GetRows(xlsx.GetSheetName(xlsx.GetActiveSheetIndex()))
	if err != nil {
		fmt.Println("read file error")
	}
	var dataList []Data
	readData(rows,&dataList)
	parseData(&dataList)
}

func timeFormat(t string) (time.Time, error) {
	t1, err := time.Parse("15:04", t)
	if err != nil {
		return time.Time{}, err
	}
	return t1, nil
}
