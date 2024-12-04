package entity

import (
	contract2 "leveling/internal/contract"
	"leveling/internal/server/contract"
	"leveling/internal/server/repository/dao"
	"leveling/internal/server/service"
	"leveling/internal/server/weapons"
)

const (
	roundTimeSecond = 1
	globalCoolDown  = 1.5
)

type Hero struct {
	name           string
	health         int
	strength       int
	mainHand       *contract.IWeapon
	attackCooldown float64 // weapon auto attack cooldown
	roundCooldown  float64
	client         *contract.Client
	nextAction     *contract2.ActionEvent
	target         *contract.IHero
	round          *contract.Round
	subject        *contract.Subject
	isActive       bool
}

func New(data dao.Hero, client *contract.Client) *contract.IHero {
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
	hero.isActive = false
	if hero.IsDie() {
		return false
	}
	hero.roundAutoAttack(dt)
	hero.roundAction(dt)

	return hero.isActive
}

func (hero *Hero) roundAutoAttack(dt float64) {
	weapon := *hero.mainHand
	hero.attackCooldown += dt / weapon.GetSpeed()
	for {
		if hero.attackCooldown < roundTimeSecond {
			return
		}
		hero.attackCooldown -= roundTimeSecond
		hero.doAutoAttack()
	}
}

func (hero *Hero) doAutoAttack() {
	if hero.target == nil {
		return
	}
	target := (*hero.target).(*Hero)
	if target.IsDie() {
		hero.target = nil
		return
	}
	hero.isActive = true
	damage := contract.Damage((*hero.mainHand).GetPower())
	target.ApplyDamage(damage)
	// send message to client
	messageEvent(hero, damage, target)
}

func (hero *Hero) roundAction(dt float64) {
	hero.roundCooldown += dt
	for {
		if hero.roundCooldown < globalCoolDown {
			return
		}
		if hero.nextAction == nil {
			return
		}
		hero.roundCooldown -= globalCoolDown
		hero.doAction()
	}
}

func (hero *Hero) doAction() {
	defer func() {
		hero.nextAction = nil
	}()

	if hero.target == nil {
		return
	}
	target := (*hero.target).(*Hero)
	if target.IsDie() {
		hero.target = nil
		return
	}
	hero.isActive = true
	damage := contract.Damage((*hero.mainHand).GetPower() + hero.strength)
	target.ApplyDamage(damage)
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
	}
	if to.subject != nil {
		(*to.subject).Notify(to, getHurtEvent)
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
	if hero.health <= 0 {
		hero.health = 0
	}
}

func (hero *Hero) SetSubject(subject *contract.Subject) {
	hero.subject = subject
}

func (hero *Hero) Subject() contract.Subject {
	return *hero.subject
}
