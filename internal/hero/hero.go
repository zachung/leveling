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
	if hero.nextAction == nil {
		// 下次可以直接動作
		hero.roundCooldown = ROUNT_TIME_SECOND
	} else {
		for rounds := int64(roundTime / ROUNT_TIME_SECOND); rounds > 0; rounds-- {
			weapon.Attack(targets[0])
		}
		hero.nextAction = nil
		hero.roundCooldown = math.Mod(roundTime, ROUNT_TIME_SECOND)
	}
}

func (hero *Hero) ApplyDamage(from *contract.IHero, power int) {
	attacker := (*from).(*Hero)
	damage := power + attacker.strength
	health := hero.health
	hero.health -= damage
	// send message to client
	messageEvent(hero, health, damage, attacker)
}

func messageEvent(hero *Hero, health int, damage int, attacker *Hero) {
	// TODO: message type for show
	// TODO: event queue
	var message string
	// display for applied
	message = fmt.Sprintf("[red]-%v health[white] from [::u]%s[::U], remain %v", damage, attacker.name, health)
	if hero.client != nil {
		client := *hero.client
		client.Send([]byte(message))
		if hero.IsDie() {
			client.Send([]byte("You Died"))
		}
	} else {
		service.Logger().Debug(message)
	}
	// display for attacker
	message = fmt.Sprintf("attack [red]%s(%v)[white] make [red]%v[white] damage", hero.name, health, damage)
	if attacker.client != nil {
		client := *attacker.client
		client.Send([]byte(message))
		if hero.IsDie() {
			client.Send([]byte(fmt.Sprintf("%v is Died", hero.name)))
		}
	} else {
		service.Logger().Debug(message)
	}
}

func (hero *Hero) IsDie() bool {
	return hero.health <= 0
}

func (hero *Hero) SetNextAction(action *contract2.Action) {
	hero.nextAction = action
	service.Logger().Debug("%s %v\n", hero.name, string(action.Serialize()))
}
