package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// AttackMove 定義了一種攻擊手段
type AttackMove struct {
	Name           string        // 招式名稱
	EnergyCost     int           // 消耗能量
	Damage         int           // 造成傷害 (削減生命)
	HealAmount     int           // 治療量
	PreCastTime    time.Duration // 施法前搖
	CastTime       time.Duration // 施法/引導時間
	ActionCooldown time.Duration // 動作冷卻
	IsAoE          bool          // 是否為範圍攻擊
	IsChanneling   bool          // 是否為引導技能

	// 模組化技能邏輯
	CanUse      func(caster *Player, target *Player) bool
	ApplyEffect func(caster *Player, target *Player, b *Battle) []string
}

// Shell 代表一個可被附身的軀殼
type Shell struct {
	Name      string
	Health    int
	MaxHealth int
	Strength  int // 力量屬性，影響傷害
	Skills    []*AttackMove
}

// Player 代表玩家或敵人的核心靈魂
type Player struct {
	Name         string
	Energy       int
	MaxEnergy    int
	Skills       []*AttackMove // 靈魂自身的技能
	CurrentShell *Shell        // 當前附身的軀殼，可能為 nil

	// 動作狀態機 (玩家與敵人都使用)
	ActionState     string
	ActionStartTime time.Time
	EffectTime      time.Time // 施法/引導結束時間
	EffectApplied   bool
	CastingMove     *AttackMove
	CastingTarget   *Player // 施法鎖定的目標
	LastChannelTick time.Time
	SkillCooldowns  map[string]time.Time // 技能冷卻計時器

	// 史萊姆專用屬性
	IsSlime          bool
	ReplicationCount int
	MaxReplications  int
}

// PendingSpawn 用於處理延遲生成的單位
type PendingSpawn struct {
	SpawnTime        time.Time
	Name             string
	Energy           int
	Shell            *Shell
	IsSlime          bool
	ReplicationCount int
	MaxReplications  int
}

// Battle 結構封裝了一場戰鬥的所有狀態和 UI
type Battle struct {
	app *tview.Application

	player  *Player
	enemies []*Player

	currentTargetIndex  int
	logHistory          []string
	pendingSpawns       []PendingSpawn
	gameIsOver          bool
	interruptChanneling bool

	// 預約的行動
	nextPlayerAction *AttackMove

	// UI 元件
	playerStatus *tview.TextView
	targetStatus *tview.TextView
	enemyList    *tview.List
	battleLog    *tview.TextView
	instructions *tview.TextView
	mainLayout   *tview.Flex
}

// 全域技能定義，以便動態切換
var (
	meditate     *AttackMove
	soulEjection *AttackMove
	possess      *AttackMove
)

// NewPlayer 創建一個新的靈魂實例
func NewPlayer(name string, energy int, isSlime bool, replications int, skills []*AttackMove) *Player {
	return &Player{
		Name:             name,
		Energy:           energy,
		MaxEnergy:        energy,
		ActionState:      "Idle",
		SkillCooldowns:   make(map[string]time.Time),
		IsSlime:          isSlime,
		ReplicationCount: replications,
		MaxReplications:  2,
		Skills:           skills,
	}
}

// NewShell 創建一個新的軀殼實例
func NewShell(name string, health int, strength int, skills []*AttackMove) *Shell {
	return &Shell{
		Name:      name,
		Health:    health,
		MaxHealth: health,
		Strength:  strength,
		Skills:    skills,
	}
}

// LoseHealth 減少軀殼的生命
func (s *Shell) LoseHealth(amount int) {
	s.Health -= amount
	if s.Health < 0 {
		s.Health = 0
	}
}

// Heal 恢復軀殼的生命
func (s *Shell) Heal(amount int) {
	s.Health += amount
	if s.Health > s.MaxHealth {
		s.Health = s.MaxHealth
	}
}

// IsDefeated 檢查軀殼是否被摧毀
func (s *Shell) IsDefeated() bool {
	return s.Health <= 0
}

// GainEnergy 為玩家增加能量
func (p *Player) GainEnergy(amount int) {
	p.Energy += amount
	if p.Energy > p.MaxEnergy {
		p.Energy = p.MaxEnergy
	}
}

// LoseEnergy 減少玩家的能量
func (p *Player) LoseEnergy(amount int) {
	p.Energy -= amount
	if p.Energy < 0 {
		p.Energy = 0
	}
}

// Attack 讓玩家驅動軀殼攻擊目標
func (p *Player) Attack(mainTarget *Player, allPossibleTargets []*Player, move *AttackMove) []string {
	if p.CurrentShell == nil {
		return []string{"靈體狀態無法攻擊！"}
	}

	p.LoseEnergy(move.EnergyCost)
	logs := []string{
		fmt.Sprintf("➡️ %s 的 [%s] 擊中了目標！", p.Name, move.Name),
		fmt.Sprintf("   %s 消耗了 %d 點能量。", p.Name, move.EnergyCost),
	}

	finalDamage := move.Damage + p.CurrentShell.Strength

	if move.IsAoE {
		logs = append(logs, "   這是一個範圍攻擊！")
		for _, target := range allPossibleTargets {
			if target != p && target.CurrentShell != nil && !target.CurrentShell.IsDefeated() {
				target.CurrentShell.LoseHealth(finalDamage)
				logs = append(logs, fmt.Sprintf("   對 %s 的軀殼造成了 %d 點傷害！", target.Name, finalDamage))
			}
		}
	} else {
		if mainTarget.CurrentShell == nil || mainTarget.CurrentShell.IsDefeated() {
			return []string{fmt.Sprintf("目標 %s 的軀殼已無靈魂，攻擊無效。", mainTarget.Name)}
		}
		mainTarget.CurrentShell.LoseHealth(finalDamage)
		logs = append(logs, fmt.Sprintf("   對 %s 的軀殼造成了 %d 點傷害！ (%d 基礎 + %d 力量)", mainTarget.Name, finalDamage, move.Damage, p.CurrentShell.Strength))
	}

	return logs
}

// StartAction 開始一個新的動作 (施法或引導)
func (p *Player) StartAction(skill *AttackMove, target *Player, now time.Time) {
	p.ActionState = "Casting"
	p.CastingMove = skill
	p.CastingTarget = target
	p.ActionStartTime = now
	p.EffectTime = now.Add(skill.PreCastTime)
	p.EffectApplied = false
	if skill.IsChanneling && skill.Name == "冥想" {
		p.ActionState = "Channeling"
		p.LastChannelTick = now
	}
}

// InterruptAction 中斷當前動作
func (p *Player) InterruptAction(now time.Time) []string {
	var logs []string
	if p.ActionState == "Casting" || p.ActionState == "Channeling" {
		logs = append(logs, fmt.Sprintf("[orange]%s 中斷了 [%s]。[-:-:-]", p.Name, p.CastingMove.Name))
		// 如果是引導技能，中斷時進入冷卻
		if p.ActionState == "Channeling" {
			p.SkillCooldowns[p.CastingMove.Name] = now.Add(p.CastingMove.ActionCooldown)
		}
		p.ActionState = "Idle"
	}
	return logs
}

// Update 處理角色單個 tick 的所有邏輯
func (p *Player) Update(now time.Time, b *Battle, isPlayerControlled bool) ([]string, bool) {
	var logsThisTick []string
	var actionTaken bool

	// --- 狀態機邏輯 ---
	switch p.ActionState {
	case "Idle":
		if isPlayerControlled {
			if b.nextPlayerAction != nil {
				p.StartAction(b.nextPlayerAction, p.CastingTarget, now)
				b.nextPlayerAction = nil
				actionTaken = true
			}
		} else { // AI 邏輯
			if p.CurrentShell != nil && !p.CurrentShell.IsDefeated() && b.player.CurrentShell != nil {
				if len(p.CurrentShell.Skills) > 0 {
					skill := p.CurrentShell.Skills[0]
					if now.After(p.SkillCooldowns[skill.Name]) {
						p.StartAction(skill, b.player, now)
						actionTaken = true
					}
				}
			}
			// 未來可擴充 AI 在靈體狀態下的行為
		}

	case "Casting":
		if !p.EffectApplied && now.After(p.EffectTime) {
			if p.CastingMove.IsChanneling {
				p.ActionState = "Channeling"
				p.ActionStartTime = now
				p.EffectTime = now.Add(p.CastingMove.CastTime)
				p.LastChannelTick = now
				p.LoseEnergy(p.CastingMove.EnergyCost)
				actionTaken = true
			} else {
				if p.CastingMove.ApplyEffect != nil {
					logsThisTick = append(logsThisTick, p.CastingMove.ApplyEffect(p, p.CastingTarget, b)...)
				}
				p.SkillCooldowns[p.CastingMove.Name] = now.Add(p.CastingMove.ActionCooldown)
				p.ActionState = "Idle"
				actionTaken = true
			}
		}

	case "Channeling":
		switch p.CastingMove.Name {
		case "冥想":
			if now.Sub(p.LastChannelTick) >= 500*time.Millisecond {
				p.GainEnergy(5)
				p.LastChannelTick = now
				actionTaken = true
			}
		case "踐踏":
			if now.After(p.EffectTime) {
				p.SkillCooldowns[p.CastingMove.Name] = now.Add(p.CastingMove.ActionCooldown)
				p.ActionState = "Idle"
				actionTaken = true
			} else if now.Sub(p.LastChannelTick) >= 500*time.Millisecond {
				logsThisTick = append(logsThisTick, "[yellow]踐踏造成了範圍傷害！[-:-:-]")
				finalDamage := p.CastingMove.Damage + p.CurrentShell.Strength
				for _, t := range b.enemies {
					if t.CurrentShell != nil && !t.CurrentShell.IsDefeated() {
						t.CurrentShell.LoseHealth(finalDamage)
						logsThisTick = append(logsThisTick, fmt.Sprintf("   對 %s 的軀殼造成了 %d 點傷害！", t.Name, finalDamage))
					}
				}
				p.LastChannelTick = now
				actionTaken = true
			}
		}
	}

	// --- 被動效果和狀態檢查 ---
	if p.CurrentShell == nil && p.ActionState == "Idle" {
		p.GainEnergy(1) // 靈體狀態下回能
	}
	if p.CurrentShell != nil && p.CurrentShell.IsDefeated() {
		logsThisTick = append(logsThisTick, fmt.Sprintf("[orange]%s 的軀殼被摧毀了！[-:-:-]", p.Name))
		p.CurrentShell = nil
		p.ActionState = "Idle"
		actionTaken = true
	}

	return logsThisTick, actionTaken
}

// createProgressBar 創建一個文字進度條
func createProgressBar(startTime, endTime time.Time, width int, color string) string {
	totalDuration := endTime.Sub(startTime)
	if totalDuration <= 0 {
		return strings.Repeat(" ", width+2)
	}
	elapsed := time.Since(startTime)
	progress := float64(elapsed) / float64(totalDuration)
	if progress < 0 {
		progress = 0
	}
	if progress > 1 {
		progress = 1
	}

	filledWidth := int(progress * float64(width))

	bar := fmt.Sprintf("[%s]", color)
	bar += strings.Repeat("█", filledWidth)
	bar += "[gray]"
	bar += strings.Repeat("░", width-filledWidth)
	bar += "[-:-:-]"
	return bar
}

// createDecreasingProgressBar 創建一個由右至左減少的文字進度條
func createDecreasingProgressBar(startTime, endTime time.Time, width int, color string) string {
	totalDuration := endTime.Sub(startTime)
	if totalDuration <= 0 {
		return strings.Repeat(" ", width+2)
	}
	remaining := time.Until(endTime)
	progress := float64(remaining) / float64(totalDuration)
	if progress < 0 {
		progress = 0
	}
	if progress > 1 {
		progress = 1
	}

	filledWidth := int(progress * float64(width))

	bar := fmt.Sprintf("[%s]", color)
	bar += strings.Repeat("█", filledWidth)
	bar += "[gray]"
	bar += strings.Repeat("░", width-filledWidth)
	bar += "[-:-:-]"
	return bar
}

// createValueBar 創建一個基於數值的進度條
func createValueBar(current, max, width int, color string) string {
	if max <= 0 {
		return ""
	}
	progress := float64(current) / float64(max)
	if progress < 0 {
		progress = 0
	}
	if progress > 1 {
		progress = 1
	}

	filledWidth := int(progress * float64(width))

	bar := fmt.Sprintf("[%s]", color)
	bar += strings.Repeat("█", filledWidth)
	bar += "[gray]"
	bar += strings.Repeat("░", width-filledWidth)
	bar += "[-:-:-]"
	return bar
}

// getActiveSkills 返回玩家當前可用的所有技能列表
func (p *Player) getActiveSkills() []*AttackMove {
	var activeSkills []*AttackMove
	// 添加基礎靈魂技能
	for _, skill := range p.Skills {
		if skill.Name == "靈魂行動" {
			if p.CurrentShell == nil {
				activeSkills = append(activeSkills, possess)
			} else {
				activeSkills = append(activeSkills, soulEjection)
			}
		} else {
			activeSkills = append(activeSkills, skill)
		}
	}

	if p.CurrentShell != nil {
		activeSkills = append(activeSkills, p.CurrentShell.Skills...)
	}
	return activeSkills
}

// GetPlayerStatusText 獲取玩家狀態文字
func (p *Player) GetPlayerStatusText(b *Battle) string {
	var status strings.Builder
	status.WriteString(fmt.Sprintf("[::b]%s\n", p.Name))
	status.WriteString(fmt.Sprintf("%s\n", strings.Repeat("─", len(p.Name)+4)))

	if p.CurrentShell != nil {
		status.WriteString(fmt.Sprintf("[red]生命: %d / %d[-:-:-]\n", p.CurrentShell.Health, p.CurrentShell.MaxHealth))
		status.WriteString(createValueBar(p.CurrentShell.Health, p.CurrentShell.MaxHealth, 20, "red") + "\n")
	}

	status.WriteString(fmt.Sprintf("[blue]能量: %d / %d[-:-:-]\n", p.Energy, p.MaxEnergy))
	status.WriteString(createValueBar(p.Energy, p.MaxEnergy, 20, "blue") + "\n")

	if p.CurrentShell != nil {
		status.WriteString(fmt.Sprintf("[orange]力量: %d[-:-:-]\n", p.CurrentShell.Strength))
	} else {
		status.WriteString("[purple]狀態: 靈體[-:-:-]\n")
	}

	// 技能列
	status.WriteString("\n" + strings.Repeat("─", 24) + "\n")
	now := time.Now()

	activeSkills := p.getActiveSkills()
	skillKeys := []rune{'q', 'w', 'e', 'r', 'a', 's', 'd', 'f'}

	for i, skill := range activeSkills {
		if i >= len(skillKeys) {
			break
		}

		key := skillKeys[i]
		leftPart := fmt.Sprintf("(%c) %-8s", key, skill.Name)
		var rightPart string

		isCastingThisSkill := p.ActionState != "Idle" && p.CastingMove != nil && p.CastingMove.Name == skill.Name
		if isCastingThisSkill {
			if p.ActionState == "Casting" {
				rightPart = createProgressBar(p.ActionStartTime, p.EffectTime, 10, "yellow")
			} else { // Channeling
				if skill.CastTime > 0 {
					rightPart = createDecreasingProgressBar(p.ActionStartTime, p.EffectTime, 10, "blue")
				} else {
					rightPart = "[blue]" + strings.Repeat("█", 10) + "[-:-:-]"
				}
			}
		} else {
			cd, onCD := p.SkillCooldowns[skill.Name]
			if onCD && now.Before(cd) {
				rightPart = createProgressBar(cd.Add(-skill.ActionCooldown), cd, 10, "red")
			} else {
				rightPart = "[green]準備就緒[-:-:-]"
			}
		}
		status.WriteString(fmt.Sprintf("%-14s %s\n", leftPart, rightPart))
	}

	return status.String()
}

// GetEnemyStatusText 獲取單一敵人狀態文字
func (p *Player) GetEnemyStatusText() string {
	if p.CurrentShell == nil {
		return fmt.Sprintf("[::b]%s\n\n[gray]靈體狀態[-:-:-]", p.Name)
	}
	var status strings.Builder
	shellName := strings.Split(p.CurrentShell.Name, " ")[0]
	if p.CurrentShell.IsDefeated() {
		status.WriteString(fmt.Sprintf("[::b]%s\n\n[purple]無主軀殼[-:-:-]", shellName))
	} else {
		status.WriteString(fmt.Sprintf("[::b]%s\n", p.Name))
		status.WriteString(fmt.Sprintf("%s\n", strings.Repeat("─", len(p.Name)+4)))
		status.WriteString(fmt.Sprintf("[red]生命: %d / %d[-:-:-]\n", p.CurrentShell.Health, p.CurrentShell.MaxHealth))
		status.WriteString(createValueBar(p.CurrentShell.Health, p.CurrentShell.MaxHealth, 20, "red") + "\n")
		status.WriteString(fmt.Sprintf("[orange]力量: %d[-:-:-]\n", p.CurrentShell.Strength))

		switch p.ActionState {
		case "Casting":
			status.WriteString(fmt.Sprintf("[yellow]施法中: %s (%.1fs)[-:-:-]\n", p.CastingMove.Name, time.Until(p.EffectTime).Seconds()))
			status.WriteString(createProgressBar(p.ActionStartTime, p.EffectTime, 20, "yellow"))
		case "Idle":
			status.WriteString("[green]狀態: 可行動[-:-:-]")
		}
	}
	return status.String()
}

// NewBattle 創建一個新的戰鬥實例
func NewBattle(app *tview.Application) *Battle {
	// --- 遊戲設定 ---
	// 軀殼技能
	slash := &AttackMove{Name: "揮砍", EnergyCost: 10, Damage: 15, PreCastTime: 0, ActionCooldown: 1 * time.Second}
	heavyStrike := &AttackMove{Name: "強力一擊", EnergyCost: 35, Damage: 80, PreCastTime: 1 * time.Second, ActionCooldown: 2 * time.Second}
	stomp := &AttackMove{Name: "踐踏", EnergyCost: 20, Damage: 4, PreCastTime: 1 * time.Second, CastTime: 2 * time.Second, ActionCooldown: 4 * time.Second, IsAoE: true, IsChanneling: true}
	bite := &AttackMove{Name: "啃咬", EnergyCost: 10, Damage: 12, PreCastTime: 0, ActionCooldown: 1 * time.Second}
	heal := &AttackMove{Name: "治療", EnergyCost: 25, HealAmount: 100, PreCastTime: 1 * time.Second, ActionCooldown: 1 * time.Second}

	// 靈魂技能 - 全域變數初始化
	meditate = &AttackMove{Name: "冥想", IsChanneling: true}
	soulEjection = &AttackMove{Name: "靈魂出竅", EnergyCost: 10, PreCastTime: 500 * time.Millisecond, ActionCooldown: 1 * time.Second}
	possess = &AttackMove{Name: "附身", EnergyCost: 60, PreCastTime: 2 * time.Second, ActionCooldown: 1 * time.Second}

	// 為技能綁定模組化函式
	slash.CanUse = func(c *Player, t *Player) bool {
		return c.CurrentShell != nil && t != nil && t.CurrentShell != nil && !t.CurrentShell.IsDefeated()
	}
	slash.ApplyEffect = func(c *Player, t *Player, b *Battle) []string { return c.Attack(t, b.enemies, c.CastingMove) }
	heavyStrike.CanUse = slash.CanUse // 邏輯相同
	heavyStrike.ApplyEffect = slash.ApplyEffect
	bite.CanUse = slash.CanUse
	bite.ApplyEffect = slash.ApplyEffect
	stomp.CanUse = func(c *Player, t *Player) bool { return c.CurrentShell != nil }   // AOE 不需要特定目標
	stomp.ApplyEffect = func(c *Player, t *Player, b *Battle) []string { return nil } // 引導技能在 gameLoop 中處理
	heal.CanUse = func(c *Player, t *Player) bool { return c.CurrentShell != nil }
	heal.ApplyEffect = func(c *Player, t *Player, b *Battle) []string {
		c.CurrentShell.Heal(c.CastingMove.HealAmount)
		c.LoseEnergy(c.CastingMove.EnergyCost)
		return []string{fmt.Sprintf("[green]你治療了自己 %d 點生命！[-:-:-]", c.CastingMove.HealAmount)}
	}
	meditate.CanUse = func(c *Player, t *Player) bool { return c.CurrentShell != nil }
	soulEjection.CanUse = func(c *Player, t *Player) bool { return c.CurrentShell != nil }
	soulEjection.ApplyEffect = func(c *Player, t *Player, b *Battle) []string {
		c.CurrentShell = nil
		c.LoseEnergy(c.CastingMove.EnergyCost)
		return []string{"[purple]你施展了靈魂出竅，脫離了當前的軀殼！[-:-:-]"}
	}
	possess.CanUse = func(c *Player, t *Player) bool {
		return c.CurrentShell == nil && t != nil && t.CurrentShell != nil && t.CurrentShell.IsDefeated()
	}
	possess.ApplyEffect = func(c *Player, t *Player, b *Battle) []string {
		logs := []string{}
		if t != nil && t.CurrentShell != nil && t.CurrentShell.IsDefeated() {
			logs = append(logs, fmt.Sprintf("[green]你以靈體狀態，成功附身到 %s 的軀殼上！[-:-:-]", t.Name))
			c.LoseEnergy(c.CastingMove.EnergyCost)
			t.CurrentShell.Health = t.CurrentShell.MaxHealth
			c.CurrentShell = t.CurrentShell
			t.CurrentShell = nil
		}
		return logs
	}

	player := NewPlayer("英雄", 100, false, 0, []*AttackMove{meditate, {Name: "靈魂行動"}}) // 使用一個佔位技能
	player.CurrentShell = NewShell("人類軀殼", 500, 5, []*AttackMove{slash, heavyStrike, stomp, bite, heal})

	enemies := []*Player{
		NewPlayer("哥布林", 999, false, 0, nil),
		NewPlayer("史萊姆", 999, true, 0, nil),
		NewPlayer("骷髏兵", 999, false, 0, nil),
	}
	enemies[0].CurrentShell = NewShell("哥布林軀殼", 80, 2, []*AttackMove{stomp})
	enemies[1].CurrentShell = NewShell("凝膠軀殼", 60, 5, []*AttackMove{bite})
	enemies[2].CurrentShell = NewShell("骸骨軀殼", 120, 8, []*AttackMove{stomp})

	b := &Battle{
		app:        app,
		player:     player,
		enemies:    enemies,
		logHistory: []string{"戰鬥開始！"},
	}

	b.playerStatus = tview.NewTextView()
	b.playerStatus.SetDynamicColors(true).SetTextAlign(tview.AlignLeft).SetBorder(true).SetTitle("你的狀態")
	b.targetStatus = tview.NewTextView()
	b.targetStatus.SetDynamicColors(true).SetTextAlign(tview.AlignCenter).SetBorder(true).SetTitle("鎖定目標")
	b.enemyList = tview.NewList()
	b.enemyList.ShowSecondaryText(false).SetBorder(true).SetTitle("敵人清單")
	b.battleLog = tview.NewTextView()
	b.battleLog.SetDynamicColors(true).SetScrollable(true).SetBorder(true).SetTitle("戰鬥日誌")
	b.instructions = tview.NewTextView()
	b.instructions.SetDynamicColors(true)

	return b
}

// updateAllViews 更新所有 UI 元件
func (b *Battle) updateAllViews() {
	b.playerStatus.SetText(b.player.GetPlayerStatusText(b))
	if b.currentTargetIndex < len(b.enemies) {
		b.targetStatus.SetText(b.enemies[b.currentTargetIndex].GetEnemyStatusText())
	}

	enemyListChanged := newEnemyListChanged(b)
	b.enemyList.SetChangedFunc(nil)
	b.enemyList.Clear()
	for i, enemy := range b.enemies {
		var status string
		if enemy.CurrentShell == nil {
			status = "[gray]靈體"
		} else if enemy.CurrentShell.IsDefeated() {
			status = "[purple]無主軀殼"
		} else {
			status = createValueBar(enemy.CurrentShell.Health, enemy.CurrentShell.MaxHealth, 10, "red")
		}

		prefix := "  "
		if i == b.currentTargetIndex {
			prefix = "[red]>>[-:-:-]"
		}
		mainText := fmt.Sprintf("%s %-12s %s", prefix, enemy.Name, status)

		b.enemyList.AddItem(mainText, "", 0, nil)
	}
	if b.currentTargetIndex < len(b.enemies) {
		b.enemyList.SetCurrentItem(b.currentTargetIndex)
	}
	b.enemyList.SetChangedFunc(enemyListChanged)

	b.instructions.SetText("(1-9)選敵 | (Tab/Shift+Tab)切換 | (Ctrl+C)離開")
}

// gameLoop 是戰鬥的主迴圈
func (b *Battle) gameLoop() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		if b.gameIsOver {
			continue
		}

		now := time.Now()
		var allLogsThisTick []string
		var anyActionTaken bool

		// 玩家中斷處理
		if b.interruptChanneling {
			logs, interrupted := b.player.InterruptAction(now), true
			allLogsThisTick = append(allLogsThisTick, logs...)
			anyActionTaken = anyActionTaken || interrupted
			b.interruptChanneling = false
		}

		// 更新玩家
		playerLogs, playerAction := b.player.Update(now, b, true)
		allLogsThisTick = append(allLogsThisTick, playerLogs...)
		anyActionTaken = anyActionTaken || playerAction

		// 更新敵人並檢查勝利條件
		allEnemiesDefeated := true
		for _, enemy := range b.enemies {
			if enemy.CurrentShell != nil {
				allEnemiesDefeated = false
			}

			enemyLogs, enemyAction := enemy.Update(now, b, false)
			allLogsThisTick = append(allLogsThisTick, enemyLogs...)
			anyActionTaken = anyActionTaken || enemyAction

			// 處理史萊姆分裂
			if enemy.CurrentShell != nil && enemy.CurrentShell.IsDefeated() && enemy.IsSlime && enemy.ReplicationCount < enemy.MaxReplications {
				newShell := *enemy.CurrentShell
				newShell.Health = newShell.MaxHealth / 2
				newShell.MaxHealth = newShell.MaxHealth / 2
				b.pendingSpawns = append(b.pendingSpawns, PendingSpawn{
					SpawnTime: now.Add(2 * time.Second), Name: fmt.Sprintf("%s 分裂體", enemy.Name), Energy: 999, Shell: &newShell, IsSlime: true, ReplicationCount: enemy.ReplicationCount + 1, MaxReplications: enemy.MaxReplications,
				})
				enemy.ReplicationCount = enemy.MaxReplications // 避免重複觸發
			}
		}

		// 處理重生
		remainingSpawns := []PendingSpawn{}
		for _, spawn := range b.pendingSpawns {
			if now.After(spawn.SpawnTime) {
				newEnemy := NewPlayer(spawn.Name, spawn.Energy, spawn.IsSlime, spawn.ReplicationCount, nil)
				newEnemy.CurrentShell = spawn.Shell
				b.enemies = append(b.enemies, newEnemy)
				allLogsThisTick = append(allLogsThisTick, fmt.Sprintf("[green]%s 重生了！[-:-:-]", spawn.Name))
				anyActionTaken = true
			} else {
				remainingSpawns = append(remainingSpawns, spawn)
			}
		}
		b.pendingSpawns = remainingSpawns

		// 勝利條件判斷
		if allEnemiesDefeated && !b.gameIsOver {
			allLogsThisTick = append(allLogsThisTick, "", "[::b][green]勝利！你擊敗了所有敵人的軀殼！ 按(Ctrl+C)離開。")
			b.gameIsOver = true
			anyActionTaken = true
		}

		if len(allLogsThisTick) > 0 {
			b.logHistory = append(b.logHistory, allLogsThisTick...)
			maxLogLines := 30
			if len(b.logHistory) > maxLogLines {
				b.logHistory = b.logHistory[len(b.logHistory)-maxLogLines:]
			}
		}

		b.app.QueueUpdateDraw(func() {
			if anyActionTaken {
				b.updateAllViews()
				if len(allLogsThisTick) > 0 {
					b.battleLog.SetText(strings.Join(b.logHistory, "\n"))
					b.battleLog.ScrollToEnd()
				}
			} else {
				b.playerStatus.SetText(b.player.GetPlayerStatusText(b))
				if b.currentTargetIndex < len(b.enemies) {
					b.targetStatus.SetText(b.enemies[b.currentTargetIndex].GetEnemyStatusText())
				}
			}
		})
	}
}

// SetupUI 設置並返回戰鬥畫面的根元件
func (b *Battle) SetupUI() tview.Primitive {
	b.updateAllViews()

	rightPanel := tview.NewFlex().SetDirection(tview.FlexRow).AddItem(b.enemyList, 0, 1, true).AddItem(b.targetStatus, 10, 0, false)
	mainFlex := tview.NewFlex().AddItem(b.playerStatus, 0, 1, false).AddItem(rightPanel, 0, 1, true)
	b.mainLayout = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(mainFlex, 0, 2, true).
		AddItem(b.battleLog, 0, 1, false).
		AddItem(b.instructions, 1, 0, false)

	b.setupInputHandling()
	return b.mainLayout
}

func newEnemyListChanged(b *Battle) func(index int, mainText string, secondaryText string, shortcut rune) {
	return func(index int, mainText string, secondaryText string, shortcut rune) {
		b.currentTargetIndex = index
		b.updateAllViews()
	}
}

// setupInputHandling 設置按鍵輸入處理
func (b *Battle) setupInputHandling() {
	enemyListChanged := newEnemyListChanged(b)
	b.enemyList.SetChangedFunc(enemyListChanged)

	b.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlC:
			b.app.Stop()
			return nil
		case tcell.KeyEsc:
			if b.player.ActionState == "Channeling" || b.player.ActionState == "Casting" {
				b.interruptChanneling = true
			}
			return nil
		case tcell.KeyTab:
			nextIndex := (b.currentTargetIndex + 1) % len(b.enemies)
			for nextIndex != b.currentTargetIndex {
				if b.enemies[nextIndex].CurrentShell != nil {
					b.currentTargetIndex = nextIndex
					b.updateAllViews()
					return nil
				}
				nextIndex = (nextIndex + 1) % len(b.enemies)
			}
			return nil
		case tcell.KeyBacktab: // Shift+Tab
			nextIndex := (b.currentTargetIndex - 1 + len(b.enemies)) % len(b.enemies)
			for nextIndex != b.currentTargetIndex {
				if b.enemies[nextIndex].CurrentShell != nil {
					b.currentTargetIndex = nextIndex
					b.updateAllViews()
					return nil
				}
				nextIndex = (nextIndex - 1 + len(b.enemies)) % len(b.enemies)
			}
			return nil
		}

		if b.gameIsOver {
			return event
		}

		if b.player.ActionState == "Channeling" {
			b.interruptChanneling = true
		} else if b.player.ActionState != "Idle" || b.nextPlayerAction != nil {
			return event
		}

		runeKey := event.Rune()
		if runeKey >= '1' && runeKey <= '9' {
			index := int(runeKey - '1')
			if index < len(b.enemies) {
				b.currentTargetIndex = index
				b.updateAllViews()
			}
			return event
		}

		now := time.Now()
		target := b.enemies[b.currentTargetIndex]

		activeSkills := b.player.getActiveSkills()
		skillKeyMap := map[rune]int{'q': 0, 'w': 1, 'e': 2, 'r': 3, 'a': 4, 's': 5, 'd': 6, 'f': 7}

		if skillIndex, ok := skillKeyMap[runeKey]; ok {
			if len(activeSkills) > skillIndex {
				skill := activeSkills[skillIndex]

				if !(now.After(b.player.SkillCooldowns[skill.Name])) {
					return event
				}

				if skill.CanUse(b.player, target) {
					b.nextPlayerAction = skill
					b.player.CastingTarget = target // CanUse 已經處理了目標合法性
				}
			}
		}

		return event
	})
}

func main() {
	app := tview.NewApplication()
	battle := NewBattle(app)

	go battle.gameLoop()

	if err := app.SetRoot(battle.SetupUI(), true).SetFocus(battle.mainLayout).Run(); err != nil {
		panic(err)
	}
}
