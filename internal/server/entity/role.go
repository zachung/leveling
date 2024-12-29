package entity

import (
	"golang.org/x/image/math/f64"
	contract2 "leveling/internal/contract"
	"leveling/internal/server/contract"
	"leveling/internal/server/repository/dao"
	"leveling/internal/server/weapons"
)

const (
	roundTimeSecond = 1
	globalCoolDown  = 1.5
)

type Hero struct {
	name          string
	health        int
	strength      int
	mainHand      contract.IWeapon
	roundCooldown float64
	client        contract.Client
	target        contract.IHero
	round         contract.Round
	subject       contract.Subject
	isActive      bool

	// operations
	operations map[contract2.RoleEvent][]func(...any)

	// abilities
	abilities map[AbilityType]Ability

	// move
	position f64.Vec2
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
		position:      data.Position,
		operations:    make(map[contract2.RoleEvent][]func(...any)),
		abilities:     make(map[AbilityType]Ability),
	}
	// TODO: hero.position & vector move to MoveAbility
	hero.abilities[AutoAttack] = NewAutoAttackAbility(hero)
	hero.abilities[Action] = NewActionAbility(hero)
	hero.abilities[Movement] = NewMoveAbility(hero)

	return hero
}

func (hero *Hero) Update(dt float64) bool {
	if hero.IsDie() {
		return false
	}
	for _, a := range hero.abilities {
		a.Update(dt)
	}

	isActive := hero.isActive
	hero.isActive = false

	return isActive
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

func (hero *Hero) SetAction(action contract2.Message) {
	switch action.(type) {
	case contract2.MoveEvent:
		event := action.(contract2.MoveEvent)
		// TODO: Right 只是暫時寫的
		for _, f := range hero.operations[contract2.Right] {
			f(event)
		}
	case contract2.ActionEvent:
		event := action.(contract2.ActionEvent)
		for _, f := range hero.operations[event.Id] {
			f(!event.IsEnable)
		}
		switch event.Id {
		case contract2.Skill1:
			if hero.target == nil {
				return
			}
			hero.abilities[Action].(*ActionAbility).nextAction = &event
		case contract2.CancelAction:
			hero.abilities[Action].(*ActionAbility).nextAction = nil
			hero.target = nil
		default:
			return
		}
	}
	hero.isActive = true
	hero.subject.Notify(hero.getCurrentState())
}

func (hero *Hero) GetName() string {
	return hero.name
}

func (hero *Hero) GetHealth() int {
	return hero.health
}

func (hero *Hero) SetTarget(name string) {
	hero.target = hero.round.GetHero(name)
	hero.subject.Notify(hero.getCurrentState())
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
		IsAutoAttack: hero.abilities[AutoAttack].(*AutoAttackAbility).isAutoAttack,
		Position:     hero.position,
		Vector:       hero.abilities[Movement].(*MoveAbility).vector,
	}
	action := hero.abilities[Action].(*ActionAbility).nextAction
	if action != nil {
		event.Action = *action
	}
	if hero.target != nil {
		event.Target = contract2.Hero{
			Name:   hero.target.GetName(),
			Health: hero.target.GetHealth(),
		}
	}

	return event
}

func (hero *Hero) GetState() contract2.Hero {
	var h contract2.Hero
	h.Name = hero.name
	h.Health = hero.health
	h.Position = hero.position
	h.Vector = hero.abilities[Movement].(*MoveAbility).vector
	if hero.target != nil {
		h.Target = &contract2.Hero{
			Name:   hero.target.GetName(),
			Health: hero.target.GetHealth(),
		}
	}

	return h
}

func (hero *Hero) AddOperationListener(k contract2.RoleEvent, f func(...any)) {
	hero.operations[k] = append(hero.operations[k], f)
}
