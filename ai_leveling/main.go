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
var nextPlayerPossess bool

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
	if p.Energy < move.EnergyCost {
		return []string{fmt.Sprintf("能量不足以使用 [%s]！", move.Name)}
	}

	logs := []string{fmt.Sprintf("➡️ %s 使用 [%s] 攻擊 %s！", p.Name, move.Name, target.Name)}
	p.LoseEnergy(move.EnergyCost)
	logs = append(logs, fmt.Sprintf("   %s 消耗了 %d 點能量。", p.Name, move.EnergyCost))

	if target.CurrentShell != nil {
		finalDamage := move.Damage + p.CurrentShell.Strength
		target.CurrentShell.LoseHealth(finalDamage)
		logs = append(logs, fmt.Sprintf("   對 %s 的軀殼造成了 %d 點傷害！ (%d 基礎 + %d 力量)", target.Name, finalDamage, move.Damage, p.CurrentShell.Strength))
	}
	return logs
}

// Meditate 讓玩家冥想以恢復能量
func (p *Player) Meditate() []string {
	restoreAmount := 20
	p.GainEnergy(restoreAmount)
	return []string{fmt.Sprintf("🧘 %s 進行冥想，恢復了 %d 點能量。", p.Name, restoreAmount)}
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
		}
	} else {
		status.WriteString("[purple]狀態: 靈體[-:-:-]\n")
		if nextPlayerPossess {
			status.WriteString("[cyan]預約: 附身[-:-:-]")
		}
	}
	return status.String()
}

// GetEnemyStatusText 獲取單一敵人狀態文字
func (p *Player) GetEnemyStatusText() string {
	if p.CurrentShell == nil {
		return fmt.Sprintf("[::b]%s\n\n[gray]已被摧毀[-:-:-]", p.Name)
	}
	var status strings.Builder
	status.WriteString(fmt.Sprintf("[::b]%s\n", p.Name))
	status.WriteString(fmt.Sprintf("%s\n", strings.Repeat("─", len(p.Name)+4)))
	status.WriteString(fmt.Sprintf("[red]生命: %d / %d[-:-:-]\n", p.CurrentShell.Health, p.CurrentShell.MaxHealth))
	status.WriteString(fmt.Sprintf("[orange]力量: %d[-:-:-]", p.CurrentShell.Strength))
	return status.String()
}

func main() {
	// --- 遊戲設定 ---
	slash := &AttackMove{Name: "揮砍", EnergyCost: 10, Damage: 15}
	heavyStrike := &AttackMove{Name: "強力一擊", EnergyCost: 35, Damage: 45}
	stomp := &AttackMove{Name: "踐踏", EnergyCost: 1, Damage: 8}
	bite := &AttackMove{Name: "啃咬", EnergyCost: 1, Damage: 12}
	possessionCost := 40

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

	// 將 ChangedFunc 先宣告為一個變數，方便之後移除和加回
	var enemyListChanged func(int, string, string, rune)

	// 完整的畫面更新函式
	updateAllViews := func() {
		playerStatus.SetText(player.GetPlayerStatusText())
		targetStatus.SetText(enemies[currentTargetIndex].GetEnemyStatusText())

		// 修正：在更新列表前，先移除回呼函式
		enemyList.SetChangedFunc(nil)

		enemyList.Clear()
		for i, enemy := range enemies {
			var status string
			if enemy.CurrentShell == nil {
				status = "[gray]已被摧毀"
			} else {
				status = fmt.Sprintf("生命: %d/%d", enemy.CurrentShell.Health, enemy.CurrentShell.MaxHealth)
			}
			mainText := fmt.Sprintf("%s %s", enemy.Name, status)
			if i == currentTargetIndex {
				mainText = "[red]>> " + mainText + "[-:-:-]"
			}
			enemyList.AddItem(mainText, "", 0, nil)
		}
		enemyList.SetCurrentItem(currentTargetIndex)

		// 修正：更新完列表後，再將回呼函式加回去
		enemyList.SetChangedFunc(enemyListChanged)

		// 更新指令提示
		if player.CurrentShell != nil {
			instructions.SetText(fmt.Sprintf("[yellow](1) %s | (2) %s | (m) %s | (Tab)切換 | (q)uit", slash.Name, heavyStrike.Name, "冥想"))
		} else {
			instructions.SetText(fmt.Sprintf("[yellow](p) 附身 (消耗 %d 能量) | (Tab)切換 | (q)uit", possessionCost))
		}
	}

	logHistory = append(logHistory, "戰鬥開始！")
	updateAllViews() // 初始繪製

	rightPanel := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(enemyList, 0, 1, true).
		AddItem(targetStatus, 10, 0, false)

	mainFlex := tview.NewFlex().
		AddItem(playerStatus, 0, 1, false).
		AddItem(rightPanel, 0, 1, true)

	mainLayout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(mainFlex, 0, 1, true).
		AddItem(battleLog, 12, 0, false).
		AddItem(instructions, 1, 0, false)

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

			// 玩家行動邏輯
			if player.CurrentShell != nil {
				if player.CurrentShell.IsDefeated() {
					player.CurrentShell = nil
					logsThisTick = append(logsThisTick, "[orange]你的軀殼被摧毀了！你現在是靈體狀態。[-:-:-]")
					actionTaken = true
				} else if time.Now().After(player.CurrentShell.Cooldown) {
					target := enemies[currentTargetIndex]
					if nextPlayerAction != nil && target.CurrentShell != nil {
						logsThisTick = append(logsThisTick, player.Attack(target, nextPlayerAction)...)
						player.CurrentShell.Cooldown = time.Now().Add(cooldownDuration)
						nextPlayerAction = nil
						actionTaken = true
					} else if nextPlayerMeditate {
						logsThisTick = append(logsThisTick, player.Meditate()...)
						player.CurrentShell.Cooldown = time.Now().Add(cooldownDuration)
						nextPlayerMeditate = false
						actionTaken = true
					}
				}
			} else {
				player.GainEnergy(1)
				if nextPlayerPossess {
					player.CurrentShell = NewShell("人類軀殼", 100, 5, nil)
					player.LoseEnergy(possessionCost)
					logsThisTick = append(logsThisTick, "[green]你消耗能量附身到新的軀殼上！[-:-:-]")
					nextPlayerPossess = false
					actionTaken = true
				}
			}

			// 敵人 AI 行動邏輯
			allEnemiesDefeated := true
			for _, enemy := range enemies {
				if enemy.CurrentShell != nil {
					if enemy.CurrentShell.IsDefeated() {
						enemy.CurrentShell = nil
						logsThisTick = append(logsThisTick, fmt.Sprintf("[red]%s 的軀殼已被摧毀！[-:-:-]", enemy.Name))
						actionTaken = true
					} else {
						allEnemiesDefeated = false
						if time.Now().After(enemy.CurrentShell.Cooldown) && player.CurrentShell != nil {
							logsThisTick = append(logsThisTick, "")
							logsThisTick = append(logsThisTick, enemy.Attack(player, enemy.CurrentShell.AI_Attack)...)
							enemy.CurrentShell.Cooldown = time.Now().Add(time.Duration(20+len(enemies)) * 100 * time.Millisecond)
							actionTaken = true
						}
					}
				}
			}

			if allEnemiesDefeated && !gameIsOver {
				logsThisTick = append(logsThisTick, "", "[::b][green]勝利！你擊敗了所有敵人！ 按(q)離開。")
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

		if player.CurrentShell != nil {
			switch event.Rune() {
			case '1':
				if player.Energy >= slash.EnergyCost {
					nextPlayerAction = slash
					nextPlayerMeditate = false
				}
			case '2':
				if player.Energy >= heavyStrike.EnergyCost {
					nextPlayerAction = heavyStrike
					nextPlayerMeditate = false
				}
			case 'm':
				nextPlayerAction = nil
				nextPlayerMeditate = true
			}
		} else {
			if event.Rune() == 'p' && player.Energy >= possessionCost {
				nextPlayerPossess = true
			}
		}

		return event
	})

	if err := app.SetRoot(mainLayout, true).SetFocus(mainLayout).Run(); err != nil {
		panic(err)
	}
}
