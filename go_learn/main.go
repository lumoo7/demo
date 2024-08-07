package main

import (
	"fmt"
	"go_learn/lib"
	"gorm.io/gorm"
)

type Device struct {
	gorm.Model
	Name         string `gorm:"type:varchar(20);not null;comment:设备名称" json:"name"`
	SerialNumber string `gorm:"type:varchar(20);comment:序列号" json:"serialNumber""`
	Type         string `gorm:"type:varchar(20);comment:设备类型" json:"type"`
	IP           string `gorm:"type:varchar(20);comment:设备IP" json:"ip"`
}

type Switch struct {
	gorm.Model
	Contact     string `gorm:"type:varchar(64);comment:联系人" json:"contact"`
	Description string `gorm:"type:varchar(2000);comment:描述" json:"description"`
	Location    string `gorm:"type:varchar(128);comment:位置" json:"location"`

	// switch belongs to device
	DeviceId uint
	Device   Device
}

func main() {
	lib.InitMysql()
	err := lib.DB.AutoMigrate(&Device{}, &Switch{})
	if err != nil {
		fmt.Println(err)
	}
	pageList()
}

func batchInsert() {
	var dev = []Switch{
		Switch{
			Contact:     "Tom",
			Description: "暂无",
			Location:    "cq",
			Device: Device{
				Name:         "Device1",
				SerialNumber: "001",
				Type:         "T1",
				IP:           "192.168.1.1",
			},
		},
		Switch{
			Contact:     "Tom",
			Description: "暂无",
			Location:    "cq",
			Device: Device{
				Name:         "Device2",
				SerialNumber: "002",
				Type:         "T1",
				IP:           "192.168.1.2",
			},
		},
		Switch{
			Contact:     "Tom",
			Description: "暂无",
			Location:    "cq",
			Device: Device{
				Name:         "Device3",
				SerialNumber: "003",
				Type:         "T1",
				IP:           "192.168.1.3",
			},
		},
	}
	err := lib.DB.Create(&dev).Error
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("batch insert success")
	}
}

func pageList() {
	var res []Switch
	var count int64
	err := lib.DB.Model(&Switch{}).Preload("Device").Offset(0).Limit(10).Count(&count).Find(&res).Error
	if err != nil {
		fmt.Println(err)
	}
	for _, v := range res {
		fmt.Printf("\n%+v\n", v)
	}
}
