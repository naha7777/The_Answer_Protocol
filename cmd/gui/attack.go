package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"strings"

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
	playerNameLabel	*widget.Label
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
		playerImg:			playerImg,
		npcImg:				npcImg,
		playerHealth:		widget.NewProgressBar(),
		npcHealth:			widget.NewProgressBar(),
		playerNameLabel:	widget.NewLabel(player.name),
		npcNameLabel:		widget.NewLabel("enemy"),
		battleLog:			widget.NewLabel(""),
		player:				player,
		data:				data,
	}

	playerSide := container.NewVBox(
		b.playerNameLabel,
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
	b.npcHealth.Max = 100

	attackBtn := widget.NewButton("ATTACK", func() {b.data.conn.Send("ATTACK")})
	fleeBtn := widget.NewButton("FLEE", func() {b.data.conn.Send("FLEE")})

	bottomBar := container.NewVBox(
		container.NewCenter(b.battleLog),
		container.NewCenter(container.NewHBox(attackBtn, fleeBtn)),
	)

	topOffset := canvas.NewRectangle(color.Transparent)
	topOffset.SetMinSize(fyne.NewSize(0, 200))

	fightersArea := container.NewVBox(
		topOffset,
		fighters,
	)
	b.content = container.NewBorder(nil, bottomBar, nil, nil, container.NewCenter(fightersArea))
	return b
}

func (b *BattleUI) Content() fyne.CanvasObject {
	return b.content
}

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

	npcFile := strings.ToLower(atk.NPCName)
	b.npcImg.File = "./data/npc/" + npcFile + ".png"
	b.npcImg.Refresh()

	b.npcNameLabel.SetText(atk.NPCName)
	b.playerNameLabel.SetText(b.player.name)
	b.Update(atk)
}

func (b *BattleUI) EndFight() {
	b.inFight = false
}

func (a *AttackData) parse(msg string) error {
	rawJSON := strings.TrimSpace(msg)
	if err := json.Unmarshal([]byte(rawJSON), a); err != nil {
		fmt.Println("Parsing attack data error")
		return err
	}
	return nil
}
