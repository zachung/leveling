package object

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/math/f64"
	"image/color"
	contract2 "leveling/internal/client/contract"
	"leveling/internal/contract"
)

type Role struct {
	Position f64.Vec2
	hero     contract.Hero
}

var rectImage *ebiten.Image

func init() {
	// 創建矩形圖像
	rectImage = ebiten.NewImage(20, 20)
	rectImage.Fill(color.RGBA{255, 0, 0, 255}) // 填充紅色
}

func NewRole(hero contract.Hero) *Role {
	return &Role{
		hero: hero,
	}
}

func (r *Role) Draw(screen *ebiten.Image) {
	// 設置 GeoM 變換
	op := &ebiten.DrawImageOptions{}
	// 角色物件們位移到畫面中間
	x := r.Position[0] + contract2.ScreenWidth/2
	y := r.Position[1] + contract2.ScreenHeight/2
	op.GeoM.Translate(x, y)

	// 渲染矩形到螢幕
	screen.DrawImage(rectImage, op)
	// Name
	textOp := &text.DrawOptions{}
	textOp.GeoM.Translate(x, y-20)
	textOp.ColorScale.ScaleWithColor(color.White)
	str := fmt.Sprintf("%s(%d)", r.hero.Name, r.hero.Health)
	text.Draw(screen, str, &text.GoTextFace{
		Source: contract2.UiFaceSource,
		Size:   contract2.NormalFontSize,
	}, textOp)
}
