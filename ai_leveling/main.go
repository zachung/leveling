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
	CastTime       time.Duration // 施法時間
	ActionCooldown time.Duration // 動作冷卻
	IsAoE          bool          // 是否為範圍攻擊
	IsChanneling   bool          // 是否為引導技能
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
	CurrentShell *Shell // 當前附身的軀殼，可能為 nil

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

	currentTargetIndex int
	logHistory         []string
	pendingSpawns      []PendingSpawn
	gameIsOver         bool

	// 預約的行動
	nextPlayerAction   *AttackMove
	nextPlayerMeditate bool
	nextPlayerPossess  bool

	// UI 元件
	playerStatus *tview.TextView
	targetStatus *tview.TextView
	enemyList    *tview.List
	battleLog    *tview.TextView
	instructions *tview.TextView
	mainLayout   *tview.Flex
}

// NewPlayer 創建一個新的靈魂實例
func NewPlayer(name string, energy int, isSlime bool, replications int) *Player {
	return &Player{
		Name:             name,
		Energy:           energy,
		MaxEnergy:        energy,
		ActionState:      "Idle",
		SkillCooldowns:   make(map[string]time.Time),
		IsSlime:          isSlime,
		ReplicationCount: replications,
		MaxReplications:  2,
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

	logs := []string{fmt.Sprintf("➡️ %s 的 [%s] 擊中了目標！", p.Name, move.Name)}

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
	}

	switch p.ActionState {
	case "Casting", "Channeling":
		status.WriteString(fmt.Sprintf("[yellow]%s: %s (%.1fs)[-:-:-]\n", p.ActionState, p.CastingMove.Name, time.Until(p.EffectTime).Seconds()))
		status.WriteString(createProgressBar(p.ActionStartTime, p.EffectTime, 20, "yellow"))
	case "Idle":
		if p.CurrentShell != nil {
			status.WriteString("[green]狀態: 可行動[-:-:-]")
		} else {
			status.WriteString("[purple]狀態: 靈體[-:-:-]")
		}
	}

	if p.ActionState == "Idle" {
		var queuedAction string
		if b.nextPlayerAction != nil {
			queuedAction = fmt.Sprintf("預約: %s", b.nextPlayerAction.Name)
		} else if b.nextPlayerMeditate {
			queuedAction = "預約: 冥想"
		} else if b.nextPlayerPossess {
			queuedAction = "預約: 附身"
		}
		if queuedAction != "" {
			status.WriteString(fmt.Sprintf("\n[cyan]%s[-:-:-]", queuedAction))
		}
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
	slash := &AttackMove{Name: "揮砍", EnergyCost: 10, Damage: 15, CastTime: 0, ActionCooldown: 1 * time.Second, IsAoE: false}
	heavyStrike := &AttackMove{Name: "強力一擊", EnergyCost: 35, Damage: 80, CastTime: 1 * time.Second, ActionCooldown: 2 * time.Second, IsAoE: false}
	stomp := &AttackMove{Name: "踐踏", EnergyCost: 20, Damage: 4, CastTime: 2 * time.Second, ActionCooldown: 4 * time.Second, IsAoE: true, IsChanneling: true}
	bite := &AttackMove{Name: "啃咬", EnergyCost: 10, Damage: 12, CastTime: 0, ActionCooldown: 1 * time.Second, IsAoE: false}
	soulEjection := &AttackMove{Name: "靈魂出竅", EnergyCost: 10, CastTime: 500 * time.Millisecond}

	player := NewPlayer("英雄", 100, false, 0)
	player.CurrentShell = NewShell("人類軀殼", 500, 5, []*AttackMove{slash, heavyStrike, soulEjection})

	enemies := []*Player{
		NewPlayer("哥布林", 999, false, 0),
		NewPlayer("史萊姆", 999, true, 0),
		NewPlayer("骷髏兵", 999, false, 0),
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
	b.playerStatus.SetDynamicColors(true).SetTextAlign(tview.AlignCenter).SetBorder(true).SetTitle("你的狀態")
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
		mainText := fmt.Sprintf("%s %-8s %s", prefix, enemy.Name, status)

		b.enemyList.AddItem(mainText, "", 0, nil)
	}
	if b.currentTargetIndex < len(b.enemies) {
		b.enemyList.SetCurrentItem(b.currentTargetIndex)
	}
	b.enemyList.SetChangedFunc(enemyListChanged)

	baseInstructions := ""
	if b.currentTargetIndex < len(b.enemies) {
		target := b.enemies[b.currentTargetIndex]
		if b.player.CurrentShell != nil {
			skillText := []string{}
			skillKeys := []rune{'q', 'w'}
			now := time.Now()
			for i, skill := range b.player.CurrentShell.Skills {
				if i < len(skillKeys) {
					cd, onCD := b.player.SkillCooldowns[skill.Name]
					if onCD && now.Before(cd) {
						skillText = append(skillText, fmt.Sprintf("[gray](%c) %s (%.1fs)", skillKeys[i], skill.Name, time.Until(cd).Seconds()))
					} else {
						skillText = append(skillText, fmt.Sprintf("[yellow](%c) %s", skillKeys[i], skill.Name))
					}
				}
			}
			skillText = append(skillText, "[yellow](e) 冥想")
			baseInstructions = strings.Join(skillText, " | ")

			if target.CurrentShell != nil && target.CurrentShell.IsDefeated() {
				baseInstructions += fmt.Sprintf(" | [green](r) 附身[-:-:-]")
			} else {
				baseInstructions += " | [purple](r) 靈魂出竅[-:-:-]"
			}
		} else {
			if target.CurrentShell != nil && target.CurrentShell.IsDefeated() {
				baseInstructions = fmt.Sprintf("[green](r) 附身[-:-:-]")
			} else {
				baseInstructions = "靈體狀態：尋找無主的軀殼"
			}
		}
	}
	b.instructions.SetText(baseInstructions + " | (1-9)選敵 | (Esc)離開")
}

// gameLoop 是戰鬥的主迴圈
func (b *Battle) gameLoop() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		if b.gameIsOver {
			continue
		}

		var logsThisTick []string
		var actionTaken bool = false
		now := time.Now()

		// --- 玩家行動邏輯 ---
		switch b.player.ActionState {
		case "Idle":
			if b.nextPlayerMeditate {
				b.player.ActionState = "Channeling"
				b.player.CastingMove = &AttackMove{Name: "冥想"}
				b.player.LastChannelTick = now
				actionTaken = true
			} else if b.nextPlayerAction != nil {
				state := "Casting"
				if b.nextPlayerAction.IsChanneling {
					state = "Channeling"
					b.player.LastChannelTick = now
				}
				b.player.ActionState = state
				b.player.ActionStartTime = now
				b.player.CastingMove = b.nextPlayerAction
				b.player.EffectTime = now.Add(b.nextPlayerAction.CastTime)
				b.player.EffectApplied = false
				actionTaken = true
			} else if b.nextPlayerPossess {
				b.player.ActionState = "Casting"
				b.player.ActionStartTime = now
				b.player.CastingMove = &AttackMove{Name: "附身"}
				b.player.EffectTime = now.Add(2 * time.Second)
				b.player.EffectApplied = false
				actionTaken = true
			}
			b.nextPlayerAction, b.nextPlayerMeditate, b.nextPlayerPossess = nil, false, false

		case "Casting":
			if !b.player.EffectApplied && now.After(b.player.EffectTime) {
				target := b.player.CastingTarget
				if b.player.CastingMove.Name == "附身" {
					if target.CurrentShell != nil && target.CurrentShell.IsDefeated() {
						if b.player.CurrentShell != nil {
							logsThisTick = append(logsThisTick, fmt.Sprintf("[purple]你拋棄了 %s，附身到 %s 的軀殼上！[-:-:-]", b.player.CurrentShell.Name, target.Name))
						} else {
							logsThisTick = append(logsThisTick, fmt.Sprintf("[green]你以靈體狀態，成功附身到 %s 的軀殼上！[-:-:-]", target.Name))
						}
						b.player.LoseEnergy(60)
						target.CurrentShell.Health = target.CurrentShell.MaxHealth
						b.player.CurrentShell = target.CurrentShell
						target.CurrentShell = nil
					}
				} else if b.player.CastingMove.Name == "靈魂出竅" {
					logsThisTick = append(logsThisTick, "[purple]你施展了靈魂出竅，脫離了當前的軀殼！[-:-:-]")
					b.player.CurrentShell = nil
				} else {
					logsThisTick = append(logsThisTick, b.player.Attack(target, b.enemies, b.player.CastingMove)...)
				}
				b.player.SkillCooldowns[b.player.CastingMove.Name] = now.Add(b.player.CastingMove.ActionCooldown)
				b.player.ActionState = "Idle"
				actionTaken = true
			}

		case "Channeling":
			if b.nextPlayerAction != nil || b.nextPlayerPossess {
				b.player.ActionState = "Idle"
				actionTaken = true
			} else {
				// 修正：區分持續引導和有時間限制的引導
				switch b.player.CastingMove.Name {
				case "冥想":
					if now.Sub(b.player.LastChannelTick) >= 500*time.Millisecond {
						b.player.GainEnergy(5)
						b.player.LastChannelTick = now
						actionTaken = true
					}
				case "踐踏":
					if now.After(b.player.EffectTime) {
						b.player.SkillCooldowns[b.player.CastingMove.Name] = now.Add(b.player.CastingMove.ActionCooldown)
						b.player.ActionState = "Idle"
						actionTaken = true
					} else if now.Sub(b.player.LastChannelTick) >= 500*time.Millisecond {
						logsThisTick = append(logsThisTick, "[yellow]踐踏造成了範圍傷害！[-:-:-]")
						finalDamage := b.player.CastingMove.Damage + b.player.CurrentShell.Strength
						for _, t := range b.enemies {
							if t.CurrentShell != nil && !t.CurrentShell.IsDefeated() {
								t.CurrentShell.LoseHealth(finalDamage)
								logsThisTick = append(logsThisTick, fmt.Sprintf("   對 %s 的軀殼造成了 %d 點傷害！", t.Name, finalDamage))
							}
						}
						b.player.LastChannelTick = now
						actionTaken = true
					}
				}
			}
		}

		if b.player.CurrentShell == nil && b.player.ActionState == "Idle" {
			b.player.GainEnergy(1)
		}
		if b.player.CurrentShell != nil && b.player.CurrentShell.IsDefeated() {
			logsThisTick = append(logsThisTick, "[orange]你的軀殼被摧毀了！你現在是靈體狀態。[-:-:-]")
			b.player.CurrentShell = nil
			b.player.ActionState = "Idle"
			actionTaken = true
		}

		allEnemiesDefeated := true
		for _, enemy := range b.enemies {
			if enemy.CurrentShell != nil {
				if enemy.CurrentShell.IsDefeated() {
					if enemy.IsSlime && enemy.ReplicationCount < enemy.MaxReplications {
						newShell := *enemy.CurrentShell
						newShell.Health = newShell.MaxHealth / 2
						newShell.MaxHealth = newShell.MaxHealth / 2
						b.pendingSpawns = append(b.pendingSpawns, PendingSpawn{SpawnTime: now.Add(2 * time.Second), Name: fmt.Sprintf("%s 分裂體", enemy.Name), Energy: 999, Shell: &newShell, IsSlime: true, ReplicationCount: enemy.ReplicationCount + 1, MaxReplications: enemy.MaxReplications})
						enemy.ReplicationCount = enemy.MaxReplications
					}
				} else {
					allEnemiesDefeated = false
					if enemy.ActionState == "Idle" && b.player.CurrentShell != nil && len(enemy.CurrentShell.Skills) > 0 {
						skill := enemy.CurrentShell.Skills[0]
						if now.After(enemy.SkillCooldowns[skill.Name]) {
							enemy.ActionState = "Casting"
							enemy.ActionStartTime = now
							enemy.CastingMove = skill
							enemy.EffectTime = now.Add(enemy.CastingMove.CastTime)
							enemy.EffectApplied = false
							enemy.CastingTarget = b.player
							actionTaken = true
						}
					} else if enemy.ActionState == "Casting" {
						if !enemy.EffectApplied && now.After(enemy.EffectTime) {
							logsThisTick = append(logsThisTick, "")
							logsThisTick = append(logsThisTick, enemy.Attack(enemy.CastingTarget, []*Player{b.player}, enemy.CastingMove)...)
							enemy.EffectApplied = true
							enemy.SkillCooldowns[enemy.CastingMove.Name] = now.Add(enemy.CastingMove.ActionCooldown)
							enemy.ActionState = "Idle"
							actionTaken = true
						}
					}
				}
			}
		}

		remainingSpawns := []PendingSpawn{}
		for _, spawn := range b.pendingSpawns {
			if now.After(spawn.SpawnTime) {
				newEnemy := NewPlayer(spawn.Name, spawn.Energy, spawn.IsSlime, spawn.ReplicationCount)
				newEnemy.CurrentShell = spawn.Shell
				b.enemies = append(b.enemies, newEnemy)
				logsThisTick = append(logsThisTick, fmt.Sprintf("[green]%s 重生了！[-:-:-]", spawn.Name))
				actionTaken = true
			} else {
				remainingSpawns = append(remainingSpawns, spawn)
			}
		}
		b.pendingSpawns = remainingSpawns

		if allEnemiesDefeated && !b.gameIsOver {
			logsThisTick = append(logsThisTick, "", "[::b][green]勝利！你擊敗了所有敵人的軀殼！ 按(q)離開。")
			b.gameIsOver = true
			actionTaken = true
		}

		if len(logsThisTick) > 0 {
			b.logHistory = append(b.logHistory, logsThisTick...)
			maxLogLines := 30
			if len(b.logHistory) > maxLogLines {
				b.logHistory = b.logHistory[len(b.logHistory)-maxLogLines:]
			}
		}

		b.app.QueueUpdateDraw(func() {
			if actionTaken {
				b.updateAllViews()
				if len(logsThisTick) > 0 {
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
		case tcell.KeyEsc:
			b.app.Stop()
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

		if b.player.ActionState != "Idle" || b.nextPlayerAction != nil || b.nextPlayerMeditate || b.nextPlayerPossess {
			return event
		}

		if b.player.ActionState == "Channeling" {
			b.player.ActionState = "Idle"
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

		switch runeKey {
		case 'q':
			if b.player.CurrentShell != nil && len(b.player.CurrentShell.Skills) > 0 {
				skill := b.player.CurrentShell.Skills[0]
				if now.After(b.player.SkillCooldowns[skill.Name]) && b.player.Energy >= skill.EnergyCost && (skill.IsAoE || (target.CurrentShell != nil && !target.CurrentShell.IsDefeated())) {
					b.nextPlayerAction, b.nextPlayerMeditate, b.nextPlayerPossess = skill, false, false
					b.player.CastingTarget = target
				}
			}
		case 'w':
			if b.player.CurrentShell != nil && len(b.player.CurrentShell.Skills) > 1 {
				skill := b.player.CurrentShell.Skills[1]
				if now.After(b.player.SkillCooldowns[skill.Name]) && b.player.Energy >= skill.EnergyCost && (skill.IsAoE || (target.CurrentShell != nil && !target.CurrentShell.IsDefeated())) {
					b.nextPlayerAction, b.nextPlayerMeditate, b.nextPlayerPossess = skill, false, false
					b.player.CastingTarget = target
				}
			}
		case 'e':
			if b.player.CurrentShell != nil {
				b.nextPlayerAction, b.nextPlayerMeditate, b.nextPlayerPossess = nil, true, false
			}
		case 'r':
			if b.player.CurrentShell != nil { // 靈魂出竅
				skill := &AttackMove{Name: "靈魂出竅", EnergyCost: 10, CastTime: 500 * time.Millisecond}
				if b.player.Energy >= skill.EnergyCost {
					b.nextPlayerAction, b.nextPlayerMeditate, b.nextPlayerPossess = skill, false, false
					b.player.CastingTarget = nil // 無目標
				}
			} else { // 附身
				if target.CurrentShell != nil && target.CurrentShell.IsDefeated() && b.player.Energy >= 60 {
					b.nextPlayerAction, b.nextPlayerMeditate, b.nextPlayerPossess = nil, false, true
					b.player.CastingTarget = target
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
