package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"spider/common/db"
	_image "spider/model/image"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
)

func r18(pageUrl, savePath string) {
	c := colly.NewCollector(
		colly.Debugger(&debug.WebDebugger{}),
	)
	proxyURL, err := url.Parse("http://127.0.0.1:7890")
	if err != nil {
		panic(err)
	}
	c.WithTransport(&http.Transport{
		Proxy: http.ProxyURL(proxyURL),
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second, // 设置连接超时时间为30秒
			KeepAlive: 30 * time.Second, // 设置保持连接的时间为30秒
		}).DialContext,
		MaxIdleConns:          100,              // 设置最大空闲连接数
		IdleConnTimeout:       90 * time.Second, // 设置空闲连接超时时间
		TLSHandshakeTimeout:   10 * time.Second, // 设置TLS握手超时时间
		ExpectContinueTimeout: 1 * time.Second,  // 设置服务器期望的继续时间
	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36")
		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
		r.Headers.Set("Accept-Language", "en-US,en;q=0.5")
		r.Headers.Set("Referer", "http://example.com")
		r.Headers.Set("Connection", "keep-alive")

		fmt.Println("request")
	})

	c.OnResponse(func(r *colly.Response) {})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("error: ", err.Error())
	})

	c.OnHTML("a.item-link", func(h *colly.HTMLElement) {
		fileName := strings.ReplaceAll(h.Text, "\n", "")
		fileUrl := h.Attr("href")
		var count int64
		if err := db.DB.Model(&_image.Image{}).
			Where("folder = ?", fileName).
			Count(&count).
			Error; err != nil {
			return
		}
		fmt.Println(fileName, " ", count)
		if count > 0 {
			return
		} else {
			h.Request.Ctx.Put("folderName", fileName)
			h.Request.Visit(fileUrl)
		}
	})

	c.Limit(&colly.LimitRule{Parallelism: 10})

	c.OnHTML("#masonry > div > img", func(h *colly.HTMLElement) {
		folder := h.Request.Ctx.Get("folderName")
		imageUrl := h.Attr("data-original")
		imageName := h.Attr("title")
		var imageModel _image.Image
		if err := db.DB.Model(&_image.Image{}).
			Where("folder = ? and url = ? and name = ?", folder, imageUrl, imageName).
			Find(&imageModel).
			Error; err != nil {
			return
		}
		if imageModel.ID != 0 {
			return
		}
		childFolder := fmt.Sprintf("%s/r18/%s", savePath, folder)
		fullPath := fmt.Sprintf("%s/%s.jpg", childFolder, imageName)
		if _, err := os.Stat(childFolder); os.IsNotExist(err) {
			err := os.MkdirAll(childFolder, 0755)
			if err != nil {
				log.Fatalf("Failed to create directory: %s", err)
			}
		}
		resp, err := http.Get(imageUrl)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		out, err := os.Create(fullPath)
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		if err := db.DB.Create(&_image.Image{
			Name:   imageName,
			Url:    imageUrl,
			Folder: folder,
		}).Error; err != nil {
			log.Fatal("create image info error")
			return
		}
	})

	c.OnHTML("li.next > a", func(h *colly.HTMLElement) {
		nextPageUrl := h.Attr("href")
		h.Request.Visit(nextPageUrl)
	})

	c.Visit(pageUrl)

	c.Wait()
}

func main() {
	var savePath = "./" // "/Volumes/黄大壮/images"
	var url = "https://www.hentaiclub.net/r15/49011.html"
	go r18(url, savePath)
	select {}
}
