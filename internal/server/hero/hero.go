package hero

import (
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

func (hero *Hero) Update(dt float64) bool {
	if hero.IsDie() {
		return false
	}
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
	if target.IsDie() {
		hero.target = nil
		return
	}
	damage := contract.Damage((*hero.mainHand).GetPower() + hero.strength)
	target.ApplyDamage(damage)
	if target.health <= 0 {
		target.health = 0
	}
	// send message to client
	messageEvent(hero, damage, target)
}

func messageEvent(from *Hero, damage contract.Damage, to *Hero) {
	// TODO: event queue
	getHurtEvent := contract2.StateChangeEvent{
		Event: contract2.Event{
			Type: contract2.StateChange,
		},
		Name:         to.name,
		Health:       to.health,
		Damage:       int(damage),
		AttackerName: from.name,
	}
	if from.client != nil {
		(*from.client).Send(getHurtEvent)
	}
	if to.client != nil {
		(*to.client).Send(getHurtEvent)
	} else {
		// TODO: 分離 system entity 操作
		hero := contract.IHero(from)
		to.target = &hero
		event := contract2.ActionEvent{
			Event: contract2.Event{
				Type: contract2.Action,
			},
			Id: 1,
		}
		to.SetNextAction(&event)
	}
	if to.IsDie() {
		dieEvent := contract2.HeroDieEvent{Event: contract2.Event{Type: contract2.HeroDie}, Name: to.name}
		if from.client != nil {
			(*from.client).Send(dieEvent)
		}
		if to.client != nil {
			(*to.client).Send(dieEvent)
		}
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

func (hero *Hero) ApplyDamage(damage contract.Damage) {
	hero.health -= int(damage)
}
