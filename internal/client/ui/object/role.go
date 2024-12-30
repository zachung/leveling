package object

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/math/f64"
	"image/color"
	contract2 "leveling/internal/client/contract"
)

type Role struct {
	Position f64.Vec2
	image    *ebiten.Image
	name     string
}

var rectImage *ebiten.Image

func init() {
	// 創建矩形圖像
	rectImage = ebiten.NewImage(20, 20)
	rectImage.Fill(color.RGBA{255, 0, 0, 255}) // 填充紅色
}

func NewRole(name string) *Role {
	// TODO: draw name
	return &Role{
		name:  name,
		image: rectImage,
	}
}

func (r *Role) Draw(screen *ebiten.Image) {
	// 設置 GeoM 變換
	op := &ebiten.DrawImageOptions{}
	x := r.Position[0]
	y := r.Position[1]
	op.GeoM.Translate(x, y)

	// 渲染矩形到螢幕
	screen.DrawImage(r.image, op)
	// Name
	textOp := &text.DrawOptions{}
	textOp.GeoM.Translate(x, y-20)
	textOp.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, r.name, &text.GoTextFace{
		Source: contract2.UiFaceSource,
		Size:   contract2.NormalFontSize,
	}, textOp)
}
