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
	Name       string // 招式名稱
	EnergyCost int    // 消耗能量
	Damage     int    // 造成傷害 (削減生命)
}

// Shell 代表一個可被附身的軀殼
type Shell struct {
	Name      string
	Health    int
	MaxHealth int
	Strength  int // 力量屬性，影響傷害
	Cooldown  time.Time
	AI_Attack *AttackMove
}

// Player 代表玩家或敵人的核心靈魂
type Player struct {
	Name         string
	Energy       int
	MaxEnergy    int
	CurrentShell *Shell // 當前附身的軀殼，可能為 nil
}

// 預約的行動
var nextPlayerAction *AttackMove
var nextPlayerMeditate bool
var nextPlayerDirectPossess bool // 統一的附身行動

// NewPlayer 創建一個新的靈魂實例
func NewPlayer(name string, energy int) *Player {
	return &Player{
		Name:      name,
		Energy:    energy,
		MaxEnergy: energy,
	}
}

// NewShell 創建一個新的軀殼實例
func NewShell(name string, health int, strength int, aiAttack *AttackMove) *Shell {
	return &Shell{
		Name:      name,
		Health:    health,
		MaxHealth: health,
		Strength:  strength,
		Cooldown:  time.Now(),
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
	if p.Energy < move.EnergyCost {
		return []string{fmt.Sprintf("能量不足以使用 [%s]！", move.Name)}
	}

	logs := []string{fmt.Sprintf("➡️ %s 使用 [%s] 攻擊 %s！", p.Name, move.Name, target.Name)}
	p.LoseEnergy(move.EnergyCost)
	logs = append(logs, fmt.Sprintf("   %s 消耗了 %d 點能量。", p.Name, move.EnergyCost))

	finalDamage := move.Damage + p.CurrentShell.Strength
	target.CurrentShell.LoseHealth(finalDamage)
	logs = append(logs, fmt.Sprintf("   對 %s 的軀殼造成了 %d 點傷害！ (%d 基礎 + %d 力量)", target.Name, finalDamage, move.Damage, p.CurrentShell.Strength))

	return logs
}

// Meditate 讓玩家冥想以恢復能量
func (p *Player) Meditate() []string {
	restoreAmount := 20
	p.GainEnergy(restoreAmount)
	return []string{fmt.Sprintf("� %s 進行冥想，恢復了 %d 點能量。", p.Name, restoreAmount)}
}

// GetPlayerStatusText 獲取玩家狀態文字
func (p *Player) GetPlayerStatusText() string {
	var status strings.Builder
	status.WriteString(fmt.Sprintf("[::b]%s\n", p.Name))
	status.WriteString(fmt.Sprintf("%s\n", strings.Repeat("─", len(p.Name)+4)))
	status.WriteString(fmt.Sprintf("[blue]能量: %d / %d[-:-:-]\n", p.Energy, p.MaxEnergy))

	if p.CurrentShell != nil {
		status.WriteString(fmt.Sprintf("[red]生命: %d / %d[-:-:-]\n", p.CurrentShell.Health, p.CurrentShell.MaxHealth))
		status.WriteString(fmt.Sprintf("[orange]力量: %d[-:-:-]\n", p.CurrentShell.Strength))
		if time.Now().Before(p.CurrentShell.Cooldown) {
			status.WriteString(fmt.Sprintf("[yellow]狀態: 冷卻中 (%.1fs)[-:-:-]", time.Until(p.CurrentShell.Cooldown).Seconds()))
		} else {
			status.WriteString("[green]狀態: 可行動[-:-:-]")
		}
		if nextPlayerAction != nil {
			status.WriteString(fmt.Sprintf("\n[cyan]預約: %s[-:-:-]", nextPlayerAction.Name))
		} else if nextPlayerMeditate {
			status.WriteString("\n[cyan]預約: 冥想[-:-:-]")
		} else if nextPlayerDirectPossess {
			status.WriteString("\n[cyan]預約: 轉附身[-:-:-]")
		}
	} else {
		status.WriteString("[purple]狀態: 靈體[-:-:-]\n")
		if nextPlayerDirectPossess {
			status.WriteString("\n[cyan]預約: 附身[-:-:-]")
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
		status.WriteString(fmt.Sprintf("[orange]力量: %d[-:-:-]", p.CurrentShell.Strength))
	}
	return status.String()
}

func main() {
	// --- 遊戲設定 ---
	slash := &AttackMove{Name: "揮砍", EnergyCost: 10, Damage: 15}
	heavyStrike := &AttackMove{Name: "強力一擊", EnergyCost: 35, Damage: 80} // 提高傷害
	stomp := &AttackMove{Name: "踐踏", EnergyCost: 1, Damage: 8}
	bite := &AttackMove{Name: "啃咬", EnergyCost: 1, Damage: 12}
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
	enemyList.ShowSecondaryText(false).SetBorder(true).SetTitle("敵人清單 (用 ↑/↓ 選擇)")
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
				status = fmt.Sprintf("生命: %d/%d", enemy.CurrentShell.Health, enemy.CurrentShell.MaxHealth)
			}
			mainText := fmt.Sprintf("%s (%s)", enemy.Name, status)
			if i == currentTargetIndex {
				mainText = "[red]>> " + mainText + "[-:-:-]"
			}
			enemyList.AddItem(mainText, "", 0, nil)
		}
		enemyList.SetCurrentItem(currentTargetIndex)
		enemyList.SetChangedFunc(enemyListChanged)

		baseInstructions := ""
		target := enemies[currentTargetIndex]
		if player.CurrentShell != nil {
			baseInstructions = fmt.Sprintf("[yellow](1) %s | (2) %s | (m) %s", slash.Name, heavyStrike.Name, "冥想")
			if target.CurrentShell != nil && target.CurrentShell.IsDefeated() {
				baseInstructions += fmt.Sprintf(" | [green](x) 轉附身 (耗%d)[-:-:-]", directPossessionCost)
			}
		} else { // 靈體狀態
			if target.CurrentShell != nil && target.CurrentShell.IsDefeated() {
				baseInstructions = fmt.Sprintf("[green](x) 附身 (耗%d)[-:-:-]", directPossessionCost)
			} else {
				baseInstructions = "靈體狀態：尋找無主的軀殼"
			}
		}
		instructions.SetText(baseInstructions + " | (Tab)切換 | (q)uit")
	}

	logHistory = append(logHistory, "戰鬥開始！")
	updateAllViews()

	rightPanel := tview.NewFlex().SetDirection(tview.FlexRow).AddItem(enemyList, 0, 1, true).AddItem(targetStatus, 10, 0, false)
	mainFlex := tview.NewFlex().AddItem(playerStatus, 0, 1, false).AddItem(rightPanel, 0, 1, true)
	mainLayout := tview.NewFlex().SetDirection(tview.FlexRow).AddItem(mainFlex, 0, 1, true).AddItem(battleLog, 12, 0, false).AddItem(instructions, 1, 0, false)

	// --- 遊戲邏輯與主迴圈 ---
	cooldownDuration := 1 * time.Second
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

			// --- 玩家行動邏輯 ---
			target := enemies[currentTargetIndex]

			// 1. 處理附身 (可從任何狀態發起)
			if nextPlayerDirectPossess {
				if target.CurrentShell != nil && target.CurrentShell.IsDefeated() && player.Energy >= directPossessionCost {
					if player.CurrentShell != nil {
						logsThisTick = append(logsThisTick, fmt.Sprintf("[purple]你拋棄了 %s，附身到 %s 的軀殼上！[-:-:-]", player.CurrentShell.Name, target.Name))
					} else {
						logsThisTick = append(logsThisTick, fmt.Sprintf("[green]你以靈體狀態，成功附身到 %s 的軀殼上！[-:-:-]", target.Name))
					}
					player.LoseEnergy(directPossessionCost)

					target.CurrentShell.Health = target.CurrentShell.MaxHealth // 恢復軀殼
					player.CurrentShell = target.CurrentShell
					target.CurrentShell = nil // 敵方靈魂被永久驅逐

					player.CurrentShell.Cooldown = time.Now().Add(cooldownDuration)
					actionTaken = true
				}
			} else if player.CurrentShell != nil { // 2. 處理其他需要軀殼的行動
				if time.Now().After(player.CurrentShell.Cooldown) {
					if nextPlayerAction != nil {
						logsThisTick = append(logsThisTick, player.Attack(target, nextPlayerAction)...)
						player.CurrentShell.Cooldown = time.Now().Add(cooldownDuration)
						actionTaken = true
					} else if nextPlayerMeditate {
						logsThisTick = append(logsThisTick, player.Meditate()...)
						player.CurrentShell.Cooldown = time.Now().Add(cooldownDuration)
						actionTaken = true
					}
				}
			}

			// 重置所有預約
			nextPlayerAction, nextPlayerMeditate, nextPlayerDirectPossess = nil, false, false

			// 3. 處理被動效果
			if player.CurrentShell == nil { // 靈體狀態
				player.GainEnergy(1)
			}

			// --- 敵人 AI 行動邏輯 ---
			allEnemiesDefeated := true
			for _, enemy := range enemies {
				if enemy.CurrentShell != nil {
					if enemy.CurrentShell.IsDefeated() {
						// 軀殼已被擊敗，等待被附身或消失
					} else {
						allEnemiesDefeated = false
						if time.Now().After(enemy.CurrentShell.Cooldown) && player.CurrentShell != nil && !player.CurrentShell.IsDefeated() {
							logsThisTick = append(logsThisTick, "")
							logsThisTick = append(logsThisTick, enemy.Attack(player, enemy.CurrentShell.AI_Attack)...)
							enemy.CurrentShell.Cooldown = time.Now().Add(time.Duration(20+len(enemies)) * 100 * time.Millisecond)
							actionTaken = true
							if player.CurrentShell.IsDefeated() {
								logsThisTick = append(logsThisTick, "[orange]你的軀殼被摧毀了！你現在是靈體狀態。[-:-:-]")
								player.CurrentShell = nil
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
					battleLog.SetText(strings.Join(logHistory, "\n"))
					battleLog.ScrollToEnd()
				} else {
					playerStatus.SetText(player.GetPlayerStatusText())
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
		if event.Rune() == 'q' {
			app.Stop()
			return event
		}
		if gameIsOver {
			return event
		}

		if app.GetFocus() == enemyList {
			if event.Key() == tcell.KeyUp || event.Key() == tcell.KeyDown {
				return event
			}
		}

		switch event.Rune() {
		case '1':
			if player.CurrentShell != nil && player.Energy >= slash.EnergyCost {
				nextPlayerAction, nextPlayerMeditate, nextPlayerDirectPossess = slash, false, false
			}
		case '2':
			if player.CurrentShell != nil && player.Energy >= heavyStrike.EnergyCost {
				nextPlayerAction, nextPlayerMeditate, nextPlayerDirectPossess = heavyStrike, false, false
			}
		case 'm':
			if player.CurrentShell != nil {
				nextPlayerAction, nextPlayerMeditate, nextPlayerDirectPossess = nil, true, false
			}
		case 'x':
			target := enemies[currentTargetIndex]
			if target.CurrentShell != nil && target.CurrentShell.IsDefeated() && player.Energy >= directPossessionCost {
				nextPlayerAction, nextPlayerMeditate, nextPlayerDirectPossess = nil, false, true
			}
		}

		return event
	})

	if err := app.SetRoot(mainLayout, true).SetFocus(mainLayout).Run(); err != nil {
		panic(err)
	}
}
