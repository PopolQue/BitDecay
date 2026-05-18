package ui

import "image"

var (
	// Current Resolution
	ScreenWidth  = 1280
	ScreenHeight = 768
)

func IsPortrait() bool {
	return ScreenHeight > ScreenWidth
}

func GetWidgetRect() image.Rectangle {
	if IsPortrait() {
		return image.Rect(10, 100, ScreenWidth-10, ScreenHeight-10)
	}
	return image.Rect((ScreenWidth-1000)/2, (ScreenHeight-600)/2+50, (ScreenWidth-1000)/2+1000, (ScreenHeight-600)/2+650)
}

func GetMetricsRect() image.Rectangle {
	wr := GetWidgetRect()
	if IsPortrait() {
		return image.Rect(wr.Min.X+10, wr.Min.Y+10, wr.Max.X-10, wr.Min.Y+200)
	}
	return image.Rect(wr.Min.X+20, wr.Min.Y+20, wr.Max.X-20, wr.Min.Y+140)
}

func GetClickerRect() image.Rectangle {
	wr := GetWidgetRect()
	mr := GetMetricsRect()
	if IsPortrait() {
		return image.Rect(wr.Min.X+10, mr.Max.Y+20, wr.Max.X-10, mr.Max.Y+140)
	}
	return image.Rect(wr.Min.X+20, wr.Min.Y+160, wr.Min.X+320, wr.Min.Y+280)
}

func GetHardwareRect() image.Rectangle {
	wr := GetWidgetRect()
	if IsPortrait() {
		return image.Rect(wr.Min.X+10, ScreenHeight-400, wr.Max.X-10, ScreenHeight-220)
	}
	return image.Rect(wr.Min.X+340, wr.Min.Y+160, wr.Min.X+650, wr.Min.Y+500)
}

func GetUpgradeRect() image.Rectangle {
	wr := GetWidgetRect()
	if IsPortrait() {
		return image.Rect(wr.Min.X+10, ScreenHeight-210, wr.Max.X-10, ScreenHeight-80)
	}
	return image.Rect(wr.Min.X+670, wr.Min.Y+160, wr.Max.X-20, wr.Min.Y+500)
}

func GetLogRect() image.Rectangle {
	wr := GetWidgetRect()
	cr := GetClickerRect()
	if IsPortrait() {
		return image.Rect(wr.Min.X+10, cr.Max.Y+100, wr.Max.X-10, cr.Max.Y+220)
	}
	return image.Rect(wr.Min.X+20, wr.Min.Y+360, wr.Min.X+320, wr.Min.Y+510)
}

func GetPacketRect() image.Rectangle {
	wr := GetWidgetRect()
	if IsPortrait() {
		return image.Rect(wr.Min.X+10, ScreenHeight-70, wr.Min.X+150, ScreenHeight-10)
	}
	return image.Rect(wr.Min.X+20, wr.Min.Y+520, wr.Min.X+320, wr.Min.Y+580)
}

func GetRebootRect() image.Rectangle {
	wr := GetWidgetRect()
	if IsPortrait() {
		return image.Rect(wr.Max.X-150, ScreenHeight-70, wr.Max.X-10, ScreenHeight-10)
	}
	return image.Rect(wr.Min.X+340, wr.Min.Y+520, wr.Max.X-20, wr.Min.Y+580)
}

