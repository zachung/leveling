package hero

import (
	"fmt"
	contract2 "leveling/internal/contract"
	"leveling/internal/entity"
	"leveling/internal/server/contract"
	"leveling/internal/server/service"
	"leveling/internal/weapons"
	"math"
)

const ROUNT_TIME_SECOND = 1

type Hero struct {
	name          string
	health        int
	strength      int
	mainHand      *contract.IWeapon
	roundCooldown float64 // weapon auto attack cooldown
	client        *contract.Client
	nextAction    *contract2.Action
}

func New(data entity.Hero, client *contract.Client) *contract.IHero {
	weapon := weapons.NewWeapon(data.MainHand)
	hero := &Hero{
		name:          data.Name,
		health:        data.Health,
		strength:      data.Strength,
		mainHand:      &weapon,
		roundCooldown: 0,
		client:        client,
	}
	iHero := contract.IHero(hero)
	weapon.SetHolder(&iHero)

	return &iHero
}

func (hero *Hero) Attack(dt float64, targets []*contract.IHero) {
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

func (hero *Hero) ApplyDamage(from *contract.IHero, power int) {
	attacker := (*from).(*Hero)
	damage := power + attacker.strength
	health := hero.health
	hero.health -= damage
	// send message to client
	// TODO: message type for show
	// TODO: event queue
	if hero.client != nil {
		message := fmt.Sprintf("[red]%s(%v)[white] take [red]%v damage[white] attacked by [::u]%s[::U]", hero.name, health, damage, attacker.name)
		if hero.IsDie() {
			message = message + fmt.Sprintf(", %v is Died", hero.name)
		}
		client := *hero.client
		client.Send([]byte(message))
	} else {
		message := fmt.Sprintf("%s(%v) take %v damage attacked by %s", hero.name, health, damage, attacker.name)
		if hero.IsDie() {
			message = message + fmt.Sprintf(", %v is Died", hero.name)
		}
		service.Logger().Debug("%v\n", message)
	}
}

func (hero *Hero) IsDie() bool {
	return hero.health <= 0
}

func (hero *Hero) SetNextAction(action *contract2.Action) {
	hero.nextAction = action
	service.Logger().Info("%s %v\n", hero.name, string(action.Serialize()))
}
