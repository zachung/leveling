package main

import (
	"fmt"
	"strings"

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
}

// NewPlayer 是一個工廠函數，用於創建一個新的玩家實例
func NewPlayer(name string, health, energy int) *Player {
	return &Player{
		Name:      name,
		Health:    health,
		MaxHealth: health,
		Energy:    energy,
		MaxEnergy: energy,
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
	status.WriteString(fmt.Sprintf("[red]生命: %d / %d[-:-:-]\n", p.Health, p.MaxHealth))  // 紅色顯示生命
	status.WriteString(fmt.Sprintf("[blue]能量: %d / %d[-:-:-]\n", p.Energy, p.MaxEnergy)) // 藍色顯示能量
	return status.String()
}

func main() {
	// --- 遊戲設定 ---
	// 定義所有可用的招式
	slash := &AttackMove{Name: "揮砍", EnergyCost: 10, Damage: 15}
	heavyStrike := &AttackMove{Name: "強力一擊", EnergyCost: 35, Damage: 45}
	stomp := &AttackMove{Name: "踐踏", EnergyCost: 5, Damage: 10}

	// 建立玩家和敵人，傳入初始生命和初始能量
	player := NewPlayer("英雄", 100, 50)
	monster := NewPlayer("哥布林", 80, 20)

	// 為敵人設定預設攻擊招式
	monster.EquipAttack(stomp)

	// --- TUI 介面設定 ---
	app := tview.NewApplication()

	// 戰鬥日誌歷史記錄
	var logHistory []string
	const maxLogLines = 100

	// 建立顯示元件
	playerStatus := tview.NewTextView().SetDynamicColors(true).SetTextAlign(tview.AlignCenter)
	playerStatus.SetBorder(true).SetTitle("你的狀態")

	monsterStatus := tview.NewTextView().SetDynamicColors(true).SetTextAlign(tview.AlignCenter)
	monsterStatus.SetBorder(true).SetTitle("敵人狀態")

	battleLog := tview.NewTextView().SetDynamicColors(true).SetScrollable(true)
	battleLog.SetBorder(true).SetTitle("戰鬥日誌 (可用方向鍵捲動)")

	instructions := tview.NewTextView().SetDynamicColors(true)
	instructions.SetText(
		fmt.Sprintf("[yellow](1) %s [white](耗%d傷%d) | [yellow](2) %s [white](耗%d傷%d) | [yellow](m) %s [white]| [yellow](q)uit",
			slash.Name, slash.EnergyCost, slash.Damage,
			heavyStrike.Name, heavyStrike.EnergyCost, heavyStrike.Damage,
			"冥想"),
	)

	// 更新狀態畫面的函式
	updateStatusViews := func() {
		playerStatus.SetText(player.GetStatusText())
		monsterStatus.SetText(monster.GetStatusText())
	}

	// 初始畫面
	updateStatusViews()
	logHistory = append(logHistory, "戰鬥開始！")
	battleLog.SetText(strings.Join(logHistory, "\n"))

	// 版面配置
	flex := tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(playerStatus, 0, 1, false).
			AddItem(monsterStatus, 0, 1, false), 0, 1, false).
		AddItem(battleLog, 0, 2, false)

	mainLayout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(flex, 0, 1, false).
		AddItem(instructions, 2, 0, false)

	// --- 輸入處理 ---
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if player.IsDefeated() || monster.IsDefeated() {
			if event.Rune() == 'q' {
				app.Stop()
			}
			return event
		}

		var currentTurnLogs []string
		playerActionTaken := true

		switch event.Rune() {
		case '1':
			currentTurnLogs = append(currentTurnLogs, player.Attack(monster, slash)...)
		case '2':
			currentTurnLogs = append(currentTurnLogs, player.Attack(monster, heavyStrike)...)
		case 'm':
			currentTurnLogs = append(currentTurnLogs, player.Meditate()...)
		case 'q':
			app.Stop()
			return event
		default:
			playerActionTaken = false // 如果按了無效鍵，則不算玩家行動
		}

		// 如果玩家有行動，才輪到怪物行動
		if playerActionTaken {
			if monster.IsDefeated() {
				currentTurnLogs = append(currentTurnLogs, "", "[::b][green]恭喜！你擊敗了哥布林！ 按(q)離開。")
			} else {
				// 怪物回合 (敵人使用預設招式)
				currentTurnLogs = append(currentTurnLogs, "") // 加入空行
				currentTurnLogs = append(currentTurnLogs, monster.Attack(player, monster.AI_Attack)...)
				if player.IsDefeated() {
					currentTurnLogs = append(currentTurnLogs, "", "[::b][red]你被哥布林擊敗了... 按(q)離開。")
				}
			}

			logHistory = append(logHistory, currentTurnLogs...)
			if len(logHistory) > maxLogLines {
				logHistory = logHistory[len(logHistory)-maxLogLines:]
			}

			// 更新日誌畫面並捲動到底部
			battleLog.SetText(strings.Join(logHistory, "\n"))
			battleLog.ScrollToEnd()

			updateStatusViews()
		}
		return event
	})

	if err := app.SetRoot(mainLayout, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
