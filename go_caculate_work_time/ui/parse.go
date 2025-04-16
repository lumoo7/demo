package ui

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

type PrintData struct {
	Day        string
	CheckInAt  string
	CheckOutAt string
	Hours      float32
	IsBound    bool
}
type All struct {
	Hous       float32
	LateTimes  int8
	BoundTimes int8
}
type Res struct {
	Data []PrintData
	All  All
}

func readData(r [][]string, d *[]Data) {
	for i, row := range r {
		if i < 4 {
			continue
		}
		checkDay := row[0]
		checkInAt := row[8]
		checkOutAt := row[9]
		checkCount := row[10]
		if checkInAt == "未打卡" || checkOutAt == "未打卡" || checkCount != "2" {
			continue
		}
		*d = append(*d, Data{
			CheckInDayAt: checkDay,
			CheckInAt:    checkInAt,
			CheckOutAt:   checkOutAt,
			Count:        checkCount,
		})
	}
}

func parseData(d *[]Data) Res {
	var res Res
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
	// 统计迟到和弹性打卡次数
	for _, v := range *d {
		isBound := false
		var morningSurplus float64
		var nightSurplus float64
		// 上班打卡
		t1, _ := timeFormat(v.CheckInAt)
		count, err := strconv.Atoi(v.Count)
		if err != nil {
			fmt.Println(err)
		}
		if t1.Before(morningBaseAt) && count == 2 {
			morningSurplus = morningBaseAt.Sub(t1).Minutes()
		} else if t1.After(morningBaseAt) && count == 2 {
			isBound = true
			bounceCount += 1
			morningSurplus = morningBaseAt.Sub(t1).Minutes()
		}
		// 下班打卡
		t2, _ := timeFormat(v.CheckOutAt)
		if t2.Before(nightBaseAt) {
			continue
		} else if t2.After(nightBaseAt) && count == 2 {
			nightSurplus = t2.Sub(nightBaseAt).Minutes()
		}
		todayMinutes = morningSurplus + nightSurplus
		totalMinutes += todayMinutes
		res.Data = append(res.Data, PrintData{
			Day:        v.CheckInDayAt,
			CheckInAt:  v.CheckInAt,
			CheckOutAt: v.CheckOutAt,
			Hours:      float32(todayMinutes / 60),
			IsBound:    isBound,
		})
	}
	if bounceCount > 8 {
		lateCount = bounceCount - 8
		bounceCount = 8
	}
	res.All = All{
		Hous:       float32(totalMinutes / 60),
		LateTimes:  int8(lateCount),
		BoundTimes: int8(bounceCount),
	}
	return res
}

func timeFormat(t string) (time.Time, error) {
	t1, err := time.Parse("15:04", t)
	if err != nil {
		return time.Time{}, err
	}
	return t1, nil
}

func parseFile(path string) Res {
	xlsx, err := excelize.OpenFile(path)
	if err != nil {
		return Res{}
	}
	rows, err := xlsx.GetRows(xlsx.GetSheetName(xlsx.GetActiveSheetIndex()))
	if err != nil {
		return Res{}
	}
	var dataList []Data
	readData(rows, &dataList)
	return parseData(&dataList)

}
