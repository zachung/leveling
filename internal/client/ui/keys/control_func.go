package keys

import (
	"github.com/hajimehoshi/ebiten/v2"
	"leveling/internal/contract"
)

var KeyMap = map[ebiten.Key]contract.RoleEvent{
	ebiten.Key1:      contract.SetAutoAttack,
	ebiten.Key2:      contract.Skill1,
	ebiten.Key3:      contract.Skill2,
	ebiten.KeyW:      contract.Up,
	ebiten.KeyA:      contract.Left,
	ebiten.KeyS:      contract.Down,
	ebiten.KeyD:      contract.Right,
	ebiten.KeyEscape: contract.CancelAction,
}
