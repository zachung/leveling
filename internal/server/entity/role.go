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
	mainHand       contract.IWeapon
	attackCooldown float64 // weapon auto attack cooldown
	roundCooldown  float64
	client         contract.Client
	nextAction     *contract2.ActionEvent
	target         contract.IHero
	round          contract.Round
	subject        contract.Subject
	isActive       bool
	isAutoAttack   bool
}

func NewRole(data dao.Hero, subject contract.Subject, client contract.Client) contract.IHero {
	weapon := weapons.NewWeapon(data.MainHand)
	hero := &Hero{
		name:          data.Name,
		health:        data.Health,
		strength:      data.Strength,
		mainHand:      weapon,
		roundCooldown: 0,
		client:        client,
		subject:       subject,
	}

	return hero
}

func (hero *Hero) Update(dt float64) bool {
	if hero.IsDie() {
		return false
	}
	hero.roundAutoAttack(dt)
	hero.roundAction(dt)
	isActive := hero.isActive
	hero.isActive = false

	return isActive
}

func (hero *Hero) roundAutoAttack(dt float64) {
	hero.attackCooldown += dt / hero.mainHand.GetSpeed()
	for {
		if hero.attackCooldown < roundTimeSecond {
			return
		}
		hero.attackCooldown -= roundTimeSecond
		hero.doAutoAttack()
	}
}

func (hero *Hero) doAutoAttack() {
	if !hero.isAutoAttack {
		return
	}
	if hero.target == nil {
		return
	}
	target := hero.target.(*Hero)
	if target.IsDie() {
		hero.isAutoAttack = false
		hero.target = nil
		return
	}
	hero.isActive = true
	damage := contract.Damage(hero.mainHand.GetPower())
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
	target := hero.target.(*Hero)
	if target.IsDie() {
		hero.target = nil
		return
	}
	hero.isActive = true
	damage := contract.Damage(hero.mainHand.GetPower() + hero.strength)
	target.ApplyDamage(damage)
	// send message to client
	messageEvent(hero, damage, target)
}

func messageEvent(from *Hero, damage contract.Damage, to *Hero) {
	// TODO: event queue
	// make damage event
	makeDamageEvent := contract2.MakeDamageEvent{Event: contract2.Event{Type: contract2.MakeDamage}}
	makeDamageEvent.To = contract2.Hero{Name: to.name, Health: to.health}
	makeDamageEvent.Damage = int(damage)
	from.subject.Notify(makeDamageEvent)
	// get hurt event
	getHurtEvent := contract2.GetHurtEvent{Event: contract2.Event{Type: contract2.GetHurt}}
	getHurtEvent.From = contract2.Hero{Name: from.name, Health: from.health}
	getHurtEvent.Damage = int(damage)
	to.subject.Notify(getHurtEvent)
	from.subject.Notify(from.getCurrentState())
	to.subject.Notify(to.getCurrentState())
	// die event
	if to.IsDie() {
		dieEvent := contract2.HeroDieEvent{Event: contract2.Event{Type: contract2.HeroDie}, Name: to.name}
		from.subject.Notify(dieEvent)
		to.subject.Notify(dieEvent)
	}
}

func (hero *Hero) IsDie() bool {
	return hero.health <= 0
}

func (hero *Hero) SetNextAction(action *contract2.ActionEvent) {
	switch action.Id {
	case 1:
		hero.isAutoAttack = !hero.isAutoAttack
		event := hero.getCurrentState()
		event.Action = *action
		if hero.client != nil {
			hero.client.Send(event)
		}
	case 2:
		hero.nextAction = action
		event := hero.getCurrentState()
		event.Action = *action
		if hero.client != nil {
			hero.client.Send(event)
		}
		service.Logger().Debug("%s %+v\n", hero.name, action)
	}
}

func (hero *Hero) GetName() string {
	return hero.name
}

func (hero *Hero) GetHealth() int {
	return hero.health
}

func (hero *Hero) SetTarget(name string) {
	hero.target = hero.round.GetHero(name)
}

func (hero *Hero) GetTarget() contract.IHero {
	return hero.target
}

func (hero *Hero) SetRound(round contract.Round) {
	hero.round = round
}

func (hero *Hero) ApplyDamage(damage contract.Damage) {
	hero.health -= int(damage)
	if hero.health <= 0 {
		hero.health = 0
	}
}

func (hero *Hero) SetSubject(subject contract.Subject) {
	hero.subject = subject
}

func (hero *Hero) Subject() contract.Subject {
	return hero.subject
}

func (hero *Hero) getCurrentState() contract2.StateChangeEvent {
	event := contract2.StateChangeEvent{
		Event: contract2.Event{
			Type: contract2.StateChange,
		},
		Name:         hero.name,
		Health:       hero.health,
		IsAutoAttack: hero.isAutoAttack,
	}
	if hero.target != nil {
		event.Target = contract2.Hero{
			Name:   hero.target.GetName(),
			Health: hero.target.GetHealth(),
		}
	}

	return event
}

func (hero *Hero) SetAutoAttack(isAutoAttack bool) {
	hero.isAutoAttack = isAutoAttack
}
