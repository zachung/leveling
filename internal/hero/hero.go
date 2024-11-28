package hero

import (
	"fmt"
	"leveling/internal/entity"
	"leveling/internal/server/constract"
	"leveling/internal/server/service"
	"leveling/internal/weapons"
	"math"
)

const ROUNT_TIME_SECOND = 1

type Hero struct {
	name          string
	health        int
	strength      int
	mainHand      *constract.IWeapon
	roundCooldown float64 // weapon auto attack cooldown
}

func New(data entity.Hero) *constract.IHero {
	weapon := weapons.NewWeapon(data.MainHand)
	hero := &Hero{
		name:          data.Name,
		health:        data.Health,
		strength:      data.Strength,
		mainHand:      &weapon,
		roundCooldown: 0,
	}
	iHero := constract.IHero(hero)
	weapon.SetHolder(&iHero)

	return &iHero
}

func (hero *Hero) Attack(dt float64, targets []*constract.IHero) {
	weapon := *hero.mainHand
	hero.roundCooldown += dt / weapon.GetSpeed()
	if hero.roundCooldown < ROUNT_TIME_SECOND {
		return
	}
	roundTime := hero.roundCooldown
	for rounds := int64(roundTime / ROUNT_TIME_SECOND); rounds > 0; rounds-- {
		weapon.Attack(targets[0])
	}
	hero.roundCooldown = math.Mod(roundTime, ROUNT_TIME_SECOND)
}

func (hero *Hero) ApplyDamage(from *constract.IHero, power int) {
	attacker := (*from).(*Hero)
	damage := power + attacker.strength
	message := fmt.Sprintf("%s(%v) take %v damage attacked by %s", hero.name, hero.health, damage, attacker.name)
	hero.health -= damage
	if hero.IsDie() {
		message = message + fmt.Sprintf(", %v is Died", hero.name)
	}
	// TODO: send message to client
	service.Logger().Info("%v\n", message)
}

func (hero *Hero) IsDie() bool {
	return hero.health <= 0
}
