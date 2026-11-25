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

// AttackData holds the parsed JSON payload returned by the server during a
// combat turn, including current and maximum health points, damage dealt,
// damage type, and combat status.
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

// BattleUI manages the arena combat interface, including fighter avatars,
// health bars, action buttons, and the combat log.
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
	npcHealth		*HealthBar
	playerHealth	*HealthBar
}

// NewBattleUI builds and initializes the arena combat layout with player and
// NPC displays, custom health bars, and combat action controls.
func NewBattleUI(player *Player, data *Data) *BattleUI {
	playerImg := canvas.NewImageFromFile("./data/npc/" + player.gender + "attack.png")
	playerImg.FillMode = canvas.ImageFillContain
	playerImg.SetMinSize(fyne.NewSize(300, 400))

	npcImg := &canvas.Image{}
	npcImg.FillMode = canvas.ImageFillContain
	npcImg.SetMinSize(fyne.NewSize(300, 400))

	b := &BattleUI{
		playerImg:			playerImg,
		npcImg:				npcImg,
		playerHealth:		NewHealthBar(float64(player.health)),
		npcHealth:			NewHealthBar(0),
		playerNameLabel:	widget.NewLabel(player.name),
		npcNameLabel:		widget.NewLabel("enemy"),
		battleLog:			widget.NewLabel(""),
		player:				player,
		data:				data,
	}

	playerSide := container.NewVBox(
		b.playerNameLabel,
		b.playerHealth.Content(),
		b.playerImg,
	)

	npcSide := container.NewVBox(
		b.npcNameLabel,
		b.npcHealth.Content(),
		b.npcImg,
	)

	fighters := container.NewHBox(
		layout.NewSpacer(),
		playerSide,
		layout.NewSpacer(),
		npcSide,
		layout.NewSpacer(),
	)

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

// Content returns the root CanvasObject of the battle screen.
func (b *BattleUI) Content() fyne.CanvasObject {
	return b.content
}

// Update refreshes the player and NPC health bars and updates the combat log
// with the latest turn's damage details.
func (b *BattleUI) Update(atk AttackData) {
	b.player.health = atk.PlayerHP

	b.playerHealth.SetMax(float64(atk.PlayerMax_HP))
	b.playerHealth.SetValue(float64(atk.PlayerHP))

	b.npcHealth.SetMax(float64(atk.NPCMaxHP))
	b.npcHealth.SetValue(float64(atk.NPCHP))
	b.battleLog.SetText(fmt.Sprintf("You deal %d damage (%s) to %s !", atk.Damage, atk.DamageType, atk.NPCName))
}

// InitFight sets the combat state to active, loads the player and enemy
// sprites and names, and initializes the combat UI with the first turn data.
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

// EndFight marks the current combat session as finished.
func (b *BattleUI) EndFight() {
	b.inFight = false
}

// parse unmarshals a raw JSON attack response string into the AttackData struct.
func (a *AttackData) parse(msg string) error {
	rawJSON := strings.TrimSpace(msg)
	if err := json.Unmarshal([]byte(rawJSON), a); err != nil {
		fmt.Println("Parsing attack data error")
		return err
	}
	return nil
}
