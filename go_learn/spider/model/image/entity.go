package _image

import (
	"log"
	"spider/common/db"

	"gorm.io/gorm"
)

type Image struct {
	gorm.Model
	Name   string `gorm:"type:varchar(50);comment:图片名" json:"name"`
	Url    string `gorm:"type:varchar(200);comment:图片url" json:"url"`
	Folder string `gorm:"type:varchar(100);comment:目录名" json:"folder"`
}

type ImageStu struct {
	ID     uint   `json:"ID"`
	Name   string `json:"name"`
	Url    string `json:"url"`
	Folder string `json:"folder"`
}

func (i *Image) AutoMigrate() {
	if err := db.DB.AutoMigrate(&Image{}); err != nil {
		log.Printf("auto migrate entity Image error. %s", err.Error())
	}
	return
}

func (i *Image) Save(p *Image) (*Image, error) {
	if err := db.DB.Save(p).Error; err != nil {
		return nil, err
	}
	return p, nil
}

func (i *Image) Find(stu *ImageStu) (*Image, error) {
	var findImage = new(Image)
	if err := i.where(stu).Find(findImage).Error; err != nil {
		return nil, err
	}
	return findImage, nil
}

func (i *Image) where(stu *ImageStu) *gorm.DB {
	conn := db.DB.Model(&Image{})
	if stu.ID != 0 {
		conn = conn.Where("id = ?", stu.ID)
	}
	if len(stu.Name) != 0 {
		conn = conn.Where("name = ?", stu.Name)
	}
	if len(stu.Url) != 0 {
		conn = conn.Where("url = ?", stu.Url)
	}
	if len(stu.Folder) != 0 {
		conn = conn.Where("folder = ?", stu.Folder)
	}
	return conn
}
