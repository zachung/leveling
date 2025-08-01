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

// Player 定義了遊戲中的玩家角色
type Player struct {
	Name      string      // 玩家名稱
	Health    int         // 當前生命
	MaxHealth int         // 最大生命
	Energy    int         // 當前能量
	MaxEnergy int         // 最大能量
	AI_Attack *AttackMove // 電腦(敵人)AI使用的攻擊招式
	Cooldown  time.Time   // 行動冷卻時間
}

// NewPlayer 是一個工廠函數，用於創建一個新的玩家實例
func NewPlayer(name string, health, energy int) *Player {
	return &Player{
		Name:      name,
		Health:    health,
		MaxHealth: health,
		Energy:    energy,
		MaxEnergy: energy,
		Cooldown:  time.Now(), // 初始狀態為可立即行動
	}
}

// EquipAttack 讓玩家裝備一個攻擊招式 (主要供敵人AI使用)
func (p *Player) EquipAttack(attack *AttackMove) {
	p.AI_Attack = attack
}

// LoseHealth 減少玩家的生命，但不會低於 0
func (p *Player) LoseHealth(amount int) {
	p.Health -= amount
	if p.Health < 0 {
		p.Health = 0
	}
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

// IsDefeated 檢查玩家是否被擊敗 (生命為 0)
func (p *Player) IsDefeated() bool {
	return p.Health <= 0
}

// Attack 讓玩家對目標使用指定的攻擊招式，並返回戰鬥日誌
func (p *Player) Attack(target *Player, move *AttackMove) []string {
	var logs []string
	if move == nil {
		logs = append(logs, fmt.Sprintf("%s 沒有選擇任何招式！", p.Name))
		return logs
	}

	if p.Energy < move.EnergyCost {
		logs = append(logs, fmt.Sprintf("%s 想要使用 [%s]，但是能量不足！", p.Name, move.Name))
		return logs
	}

	// 執行攻擊
	logs = append(logs, fmt.Sprintf("➡️ %s 使用 [%s] 攻擊 %s！", p.Name, move.Name, target.Name))
	p.LoseEnergy(move.EnergyCost)
	logs = append(logs, fmt.Sprintf("   %s 消耗了 %d 點能量。", p.Name, move.EnergyCost))

	target.LoseHealth(move.Damage) // 修正：對目標造成生命傷害
	logs = append(logs, fmt.Sprintf("   %s 對 %s 造成了 %d 點生命傷害！", p.Name, target.Name, move.Damage))
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
	status.WriteString(fmt.Sprintf("[::b]%s\n", p.Name)) // 粗體名稱
	status.WriteString(fmt.Sprintf("%s\n", strings.Repeat("─", len(p.Name)+4)))
	status.WriteString(fmt.Sprintf("[red]生命: %d / %d[-:-:-]\n", p.Health, p.MaxHealth))
	status.WriteString(fmt.Sprintf("[blue]能量: %d / %d[-:-:-]\n", p.Energy, p.MaxEnergy))
	if time.Now().Before(p.Cooldown) {
		status.WriteString(fmt.Sprintf("[yellow]狀態: 冷卻中 (%.1fs)[-:-:-]", time.Until(p.Cooldown).Seconds()))
	} else {
		status.WriteString("[green]狀態: 可行動[-:-:-]")
	}
	return status.String()
}

func main() {
	// --- 遊戲設定 ---
	slash := &AttackMove{Name: "揮砍", EnergyCost: 10, Damage: 15}
	heavyStrike := &AttackMove{Name: "強力一擊", EnergyCost: 35, Damage: 45}
	stomp := &AttackMove{Name: "踐踏", EnergyCost: 1, Damage: 8}

	player := NewPlayer("英雄", 100, 50)
	monster := NewPlayer("哥布林", 80, 999) // 敵人能量設高，確保能一直攻擊
	monster.EquipAttack(stomp)

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
	instructions.SetText(
		fmt.Sprintf("[yellow](1) %s [white](耗%d傷%d) | [yellow](2) %s [white](耗%d傷%d) | [yellow](m) %s [white]| [yellow](q)uit",
			slash.Name, slash.EnergyCost, slash.Damage, heavyStrike.Name, heavyStrike.EnergyCost, heavyStrike.Damage, "冥想"),
	)
	updateStatusViews := func() {
		playerStatus.SetText(player.GetStatusText())
		monsterStatus.SetText(monster.GetStatusText())
	}
	logHistory = append(logHistory, "戰鬥開始！")
	battleLog.SetText(strings.Join(logHistory, "\n"))

	mainLayout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewFlex().AddItem(playerStatus, 0, 1, false).AddItem(monsterStatus, 0, 1, false), 0, 1, false).
		AddItem(battleLog, 0, 2, false).
		AddItem(instructions, 2, 0, false)

	// --- 遊戲邏輯與主迴圈 ---
	var playerAction *AttackMove
	var playerMeditate bool
	cooldownDuration := 1 * time.Second

	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		for range ticker.C {
			if player.IsDefeated() || monster.IsDefeated() {
				continue
			}

			var logsThisTick []string
			actionTaken := false

			// 玩家行動
			if time.Now().After(player.Cooldown) {
				if playerAction != nil {
					logsThisTick = append(logsThisTick, player.Attack(monster, playerAction)...)
					player.Cooldown = time.Now().Add(cooldownDuration)
					playerAction = nil
					actionTaken = true
				} else if playerMeditate {
					logsThisTick = append(logsThisTick, player.Meditate()...)
					player.Cooldown = time.Now().Add(cooldownDuration)
					playerMeditate = false
					actionTaken = true
				}
			}

			// 敵人 AI 行動
			if time.Now().After(monster.Cooldown) && !monster.IsDefeated() {
				logsThisTick = append(logsThisTick, "")
				logsThisTick = append(logsThisTick, monster.Attack(player, monster.AI_Attack)...)
				monster.Cooldown = time.Now().Add(2 * time.Second) // 讓敵人攻擊慢一點
				actionTaken = true
			}

			if actionTaken {
				if monster.IsDefeated() {
					logsThisTick = append(logsThisTick, "", "[::b][green]恭喜！你擊敗了哥布林！ 按(q)離開。")
				} else if player.IsDefeated() {
					logsThisTick = append(logsThisTick, "", "[::b][red]你被哥布林擊敗了... 按(q)離開。")
				}
				logHistory = append(logHistory, logsThisTick...)
				if len(logHistory) > maxLogLines {
					logHistory = logHistory[len(logHistory)-maxLogLines:]
				}
			}

			// 使用 QueueUpdateDraw 安全地更新 UI
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
		if player.IsDefeated() || monster.IsDefeated() {
			if event.Rune() == 'q' {
				app.Stop()
			}
			return event
		}

		if time.Now().After(player.Cooldown) { // 只在冷卻結束時接受輸入
			switch event.Rune() {
			case '1':
				playerAction = slash
			case '2':
				playerAction = heavyStrike
			case 'm':
				playerMeditate = true
			}
		}

		if event.Rune() == 'q' {
			app.Stop()
		}
		return event
	})

	if err := app.SetRoot(mainLayout, true).SetFocus(mainLayout).Run(); err != nil {
		panic(err)
	}
}
