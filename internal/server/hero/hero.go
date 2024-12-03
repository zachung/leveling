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
	target        *contract.IHero
	round         *contract.Round
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

func (hero *Hero) Attack(dt float64) bool {
	weapon := *hero.mainHand
	hero.roundCooldown += dt / weapon.GetSpeed()
	if hero.roundCooldown < ROUNT_TIME_SECOND {
		return false
	}
	roundTime := hero.roundCooldown
	if hero.nextAction == nil {
		// 下次可以直接動作
		hero.roundCooldown = ROUNT_TIME_SECOND

		return false
	} else {
		for rounds := int64(roundTime / ROUNT_TIME_SECOND); rounds > 0; rounds-- {
			hero.attackTarget()
		}
		hero.nextAction = nil
		hero.roundCooldown = math.Mod(roundTime, ROUNT_TIME_SECOND)

		return true
	}
}

func (hero *Hero) attackTarget() {
	if hero.target == nil {
		return
	}
	target := (*hero.target).(*Hero)
	damage := (*hero.mainHand).GetPower() + hero.strength
	target.health -= damage
	// send message to client
	messageEvent(target, damage, hero)
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

func (hero *Hero) SetTarget(name string) {
	hero.target = (*hero.round).GetHero(name)
}

func (hero *Hero) SetRound(round *contract.Round) {
	hero.round = round
}
