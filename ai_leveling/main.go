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

// 預約的行動
var nextPlayerAction *AttackMove
var nextPlayerMeditate bool
var nextPlayerPossess bool

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

	// 對於單體攻擊，檢查主目標是否有效
	if !move.IsAoE && (mainTarget.CurrentShell == nil || mainTarget.CurrentShell.IsDefeated()) {
		return []string{fmt.Sprintf("目標 %s 的軀殼已無靈魂，無法攻擊。", mainTarget.Name)}
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
		// 單體攻擊
		mainTarget.CurrentShell.LoseHealth(finalDamage)
		logs = append(logs, fmt.Sprintf("   對 %s 的軀殼造成了 %d 點傷害！ (%d 基礎 + %d 力量)", mainTarget.Name, finalDamage, move.Damage, p.CurrentShell.Strength))
	}

	return logs
}

// createProgressBar 創建一個文字進度條
func createProgressBar(startTime, endTime time.Time, width int, color string) string {
	totalDuration := endTime.Sub(startTime)
	if totalDuration <= 0 {
		return strings.Repeat(" ", width+2) // 返回空格以維持排版
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
func (p *Player) GetPlayerStatusText() string {
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
		if nextPlayerAction != nil {
			queuedAction = fmt.Sprintf("預約: %s", nextPlayerAction.Name)
		} else if nextPlayerMeditate {
			queuedAction = "預約: 冥想"
		} else if nextPlayerPossess {
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

func main() {
	// --- 遊戲設定 ---
	slash := &AttackMove{Name: "揮砍", EnergyCost: 10, Damage: 15, CastTime: 0, ActionCooldown: 1 * time.Second, IsAoE: false}
	heavyStrike := &AttackMove{Name: "強力一擊", EnergyCost: 35, Damage: 80, CastTime: 1 * time.Second, ActionCooldown: 2 * time.Second, IsAoE: false}
	stomp := &AttackMove{Name: "踐踏", EnergyCost: 20, Damage: 4, CastTime: 2 * time.Second, ActionCooldown: 4 * time.Second, IsAoE: true, IsChanneling: true}
	bite := &AttackMove{Name: "啃咬", EnergyCost: 10, Damage: 12, CastTime: 0, ActionCooldown: 1 * time.Second, IsAoE: false}

	possessionCastTime := 2 * time.Second
	directPossessionCost := 60

	player := NewPlayer("英雄", 100, false, 0)
	player.CurrentShell = NewShell("人類軀殼", 500, 5, []*AttackMove{slash, heavyStrike})

	enemies := []*Player{
		NewPlayer("哥布林", 999, false, 0),
		NewPlayer("史萊姆", 999, true, 0),
		NewPlayer("骷髏兵", 999, false, 0),
	}
	enemies[0].CurrentShell = NewShell("哥布林軀殼", 80, 2, []*AttackMove{stomp})
	enemies[1].CurrentShell = NewShell("凝膠軀殼", 60, 5, []*AttackMove{bite})
	enemies[2].CurrentShell = NewShell("骸骨軀殼", 120, 8, []*AttackMove{stomp})

	var currentTargetIndex int = 0

	// --- TUI 介面設定 ---
	app := tview.NewApplication()
	var logHistory []string
	const maxLogLines = 100

	playerStatus := tview.NewTextView()
	playerStatus.SetDynamicColors(true).SetTextAlign(tview.AlignCenter).SetBorder(true).SetTitle("你的狀態")
	targetStatus := tview.NewTextView()
	targetStatus.SetDynamicColors(true).SetTextAlign(tview.AlignCenter).SetBorder(true).SetTitle("鎖定目標")
	enemyList := tview.NewList()
	enemyList.ShowSecondaryText(false).SetBorder(true).SetTitle("敵人清單")
	battleLog := tview.NewTextView()
	battleLog.SetDynamicColors(true).SetScrollable(true).SetBorder(true).SetTitle("戰鬥日誌")
	instructions := tview.NewTextView()
	instructions.SetDynamicColors(true)

	var enemyListChanged func(int, string, string, rune)

	updateAllViews := func() {
		playerStatus.SetText(player.GetPlayerStatusText())
		if currentTargetIndex < len(enemies) {
			targetStatus.SetText(enemies[currentTargetIndex].GetEnemyStatusText())
		}

		enemyList.SetChangedFunc(nil)
		enemyList.Clear()
		for i, enemy := range enemies {
			var status string
			if enemy.CurrentShell == nil {
				status = "[gray]靈體"
			} else if enemy.CurrentShell.IsDefeated() {
				status = "[purple]無主軀殼"
			} else {
				status = createValueBar(enemy.CurrentShell.Health, enemy.CurrentShell.MaxHealth, 10, "red")
			}

			prefix := "  "
			if i == currentTargetIndex {
				prefix = "[red]>>[-:-:-]"
			}
			mainText := fmt.Sprintf("%s %-8s %s", prefix, enemy.Name, status)

			enemyList.AddItem(mainText, "", 0, nil)
		}
		if currentTargetIndex < len(enemies) {
			enemyList.SetCurrentItem(currentTargetIndex)
		}
		enemyList.SetChangedFunc(enemyListChanged)

		baseInstructions := ""
		if currentTargetIndex < len(enemies) {
			target := enemies[currentTargetIndex]
			if player.CurrentShell != nil {
				skillText := []string{}
				skillKeys := []rune{'q', 'w'}
				now := time.Now()
				for i, skill := range player.CurrentShell.Skills {
					if i < len(skillKeys) {
						cd, onCD := player.SkillCooldowns[skill.Name]
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
					baseInstructions += fmt.Sprintf(" | [green](r) 附身 (耗%d)[-:-:-]", directPossessionCost)
				}
			} else {
				if target.CurrentShell != nil && target.CurrentShell.IsDefeated() {
					baseInstructions = fmt.Sprintf("[green](r) 附身 (耗%d)[-:-:-]", directPossessionCost)
				} else {
					baseInstructions = "靈體狀態：尋找無主的軀殼"
				}
			}
		}
		instructions.SetText(baseInstructions + " | (1-9)選敵 | (Esc)離開")
	}

	logHistory = append(logHistory, "戰鬥開始！")
	updateAllViews()

	rightPanel := tview.NewFlex().SetDirection(tview.FlexRow).AddItem(enemyList, 0, 1, true).AddItem(targetStatus, 10, 0, false)
	mainFlex := tview.NewFlex().AddItem(playerStatus, 0, 1, false).AddItem(rightPanel, 0, 1, true)
	mainLayout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(mainFlex, 0, 2, true).
		AddItem(battleLog, 0, 1, false).
		AddItem(instructions, 1, 0, false)

	// --- 遊戲邏輯與主迴圈 ---
	var gameIsOver bool = false
	var pendingSpawns []PendingSpawn

	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		for range ticker.C {
			if gameIsOver {
				continue
			}

			var logsThisTick []string
			var actionTaken bool = false
			now := time.Now()

			// --- 玩家行動邏輯 ---
			if currentTargetIndex < len(enemies) {
				target := enemies[currentTargetIndex]

				switch player.ActionState {
				case "Idle":
					if nextPlayerMeditate {
						player.ActionState = "Channeling"
						player.CastingMove = &AttackMove{Name: "冥想"}
						player.LastChannelTick = now
						actionTaken = true
					} else if nextPlayerAction != nil {
						state := "Casting"
						if nextPlayerAction.IsChanneling {
							state = "Channeling"
							player.LastChannelTick = now
						}
						player.ActionState = state
						player.ActionStartTime = now
						player.CastingMove = nextPlayerAction
						player.EffectTime = now.Add(nextPlayerAction.CastTime)
						player.EffectApplied = false
						actionTaken = true
					} else if nextPlayerPossess {
						player.ActionState = "Casting"
						player.ActionStartTime = now
						player.CastingMove = &AttackMove{Name: "附身"}
						player.EffectTime = now.Add(possessionCastTime)
						player.EffectApplied = false
						actionTaken = true
					}
					nextPlayerAction, nextPlayerMeditate, nextPlayerPossess = nil, false, false

				case "Casting":
					if !player.EffectApplied && now.After(player.EffectTime) {
						if player.CastingMove.Name == "附身" {
							if target.CurrentShell != nil && target.CurrentShell.IsDefeated() {
								if player.CurrentShell != nil {
									logsThisTick = append(logsThisTick, fmt.Sprintf("[purple]你拋棄了 %s，附身到 %s 的軀殼上！[-:-:-]", player.CurrentShell.Name, target.Name))
								} else {
									logsThisTick = append(logsThisTick, fmt.Sprintf("[green]你以靈體狀態，成功附身到 %s 的軀殼上！[-:-:-]", target.Name))
								}
								player.LoseEnergy(directPossessionCost)
								target.CurrentShell.Health = target.CurrentShell.MaxHealth
								player.CurrentShell = target.CurrentShell
								target.CurrentShell = nil
							}
						} else {
							logsThisTick = append(logsThisTick, player.Attack(target, enemies, player.CastingMove)...)
						}
						player.SkillCooldowns[player.CastingMove.Name] = now.Add(player.CastingMove.ActionCooldown)
						player.ActionState = "Idle"
						actionTaken = true
					}

				case "Channeling":
					if nextPlayerAction != nil || nextPlayerPossess {
						player.ActionState = "Idle"
						actionTaken = true
					} else {
						if now.After(player.EffectTime) { // 引導結束
							player.SkillCooldowns[player.CastingMove.Name] = now.Add(player.CastingMove.ActionCooldown)
							player.ActionState = "Idle"
							actionTaken = true
						} else if now.Sub(player.LastChannelTick) >= 500*time.Millisecond {
							if player.CastingMove.Name == "冥想" {
								player.GainEnergy(5)
							} else if player.CastingMove.Name == "踐踏" {
								logsThisTick = append(logsThisTick, "[yellow]踐踏造成了範圍傷害！[-:-:-]")
								finalDamage := player.CastingMove.Damage + player.CurrentShell.Strength
								for _, t := range enemies {
									if t.CurrentShell != nil && !t.CurrentShell.IsDefeated() {
										t.CurrentShell.LoseHealth(finalDamage)
										logsThisTick = append(logsThisTick, fmt.Sprintf("   對 %s 的軀殼造成了 %d 點傷害！", t.Name, finalDamage))
									}
								}
							}
							player.LastChannelTick = now
							actionTaken = true
						}
					}
				}
			}

			// 被動效果
			if player.CurrentShell == nil && player.ActionState == "Idle" {
				player.GainEnergy(1)
			}
			if player.CurrentShell != nil && player.CurrentShell.IsDefeated() {
				logsThisTick = append(logsThisTick, "[orange]你的軀殼被摧毀了！你現在是靈體狀態。[-:-:-]")
				player.CurrentShell = nil
				player.ActionState = "Idle"
				actionTaken = true
			}

			// --- 敵人 AI 行動邏輯 ---
			allEnemiesDefeated := true
			for _, enemy := range enemies {
				if enemy.CurrentShell != nil {
					if enemy.CurrentShell.IsDefeated() {
						if enemy.IsSlime && enemy.ReplicationCount < enemy.MaxReplications {
							newShell := *enemy.CurrentShell
							newShell.Health = newShell.MaxHealth / 2
							newShell.MaxHealth = newShell.MaxHealth / 2

							pendingSpawns = append(pendingSpawns, PendingSpawn{
								SpawnTime:        now.Add(2 * time.Second),
								Name:             fmt.Sprintf("%s 分裂體", enemy.Name),
								Energy:           999,
								Shell:            &newShell,
								IsSlime:          true,
								ReplicationCount: enemy.ReplicationCount + 1,
								MaxReplications:  enemy.MaxReplications,
							})
							enemy.ReplicationCount = enemy.MaxReplications
						}
					} else {
						allEnemiesDefeated = false
						if enemy.ActionState == "Idle" && player.CurrentShell != nil && len(enemy.CurrentShell.Skills) > 0 {
							skill := enemy.CurrentShell.Skills[0]
							if now.After(enemy.SkillCooldowns[skill.Name]) {
								enemy.ActionState = "Casting"
								enemy.ActionStartTime = now
								enemy.CastingMove = skill
								enemy.EffectTime = now.Add(enemy.CastingMove.CastTime)
								enemy.EffectApplied = false
								actionTaken = true
							}
						} else if enemy.ActionState == "Casting" {
							if !enemy.EffectApplied && now.After(enemy.EffectTime) {
								logsThisTick = append(logsThisTick, "")
								logsThisTick = append(logsThisTick, enemy.Attack(player, []*Player{player}, enemy.CastingMove)...)
								enemy.EffectApplied = true
								enemy.SkillCooldowns[enemy.CastingMove.Name] = now.Add(enemy.CastingMove.ActionCooldown)
								enemy.ActionState = "Idle"
								actionTaken = true
							}
						}
					}
				}
			}

			// 處理重生
			remainingSpawns := []PendingSpawn{}
			for _, spawn := range pendingSpawns {
				if now.After(spawn.SpawnTime) {
					newEnemy := NewPlayer(spawn.Name, spawn.Energy, spawn.IsSlime, spawn.ReplicationCount)
					newEnemy.CurrentShell = spawn.Shell
					enemies = append(enemies, newEnemy)
					logsThisTick = append(logsThisTick, fmt.Sprintf("[green]%s 重生了！[-:-:-]", spawn.Name))
					actionTaken = true
				} else {
					remainingSpawns = append(remainingSpawns, spawn)
				}
			}
			pendingSpawns = remainingSpawns

			if allEnemiesDefeated && !gameIsOver {
				logsThisTick = append(logsThisTick, "", "[::b][green]勝利！你擊敗了所有敵人的軀殼！ 按(q)離開。")
				gameIsOver = true
				actionTaken = true
			}

			if len(logsThisTick) > 0 {
				logHistory = append(logHistory, logsThisTick...)
				if len(logHistory) > maxLogLines {
					logHistory = logHistory[len(logHistory)-maxLogLines:]
				}
			}

			app.QueueUpdateDraw(func() {
				if actionTaken {
					updateAllViews()
					if len(logsThisTick) > 0 {
						battleLog.SetText(strings.Join(logHistory, "\n"))
						battleLog.ScrollToEnd()
					}
				} else {
					playerStatus.SetText(player.GetPlayerStatusText())
					if currentTargetIndex < len(enemies) {
						targetStatus.SetText(enemies[currentTargetIndex].GetEnemyStatusText())
					}
				}
			})
		}
	}()

	// --- 輸入處理 ---
	enemyListChanged = func(index int, mainText string, secondaryText string, shortcut rune) {
		currentTargetIndex = index
		updateAllViews()
	}
	enemyList.SetChangedFunc(enemyListChanged)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			app.Stop()
			return nil
		case tcell.KeyTab:
			currentTargetIndex = (currentTargetIndex + 1) % len(enemies)
			updateAllViews()
			return nil
		case tcell.KeyBacktab:
			currentTargetIndex = (currentTargetIndex - 1 + len(enemies)) % len(enemies)
			updateAllViews()
			return nil
		}

		if gameIsOver {
			return event
		}

		if player.ActionState == "Channeling" {
			player.ActionState = "Idle"
		}

		runeKey := event.Rune()
		if runeKey >= '1' && runeKey <= '9' {
			index := int(runeKey - '1')
			if index < len(enemies) {
				currentTargetIndex = index
				updateAllViews()
			}
			return event
		}

		now := time.Now()
		switch runeKey {
		case 'q':
			if player.CurrentShell != nil && len(player.CurrentShell.Skills) > 0 {
				skill := player.CurrentShell.Skills[0]
				if now.After(player.SkillCooldowns[skill.Name]) && player.Energy >= skill.EnergyCost {
					nextPlayerAction, nextPlayerMeditate, nextPlayerPossess = skill, false, false
				}
			}
		case 'w':
			if player.CurrentShell != nil && len(player.CurrentShell.Skills) > 1 {
				skill := player.CurrentShell.Skills[1]
				if now.After(player.SkillCooldowns[skill.Name]) && player.Energy >= skill.EnergyCost {
					nextPlayerAction, nextPlayerMeditate, nextPlayerPossess = skill, false, false
				}
			}
		case 'e':
			if player.CurrentShell != nil {
				nextPlayerAction, nextPlayerMeditate, nextPlayerPossess = nil, true, false
			}
		case 'r':
			if currentTargetIndex < len(enemies) {
				target := enemies[currentTargetIndex]
				if target.CurrentShell != nil && target.CurrentShell.IsDefeated() && player.Energy >= directPossessionCost {
					nextPlayerAction, nextPlayerMeditate, nextPlayerPossess = nil, false, true
				}
			}
		}

		return event
	})

	if err := app.SetRoot(mainLayout, true).SetFocus(mainLayout).Run(); err != nil {
		panic(err)
	}
}
