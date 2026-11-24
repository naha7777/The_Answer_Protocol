package main

import (
	"encoding/json"
	"fmt"
	// "image/color"
	// "os"
	"strings"
	// "the_answer_protocol/client"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type AttackData struct {
	PlayerHP		int		`json:"player_hp"`
	PlayerMax_HP	int		`json:"player_max_hp"`
	NPCHP			int		`json:"NPC_hp"`
	NPCMaxHP		int		`json:"NPC_max_hp"`
	Damage			int		`json:"damage"`
	NPCName			string	`json:"NPC_name"`
	DamageType		string	`json:"damage_type"`
	Status			string	`json:"status"`
}

type BattleUI struct {
	inFight			bool
	data			*Data
	player			*Player
	playerImg		*canvas.Image
	npcImg			*canvas.Image
	npcNameLabel	*widget.Label
	battleLog		*widget.Label
	content			fyne.CanvasObject
	npcHealth		*widget.ProgressBar
	playerHealth	*widget.ProgressBar
}


func NewBattleUI(player *Player, data *Data) *BattleUI {
	playerImg := canvas.NewImageFromFile(player.battleAvatar)
	playerImg.FillMode = canvas.ImageFillContain
	playerImg.SetMinSize(fyne.NewSize(300, 400))

	npcImg := &canvas.Image{}
	npcImg.FillMode = canvas.ImageFillContain
	npcImg.SetMinSize(fyne.NewSize(300, 400))

	b := &BattleUI{
		playerImg:		playerImg,
		npcImg:			npcImg,
		playerHealth:	widget.NewProgressBar(),
		npcHealth:		widget.NewProgressBar(),
		npcNameLabel:	widget.NewLabel("enemy"),
		battleLog:		widget.NewLabel(""),
		player:			player,
		data:			data,
	}

	playerSide := container.NewVBox(
		widget.NewLabel(player.name),
		b.playerHealth,
		b.playerImg,
	)

	npcSide := container.NewVBox(
		b.npcNameLabel,
		b.npcHealth,
		b.npcImg,
	)

	fighters := container.NewHBox(
		layout.NewSpacer(),
		playerSide,
		layout.NewSpacer(),
		npcSide,
		layout.NewSpacer(),
	)

	b.playerHealth.Min = 0
	b.playerHealth.Max = 100
	b.npcHealth.Min = 0
	b.npcHealth.Min = 100

	attackBtn := widget.NewButton("ATTACK", func() {b.data.conn.Send("ATTACK")})
	fleeBtn := widget.NewButton("FLEE", func() {b.data.conn.Send("FLEE")})

	bottomBar := container.NewVBox(
		container.NewCenter(b.battleLog),
		container.NewCenter(container.NewHBox(attackBtn, fleeBtn)),
	)

	// actions := container.NewHBox(attackBtn, fleeBtn)
	// bars := container.NewVBox(
	// 	widget.NewLabel(player.name),
	// 	b.playerHealth,
	// 	b.npcNameLabel,
	// 	b.npcHealth,
	// )

	// b.content = container.NewBorder(bars, actions, nil, nil ,b.battleLog)
	b.content = container.NewBorder(nil, bottomBar, nil, nil, container.NewCenter(fighters))
	return b
}

func (b *BattleUI) Content() fyne.CanvasObject {
	return b.content
}

// func (b *BattleUI) StartFight(atk AttackData) {
// 	cleanID := strings.TrimPrefix(npcID, "npc.")
// 	b.npcNameLabel.SetText(IDToName(npcID))

// 	b.npcImg.File = "./data/npc/" + cleanID + ".png"
// 	b.npcImg.Refresh()

// 	b.npcHealth.Max = npcMaxHP
// 	b.npcHealth.SetValue(npcMaxHP)
// 	b.playerHealth.SetValue(float64(b.player.health))
// 	b.battleLog.SetText("Battle is starting !!!")
// }

func (b *BattleUI) Update(atk AttackData) {
	b.player.health = atk.PlayerHP

	b.playerHealth.Max = float64(atk.PlayerMax_HP)
	b.playerHealth.SetValue(float64(atk.PlayerHP))

	b.npcHealth.Max = float64(atk.NPCMaxHP)
	b.npcHealth.SetValue(float64(atk.NPCHP))
	b.battleLog.SetText(fmt.Sprintf("You deal %d damage (%s) to %s !", atk.Damage, atk.DamageType, atk.NPCName))
}

func (b *BattleUI) InitFight(atk AttackData) {
	b.inFight = true
	b.npcNameLabel.SetText(atk.NPCName)
	b.playerImg.File = "./data/npc/" + b.player.gender + "attack.png"
	b.playerImg.Refresh()

	npcFile := strings.ToLower(strings.ReplaceAll(atk.NPCName, " ", "_"))
	b.npcImg.File = "./data/npc/" + npcFile + "attack.png"
	b.npcImg.Refresh()

	b.npcNameLabel.SetText(atk.NPCName)
	b.Update(atk)
}

func (b *BattleUI) EndFight() {
	b.inFight = false
}

func (*AttackData) atkData(msg string) {
	rawJSON := strings.TrimSpace(msg)
	if rawJSON == "" {
		return
	}
	var data AttackData
	if err := json.Unmarshal([]byte(rawJSON), &data); err != nil {
		fmt.Println("Parsing attack data error")
		return
	}
}