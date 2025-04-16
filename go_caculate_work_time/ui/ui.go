package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func CreateUI() {
	myApp := app.New()
	myWindow := myApp.NewWindow("云尖无偿加班统计")
	myWindow.Resize(fyne.NewSize(600, 400))

	// 界面组件
	fileEntry := widget.NewEntry()
	fileEntry.SetPlaceHolder("输入文件路径或拖放文件")

	// 结果列表组件（修改为带缓存的高效列表）
	var resultData []PrintData
	var statsLabel = widget.NewLabel("统计信息：")
	resultList := widget.NewList(
		func() int { return len(resultData) },
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(
				fmt.Sprintf("日期:%s    最早:%s    最晚:%s    弹性:%s    有效:%0.2fh",
					resultData[id].Day,
					resultData[id].CheckInAt,
					resultData[id].CheckOutAt,
					b1(resultData[id].IsBound),
					resultData[id].Hours,
				))
		},
	)

	// 内容显示容器
	contentContainer := container.NewMax()

	// 模拟解析逻辑
	mockParse := func(filePath string) (Res, error) {
		ds := parseFile(filePath)
		if len(ds.Data) == 0 {
			return Res{}, fmt.Errorf("解析失败：模拟错误")
		}
		return ds, nil
	}

	// 更新结果列表函数
	updateResults := func(items Res) {
		resultData = items.Data
		resultList.Refresh() // 关键：直接刷新列表数据
		// 更新统计信息
		statsText := fmt.Sprintf("剩余工时: %.2fh    弹性打卡次数: %d次    迟到次数: %d", items.All.Hous, items.All.BoundTimes, items.All.LateTimes)
		statsLabel.SetText(statsText)
	}

	// 创建操作按钮
	selectButton := widget.NewButton("选择文件", func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err == nil && reader != nil {
				fileEntry.SetText(reader.URI().Path())
			}
		}, myWindow)
	})

	parseButton := widget.NewButton("解析数据", func() {
		filePath := fileEntry.Text
		if filePath == "" {
			dialog.ShowError(fmt.Errorf("请先选择文件"), myWindow)
			return
		}

		data, err := mockParse(filePath)
		if err != nil {
			contentContainer.Objects = []fyne.CanvasObject{
				widget.NewLabelWithStyle(
					err.Error(),
					fyne.TextAlignCenter,
					fyne.TextStyle{Bold: true},
				),
			}
		} else {
			updateResults(data)
			contentContainer.Objects = []fyne.CanvasObject{
				container.NewBorder(
					widget.NewLabel("解析结果："),
					container.NewHBox( // 底部统计行
						layout.NewSpacer(),
						statsLabel,
					),
					nil, nil,
					container.NewVScroll(resultList),
				),
			}
		}
		contentContainer.Refresh()
	})

	// 关键修改：使用布局控制输入框尺寸
	topBar := container.NewHBox(
		container.New(layout.NewGridWrapLayout(fyne.NewSize(400, 40)), fileEntry),
		layout.NewSpacer(),
		container.NewHBox(
			selectButton,
			parseButton,
		),
	)

	// 主布局
	mainContainer := container.NewBorder(
		topBar,           // 顶部
		nil,              // 底部
		nil,              // 左侧
		nil,              // 右侧
		contentContainer, // 中心内容
	)

	// 拖放处理
	myWindow.SetOnDropped(func(pos fyne.Position, uris []fyne.URI) {
		if len(uris) > 0 {
			fileEntry.SetText(uris[0].Path())
		}
	})

	myWindow.SetContent(mainContainer)
	myWindow.ShowAndRun()
}

func b1(b bool) string {
	if b {
		return "✅"
	}
	return "❌"
}
