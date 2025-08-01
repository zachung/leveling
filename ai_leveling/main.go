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
	Cooldown  time.Time
	AI_Attack *AttackMove
}

// Player 代表玩家的核心靈魂
type Player struct {
	Name         string
	Energy       int
	MaxEnergy    int
	CurrentShell *Shell // 玩家當前附身的軀殼，可能為 nil
}

// 預約的行動
var nextPlayerAction *AttackMove
var nextPlayerMeditate bool
var nextPlayerPossess bool

// NewPlayer 創建一個新的玩家靈魂實例
func NewPlayer(name string, energy int) *Player {
	return &Player{
		Name:      name,
		Energy:    energy,
		MaxEnergy: energy,
	}
}

// NewShell 創建一個新的軀殼實例
func NewShell(name string, health int, aiAttack *AttackMove) *Shell {
	return &Shell{
		Name:      name,
		Health:    health,
		MaxHealth: health,
		Cooldown:  time.Now(),
		AI_Attack: aiAttack,
	}
}

// LoseHealth 減少軀殼的生命，但不會低於 0
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

// GainEnergy 為玩家增加能量，但不會超過最大值
func (p *Player) GainEnergy(amount int) {
	p.Energy += amount
	if p.Energy > p.MaxEnergy {
		p.Energy = p.MaxEnergy
	}
}

// LoseEnergy 減少玩家的能量，但不會低於 0
func (p *Player) LoseEnergy(amount int) {
	p.Energy -= amount
	if p.Energy < 0 {
		p.Energy = 0
	}
}

// Attack 讓玩家驅動軀殼對目標使用指定的攻擊招式
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
		target.CurrentShell.LoseHealth(move.Damage)
		logs = append(logs, fmt.Sprintf("   對 %s 的軀殼造成了 %d 點傷害！", target.Name, move.Damage))
	}
	return logs
}

// Meditate 讓玩家冥想以恢復能量
func (p *Player) Meditate() []string {
	restoreAmount := 20
	p.GainEnergy(restoreAmount)
	return []string{fmt.Sprintf("🧘 %s 進行冥想，恢復了 %d 點能量。", p.Name, restoreAmount)}
}

// GetStatusText 獲取格式化後的狀態文字
func (p *Player) GetStatusText() string {
	var status strings.Builder
	status.WriteString(fmt.Sprintf("[::b]%s\n", p.Name))
	status.WriteString(fmt.Sprintf("%s\n", strings.Repeat("─", len(p.Name)+4)))
	status.WriteString(fmt.Sprintf("[blue]能量: %d / %d[-:-:-]\n", p.Energy, p.MaxEnergy))

	if p.CurrentShell != nil {
		status.WriteString(fmt.Sprintf("[red]生命: %d / %d[-:-:-]\n", p.CurrentShell.Health, p.CurrentShell.MaxHealth))
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

func main() {
	// --- 遊戲設定 ---
	slash := &AttackMove{Name: "揮砍", EnergyCost: 10, Damage: 15}
	heavyStrike := &AttackMove{Name: "強力一擊", EnergyCost: 35, Damage: 45}
	stomp := &AttackMove{Name: "踐踏", EnergyCost: 1, Damage: 8}
	possessionCost := 40

	player := NewPlayer("英雄", 100)
	monster := NewPlayer("哥布林", 999) // 敵人能量無限

	player.CurrentShell = NewShell("人類軀殼", 100, nil)
	monster.CurrentShell = NewShell("哥布林軀殼", 80, stomp)

	// --- TUI 介面設定 ---
	app := tview.NewApplication()
	var logHistory []string
	const maxLogLines = 100

	playerStatus := tview.NewTextView().SetDynamicColors(true).SetTextAlign(tview.AlignCenter)
	playerStatus.SetBorder(true).SetTitle("你的狀態")
	monsterStatus := tview.NewTextView().SetDynamicColors(true).SetTextAlign(tview.AlignCenter)
	monsterStatus.SetBorder(true).SetTitle("敵人狀態")
	battleLog := tview.NewTextView().SetDynamicColors(true).SetScrollable(true)
	battleLog.SetBorder(true).SetTitle("戰鬥日誌 (可用方向鍵捲動)")
	instructions := tview.NewTextView().SetDynamicColors(true)

	updateStatusViews := func() {
		playerStatus.SetText(player.GetStatusText())
		monsterStatus.SetText(monster.GetStatusText())
		if player.CurrentShell != nil {
			instructions.SetText(fmt.Sprintf("[yellow](1) %s | (2) %s | (m) %s | (q)uit", slash.Name, heavyStrike.Name, "冥想"))
		} else {
			instructions.SetText(fmt.Sprintf("[yellow](p) 附身 (消耗 %d 能量) | (q)uit", possessionCost))
		}
	}
	logHistory = append(logHistory, "戰鬥開始！")

	mainLayout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewFlex().AddItem(playerStatus, 0, 1, false).AddItem(monsterStatus, 0, 1, false), 0, 1, false).
		AddItem(battleLog, 0, 2, false).
		AddItem(instructions, 1, 0, false)

	// --- 遊戲邏輯與主迴圈 ---
	cooldownDuration := 1 * time.Second

	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		for range ticker.C {
			var logsThisTick []string
			actionTaken := false

			// 玩家行動邏輯
			if player.CurrentShell != nil { // 有軀殼時
				if player.CurrentShell.IsDefeated() {
					player.CurrentShell = nil
					logsThisTick = append(logsThisTick, "[orange]你的軀殼被摧毀了！你現在是靈體狀態。[-:-:-]")
					actionTaken = true
				} else if time.Now().After(player.CurrentShell.Cooldown) {
					if nextPlayerAction != nil {
						logsThisTick = append(logsThisTick, player.Attack(monster, nextPlayerAction)...)
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
			} else { // 靈體狀態邏輯
				player.GainEnergy(1) // 每 100ms 恢復 1 點能量
				if nextPlayerPossess {
					player.CurrentShell = NewShell("人類軀殼", 100, nil)
					player.LoseEnergy(possessionCost)
					logsThisTick = append(logsThisTick, "[green]你消耗能量附身到新的軀殼上！[-:-:-]")
					nextPlayerPossess = false
					actionTaken = true
				}
			}

			// 敵人 AI 行動邏輯
			if monster.CurrentShell != nil && time.Now().After(monster.CurrentShell.Cooldown) {
				if monster.CurrentShell.IsDefeated() {
					monster.CurrentShell = nil // 怪物也可以被摧毀
					logsThisTick = append(logsThisTick, "[::b][green]恭喜！你摧毀了哥布林的軀殼！")
				} else if player.CurrentShell != nil { // 只有在玩家有軀殼時才攻擊
					logsThisTick = append(logsThisTick, "")
					logsThisTick = append(logsThisTick, monster.Attack(player, monster.CurrentShell.AI_Attack)...)
					monster.CurrentShell.Cooldown = time.Now().Add(2 * time.Second)
					actionTaken = true
				}
			}

			if actionTaken {
				logHistory = append(logHistory, logsThisTick...)
				if len(logHistory) > maxLogLines {
					logHistory = logHistory[len(logHistory)-maxLogLines:]
				}
			}

			app.QueueUpdateDraw(func() {
				updateStatusViews()
				if actionTaken {
					battleLog.SetText(strings.Join(logHistory, "\n"))
					battleLog.ScrollToEnd()
				}
			})
		}
	}()

	// --- 輸入處理 ---
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' {
			app.Stop()
			return event
		}

		if player.CurrentShell != nil { // 有軀殼時的輸入
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
		} else { // 靈體狀態時的輸入
			if event.Rune() == 'p' {
				if player.Energy >= possessionCost {
					nextPlayerPossess = true
				}
			}
		}
		return event
	})

	if err := app.SetRoot(mainLayout, true).SetFocus(mainLayout).Run(); err != nil {
		panic(err)
	}
}
