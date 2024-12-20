package livebili

import (
	"fmt"
	"github.com/fogleman/gg"
	"github.com/kohmebot/pkg/canvas"
	"image"
)

type FollowerImg struct {
	// 字体文件路径
	TTFPath string
	// 昵称
	NickName string
	// 头像
	Avatar image.Image
}

func NewFollowerImg(ttf string, ava image.Image, nickName string) *FollowerImg {
	return &FollowerImg{
		TTFPath:  ttf,
		NickName: nickName,
		Avatar:   ava,
	}
}

func (f *FollowerImg) DrawUpFollower(follower int, delta int) (*canvas.Canvas, error) {
	fonts := canvas.NewFonts(f.TTFPath)

	dw, dh := 744, 240
	c := canvas.NewCanvas(dw, dh)
	b := canvas.NewImgBackgroundWithBlur(f.Avatar, 150)
	if err := c.SetBackground(b); err != nil {
		return nil, err
	}

	// 矩形框
	x, y := 20.0, 20.0
	reg := canvas.NewRectangleFrame(704, 200).SetRadius(30).SetRGBA(1, 1, 1, 0.3)
	if err := c.DrawWith(reg, x, y); err != nil {
		return nil, err
	}

	// 头像
	ava := canvas.NewImageCircleFrame(f.Avatar, 40)
	if err := c.DrawWith(ava, x+60, y+60); err != nil {
		return nil, err
	}

	// 昵称
	text := canvas.NewTextFrame(f.NickName, fonts, 30, 0, 704, gg.AlignLeft).SetRGBA(0, 0, 0, 1)
	if err := c.DrawWith(text, x+60+60, y+43); err != nil {
		return nil, err
	}

	// 右上角标记
	text = canvas.NewTextFrame("涨粉了！", fonts, 15, 0, 300, gg.AlignRight).SetRGBA(1, 0, 0, 0.8)
	if err := c.DrawWith(text, x+704-360, y+30); err != nil {
		return nil, err
	}

	// 涨粉数值
	text = canvas.NewTextFrame(fmt.Sprintf("新增 %d 位粉丝", delta), fonts, 20, 0, 704, gg.AlignLeft).
		SetRGBA(0, 0, 0, 1)
	if err := c.DrawWith(text, x+30, y+60+60); err != nil {
		return nil, err
	}

	// 粉丝数标记
	text = canvas.NewTextFrame(fmt.Sprintf("粉丝数 %d", follower), fonts, 15, 0, 300, gg.AlignRight).
		SetRGBA(0, 0, 0, 0.6)
	if err := c.DrawWith(text, x+704-360, y+150); err != nil {
		return nil, err
	}

	return c, nil
}

func (f *FollowerImg) DrawDownFollower(follower int, delta int) (*canvas.Canvas, error) {
	fonts := canvas.NewFonts(f.TTFPath)

	dw, dh := 744, 240
	c := canvas.NewCanvas(dw, dh)
	b := canvas.NewImgBackgroundWithBlur(f.Avatar, 150)
	if err := c.SetBackground(b); err != nil {
		return nil, err
	}

	// 矩形框
	x, y := 20.0, 20.0
	reg := canvas.NewRectangleFrame(704, 200).SetRadius(30).SetRGBA(1, 1, 1, 0.3)
	if err := c.DrawWith(reg, x, y); err != nil {
		return nil, err
	}

	// 头像
	ava := canvas.NewImageCircleFrame(f.Avatar, 40)
	if err := c.DrawWith(ava, x+60, y+60); err != nil {
		return nil, err
	}

	// 昵称
	text := canvas.NewTextFrame(f.NickName, fonts, 30, 0, 704, gg.AlignLeft).SetRGBA(0, 0, 0, 1)
	if err := c.DrawWith(text, x+60+60, y+43); err != nil {
		return nil, err
	}

	// 右上角标记
	text = canvas.NewTextFrame("掉粉了...", fonts, 15, 0, 300, gg.AlignRight).SetRGBA(0, 0.6, 0, 0.8)
	if err := c.DrawWith(text, x+704-360, y+30); err != nil {
		return nil, err
	}

	// 涨粉数值
	text = canvas.NewTextFrame(fmt.Sprintf("失去了 %d 位粉丝", delta), fonts, 20, 0, 704, gg.AlignLeft).
		SetRGBA(0, 0, 0, 1)
	if err := c.DrawWith(text, x+30, y+60+60); err != nil {
		return nil, err
	}

	// 粉丝数标记
	text = canvas.NewTextFrame(fmt.Sprintf("粉丝数 %d", follower), fonts, 15, 0, 300, gg.AlignRight).
		SetRGBA(0, 0, 0, 0.6)
	if err := c.DrawWith(text, x+704-360, y+150); err != nil {
		return nil, err
	}

	return c, nil
}
