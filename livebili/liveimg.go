package livebili

import (
	"bytes"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"image"
	"image/jpeg"
	"time"
	"unicode/utf8"
)

type LiveImg struct {
	// 字体文件路径
	TTFPath string
	// 头像
	Avatar image.Image
	// 昵称
	NickName string
}

func NewLiveImg(ttf string, ava image.Image, nickName string) *LiveImg {
	return &LiveImg{TTFPath: ttf, Avatar: ava, NickName: nickName}
}

func (l *LiveImg) DrawOnLive(cover image.Image, liveTitle string, lastOffTime time.Time) (image.Image, error) {
	dw, dh := 744, 700
	pw, ph := 704.0, 396.0
	dc := gg.NewContext(dw, dh)
	cover = imaging.Resize(cover, int(pw), int(ph), imaging.Lanczos)
	blurred := imaging.Blur(cover, 100.0)
	blurred = imaging.Resize(blurred, dw, dh, imaging.Lanczos)

	dc.DrawImage(blurred, 0, 0)

	// 图片框的位置和尺寸
	x, y := 20.0, 20.0
	radius := 30.0           // 圆角半径
	borderWidth := 1.5       // 边框宽度
	dc.SetRGBA(1, 1, 1, 0.5) // 设置边框颜色

	drawRoundedRect(dc, x-borderWidth, y-borderWidth, pw+2*borderWidth, ph+2*borderWidth, radius+borderWidth)
	dc.Fill()
	drawImageInRoundedRect(dc, cover, x, y, pw, ph, radius)
	img := dc.Image()
	dc = gg.NewContextForImage(img)

	dc.SetRGBA(1, 1, 1, 0.3)
	err := drawLive(dc, l.TTFPath, 20, 450, l.Avatar, l.NickName, liveTitle, lastOffTime)
	if err != nil {
		return nil, err
	}

	return dc.Image(), nil
}

func (l *LiveImg) DrawOffLive(desc string, lastOnTime time.Time) (image.Image, error) {
	dw, dh := 744, 240
	dc := gg.NewContext(dw, dh)
	blurred := imaging.Blur(l.Avatar, 150.0)
	blurred = imaging.Resize(blurred, dw, dh, imaging.Lanczos)

	dc.DrawImage(blurred, 0, 0)

	dc.SetRGBA(1, 1, 1, 0.3)
	err := drawOffLive(dc, l.TTFPath, 20, 20, l.Avatar, l.NickName, desc, lastOnTime)
	if err != nil {
		return nil, err
	}
	return dc.Image(), nil
}

func drawLive(dc *gg.Context, ttf string, x, y float64, img image.Image, name string, desc string, lastOffTime time.Time) error {
	drawRoundedRect(dc, x, y, 704, 200, 30)
	dc.Fill()

	dc.SetRGB(0, 0, 0)
	err := dc.LoadFontFace(ttf, 30)
	if err != nil {
		return err
	}
	dc.DrawStringAnchored(name, x+60+60, y+60, 0, 0.5)

	err = dc.LoadFontFace(ttf, 20)
	if err != nil {
		return err
	}
	dc.SetRGBA(1, 0, 0, 0.8)
	dc.DrawStringAnchored("● Live", x+704-60, y+30, 0.5, 0.5)

	err = dc.LoadFontFace(ttf, 15)
	if err != nil {
		return err
	}
	dc.SetRGBA(0, 0, 0, 0.6)
	h, m, s := toNowDuration(lastOffTime)
	dc.DrawStringAnchored(fmt.Sprintf("时隔 %d时%d分%d秒", h, m, s), x+704-20, y+170, 1, 0.5)

	err = dc.LoadFontFace(ttf, 20)
	if err != nil {
		return err
	}
	dc.SetRGBA(0, 0, 0, 0.8)

	var rByte []byte
	rCount := 0
	for _, r := range desc {
		if rCount >= 30 {
			rByte = append(rByte, []byte("...")...)
			break
		}
		rByte = utf8.AppendRune(rByte, r)
		rCount++
	}
	desc = string(rByte)

	dc.DrawStringAnchored(desc, x+30, y+60+75, 0, 0.5)

	// 绘制头像（圆形）
	avatarRadius := 40.0 // 头像的半径，比例可调整
	avatarCenterX := x + 60
	avatarCenterY := y + 60 // 头像位置，居上

	dc.SetRGB(1, 1, 1) // 设置头像背景颜色
	dc.NewSubPath()
	dc.DrawCircle(avatarCenterX, avatarCenterY, avatarRadius)
	dc.ClosePath()
	dc.Clip()

	img = imaging.Resize(img, 80, 80, imaging.Lanczos)

	dc.DrawImageAnchored(img, int(avatarCenterX), int(avatarCenterY), 0.5, 0.5)
	return nil
}

func drawOffLive(dc *gg.Context, ttf string, x, y float64, img image.Image, name string, desc string, lastOnTime time.Time) error {
	drawRoundedRect(dc, x, y, 704, 200, 30)
	dc.Fill()

	dc.SetRGB(0, 0, 0)
	err := dc.LoadFontFace(ttf, 30)
	if err != nil {
		return err
	}
	dc.DrawStringAnchored(name, x+60+60, y+60, 0, 0.5)

	err = dc.LoadFontFace(ttf, 20)
	if err != nil {
		return err
	}
	dc.SetRGBA(0.5, 0.5, 0.5, 0.8)
	dc.DrawStringAnchored("● Live", x+704-60, y+30, 0.5, 0.5)

	err = dc.LoadFontFace(ttf, 15)
	if err != nil {
		return err
	}
	dc.SetRGBA(0, 0, 0, 0.6)
	h, m, s := toNowDuration(lastOnTime)
	dc.DrawStringAnchored(fmt.Sprintf("直播时长 %d时%d分%d秒", h, m, s), x+704-20, y+170, 1, 0.5)

	err = dc.LoadFontFace(ttf, 20)
	if err != nil {
		return err
	}
	dc.SetRGBA(0, 0, 0, 0.8)

	var rByte []byte
	rCount := 0
	for _, r := range desc {
		if rCount >= 30 {
			rByte = append(rByte, []byte("...")...)
			break
		}
		rByte = utf8.AppendRune(rByte, r)
		rCount++
	}
	desc = string(rByte)

	dc.DrawStringAnchored(desc, x+30, y+60+75, 0, 0.5)

	// 绘制头像（圆形）
	avatarRadius := 40.0 // 头像的半径，比例可调整
	avatarCenterX := x + 60
	avatarCenterY := y + 60 // 头像位置，居上

	dc.SetRGB(1, 1, 1) // 设置头像背景颜色
	dc.NewSubPath()
	dc.DrawCircle(avatarCenterX, avatarCenterY, avatarRadius)
	dc.ClosePath()
	dc.Clip()

	img = imaging.Resize(img, 80, 80, imaging.Lanczos)

	dc.DrawImageAnchored(img, int(avatarCenterX), int(avatarCenterY), 0.5, 0.5)
	return nil
}

// 绘制圆角矩形
func drawRoundedRect(dc *gg.Context, x, y, width, height, radius float64) {
	dc.NewSubPath()
	dc.MoveTo(x+radius, y)
	dc.LineTo(x+width-radius, y)
	dc.QuadraticTo(x+width, y, x+width, y+radius)
	dc.LineTo(x+width, y+height-radius)
	dc.QuadraticTo(x+width, y+height, x+width-radius, y+height)
	dc.LineTo(x+radius, y+height)
	dc.QuadraticTo(x, y+height, x, y+height-radius)
	dc.LineTo(x, y+radius)
	dc.QuadraticTo(x, y, x+radius, y)
	dc.ClosePath()
}

// 在圆角矩形中绘制图片
func drawImageInRoundedRect(dc *gg.Context, img image.Image, x, y, width, height, radius float64) {

	// 创建裁剪路径
	dc.Push() // 保存当前上下文状态
	drawRoundedRect(dc, x, y, width, height, radius)
	dc.Clip() // 裁剪为圆角矩形

	// 绘制图片
	dc.DrawImageAnchored(img, int(x+width/2), int(y+height/2), 0.5, 0.5)
	dc.Pop() // 恢复上下文状态

}

func toNowDuration(dur time.Time) (hour int, minute int, second int) {
	now := time.Now()
	sub := now.Sub(dur)
	// 将时间间隔转换为小时、分钟和秒
	hours := int(sub.Hours())
	minutes := int(sub.Minutes()) % 60
	seconds := int(sub.Seconds()) % 60

	return hours, minutes, seconds
}

// ImageToBytes 将 image.Image 转换为 bytes
func ImageToBytes(img image.Image) ([]byte, error) {
	var buf bytes.Buffer

	// 将图像编码为 JPEG 格式
	err := jpeg.Encode(&buf, img, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to encode image to JPEG: %v", err)
	}

	return buf.Bytes(), nil
}
