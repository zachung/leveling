package hero

import (
	"fmt"
	contract2 "leveling/internal/contract"
	"leveling/internal/server/contract"
	"leveling/internal/server/entity"
	"leveling/internal/server/service"
	"leveling/internal/server/weapons"
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
	nextAction    *contract2.ActionEvent
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
	hero.health -= damage
	// send message to client
	messageEvent(hero, damage, attacker)
}

func messageEvent(hero *Hero, damage int, attacker *Hero) {
	// TODO: event queue
	var message string
	// display for applied
	getHurtEvent := contract2.StateChangeEvent{
		Event: contract2.Event{
			Type: contract2.StateChange,
		},
		Name:         hero.name,
		Health:       hero.health,
		Damage:       damage,
		AttackerName: attacker.name,
	}
	if hero.client != nil {
		fromClient := *attacker.client
		toClient := *hero.client
		fromClient.Send(getHurtEvent)
		toClient.Send(getHurtEvent)
		if hero.IsDie() {
			dieEvent := contract2.HeroDieEvent{Event: contract2.Event{Type: contract2.HeroDie}, Name: hero.name}
			fromClient.Send(dieEvent)
			toClient.Send(dieEvent)
		}
	} else {
		message = fmt.Sprintf("-%v health from %s, remain %v", getHurtEvent.Damage, getHurtEvent.AttackerName, getHurtEvent.Health)
		service.Logger().Debug(message)
	}
}

func (hero *Hero) IsDie() bool {
	return hero.health <= 0
}

func (hero *Hero) SetNextAction(action *contract2.ActionEvent) {
	hero.nextAction = action
	service.Logger().Debug("%s %+v\n", hero.name, action)
}

func (hero *Hero) GetName() string {
	return hero.name
}

func (hero *Hero) GetHealth() int {
	return hero.health
}
