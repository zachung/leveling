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
}

// Shell 代表一個可被附身的軀殼
type Shell struct {
	Name      string
	Health    int
	MaxHealth int
	Strength  int // 力量屬性，影響傷害
	AI_Attack *AttackMove
}

// Player 代表玩家或敵人的核心靈魂
type Player struct {
	Name         string
	Energy       int
	MaxEnergy    int
	CurrentShell *Shell // 當前附身的軀殼，可能為 nil

	// 動作狀態機 (玩家與敵人都使用)
	ActionState      string
	ActionStartTime  time.Time
	StateFinishTime  time.Time
	EffectTime       time.Time
	EffectApplied    bool
	CastingMove      *AttackMove
	LastMeditateTick time.Time
}

// 預約的行動
var nextPlayerAction *AttackMove
var nextPlayerMeditate bool
var nextPlayerPossess bool

// NewPlayer 創建一個新的靈魂實例
func NewPlayer(name string, energy int) *Player {
	return &Player{
		Name:        name,
		Energy:      energy,
		MaxEnergy:   energy,
		ActionState: "Idle",
	}
}

// NewShell 創建一個新的軀殼實例
func NewShell(name string, health int, strength int, aiAttack *AttackMove) *Shell {
	return &Shell{
		Name:      name,
		Health:    health,
		MaxHealth: health,
		Strength:  strength,
		AI_Attack: aiAttack,
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
func (p *Player) Attack(target *Player, move *AttackMove) []string {
	if p.CurrentShell == nil {
		return []string{"靈體狀態無法攻擊！"}
	}
	if target.CurrentShell == nil || target.CurrentShell.IsDefeated() {
		return []string{fmt.Sprintf("目標 %s 的軀殼已無靈魂，無法攻擊。", target.Name)}
	}

	logs := []string{fmt.Sprintf("➡️ %s 的 [%s] 擊中了 %s！", p.Name, move.Name, target.Name)}
	p.LoseEnergy(move.EnergyCost)
	logs = append(logs, fmt.Sprintf("   %s 消耗了 %d 點能量。", p.Name, move.EnergyCost))

	finalDamage := move.Damage + p.CurrentShell.Strength
	target.CurrentShell.LoseHealth(finalDamage)
	logs = append(logs, fmt.Sprintf("   對 %s 的軀殼造成了 %d 點傷害！ (%d 基礎 + %d 力量)", target.Name, finalDamage, move.Damage, p.CurrentShell.Strength))

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

	now := time.Now()
	switch p.ActionState {
	case "Casting":
		if now.Before(p.EffectTime) {
			status.WriteString(fmt.Sprintf("[yellow]施法中: %s (%.1fs)[-:-:-]\n", p.CastingMove.Name, time.Until(p.EffectTime).Seconds()))
			status.WriteString(createProgressBar(p.ActionStartTime, p.EffectTime, 20, "yellow"))
		} else {
			status.WriteString(fmt.Sprintf("[red]冷卻中 (%.1fs)[-:-:-]\n", time.Until(p.StateFinishTime).Seconds()))
			status.WriteString(createProgressBar(p.EffectTime, p.StateFinishTime, 20, "red"))
		}
	case "Channeling":
		status.WriteString("[green]引導中: 冥想[-:-:-]")
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

		now := time.Now()
		switch p.ActionState {
		case "Casting":
			if now.Before(p.EffectTime) {
				status.WriteString(fmt.Sprintf("[yellow]施法中: %s (%.1fs)[-:-:-]\n", p.CastingMove.Name, time.Until(p.EffectTime).Seconds()))
				status.WriteString(createProgressBar(p.ActionStartTime, p.EffectTime, 20, "yellow"))
			} else {
				status.WriteString(fmt.Sprintf("[red]冷卻中 (%.1fs)[-:-:-]\n", time.Until(p.StateFinishTime).Seconds()))
				status.WriteString(createProgressBar(p.EffectTime, p.StateFinishTime, 20, "red"))
			}
		case "Idle":
			status.WriteString("[green]狀態: 可行動[-:-:-]")
		}
	}
	return status.String()
}

func main() {
	// --- 遊戲設定 ---
	slash := &AttackMove{Name: "揮砍", EnergyCost: 10, Damage: 15, CastTime: 500 * time.Millisecond, ActionCooldown: 500 * time.Millisecond}
	heavyStrike := &AttackMove{Name: "強力一擊", EnergyCost: 35, Damage: 80, CastTime: 1 * time.Second, ActionCooldown: 1 * time.Second}
	stomp := &AttackMove{Name: "踐踏", EnergyCost: 1, Damage: 8, CastTime: 500 * time.Millisecond, ActionCooldown: 2 * time.Second}
	bite := &AttackMove{Name: "啃咬", EnergyCost: 1, Damage: 12, CastTime: 500 * time.Millisecond, ActionCooldown: 2 * time.Second}

	possessionCastTime := 2 * time.Second
	possessionCooldown := 1 * time.Second
	directPossessionCost := 60

	player := NewPlayer("英雄", 100)
	player.CurrentShell = NewShell("人類軀殼", 100, 5, nil)

	enemies := []*Player{
		NewPlayer("哥布林", 999),
		NewPlayer("史萊姆", 999),
		NewPlayer("骷髏兵", 999),
	}
	enemies[0].CurrentShell = NewShell("哥布林軀殼", 80, 2, stomp)
	enemies[1].CurrentShell = NewShell("凝膠軀殼", 60, 5, bite)
	enemies[2].CurrentShell = NewShell("骸骨軀殼", 120, 8, stomp)

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
		targetStatus.SetText(enemies[currentTargetIndex].GetEnemyStatusText())

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
		enemyList.SetCurrentItem(currentTargetIndex)
		enemyList.SetChangedFunc(enemyListChanged)

		baseInstructions := ""
		target := enemies[currentTargetIndex]
		if player.CurrentShell != nil {
			baseInstructions = fmt.Sprintf("[yellow](q) %s | (w) %s | (e) %s", slash.Name, heavyStrike.Name, "冥想")
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
		instructions.SetText(baseInstructions + " | (1-3)選敵 | (Esc)離開")
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
			target := enemies[currentTargetIndex]

			switch player.ActionState {
			case "Idle":
				if nextPlayerMeditate {
					player.ActionState = "Channeling"
					player.LastMeditateTick = now
					actionTaken = true
				} else if nextPlayerAction != nil {
					player.ActionState = "Casting"
					player.ActionStartTime = now
					player.CastingMove = nextPlayerAction
					player.EffectTime = now.Add(nextPlayerAction.CastTime)
					player.StateFinishTime = now.Add(nextPlayerAction.CastTime + nextPlayerAction.ActionCooldown)
					player.EffectApplied = false
					actionTaken = true
				} else if nextPlayerPossess {
					player.ActionState = "Casting"
					player.ActionStartTime = now
					player.CastingMove = &AttackMove{Name: "附身"} // 僅用於顯示
					player.EffectTime = now.Add(possessionCastTime)
					player.StateFinishTime = now.Add(possessionCastTime + possessionCooldown)
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
						logsThisTick = append(logsThisTick, player.Attack(target, player.CastingMove)...)
					}
					player.EffectApplied = true
					actionTaken = true
				}
				if now.After(player.StateFinishTime) {
					player.ActionState = "Idle"
					actionTaken = true
				}

			case "Channeling": // 冥想
				if nextPlayerAction != nil || nextPlayerPossess {
					player.ActionState = "Idle" // 打斷引導
					actionTaken = true
				} else if now.Sub(player.LastMeditateTick) >= 500*time.Millisecond {
					player.GainEnergy(5)
					player.LastMeditateTick = now
					actionTaken = true
				}
			}

			// 被動效果
			if player.CurrentShell == nil && player.ActionState == "Idle" {
				player.GainEnergy(1) // 靈體狀態回能
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
						// 無主軀殼
					} else {
						allEnemiesDefeated = false
						switch enemy.ActionState {
						case "Idle":
							if player.CurrentShell != nil {
								enemy.ActionState = "Casting"
								enemy.ActionStartTime = now
								enemy.CastingMove = enemy.CurrentShell.AI_Attack
								enemy.EffectTime = now.Add(enemy.CastingMove.CastTime)
								enemy.StateFinishTime = now.Add(enemy.CastingMove.CastTime + enemy.CastingMove.ActionCooldown)
								enemy.EffectApplied = false
								actionTaken = true
							}
						case "Casting":
							if !enemy.EffectApplied && now.After(enemy.EffectTime) {
								logsThisTick = append(logsThisTick, "")
								logsThisTick = append(logsThisTick, enemy.Attack(player, enemy.CastingMove)...)
								enemy.EffectApplied = true
								actionTaken = true
							}
							if now.After(enemy.StateFinishTime) {
								enemy.ActionState = "Idle"
								actionTaken = true
							}
						}
					}
				}
			}

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
					targetStatus.SetText(enemies[currentTargetIndex].GetEnemyStatusText())
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
		// 優先處理通用按鍵
		switch event.Key() {
		case tcell.KeyEsc:
			app.Stop()
			return nil
		case tcell.KeyTab:
			currentTargetIndex = (currentTargetIndex + 1) % len(enemies)
			updateAllViews()
			return nil
		case tcell.KeyBacktab: // Shift+Tab
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

		// 處理 Rune (字元) 按鍵
		rune := event.Rune()
		if rune >= '1' && rune <= '9' {
			index := int(rune - '1')
			if index < len(enemies) {
				currentTargetIndex = index
				updateAllViews()
			}
			return event
		}

		switch rune {
		case 'q':
			if player.CurrentShell != nil && player.Energy >= slash.EnergyCost {
				nextPlayerAction, nextPlayerMeditate, nextPlayerPossess = slash, false, false
			}
		case 'w':
			if player.CurrentShell != nil && player.Energy >= heavyStrike.EnergyCost {
				nextPlayerAction, nextPlayerMeditate, nextPlayerPossess = heavyStrike, false, false
			}
		case 'e':
			if player.CurrentShell != nil {
				nextPlayerAction, nextPlayerMeditate, nextPlayerPossess = nil, true, false
			}
		case 'r':
			target := enemies[currentTargetIndex]
			if target.CurrentShell != nil && target.CurrentShell.IsDefeated() && player.Energy >= directPossessionCost {
				nextPlayerAction, nextPlayerMeditate, nextPlayerPossess = nil, false, true
			}
		}

		return event
	})

	if err := app.SetRoot(mainLayout, true).SetFocus(mainLayout).Run(); err != nil {
		panic(err)
	}
}
