package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

// HealthBar represents a custom UI health bar widget that dynamically changes
// its width and interpolates its fill color from green to orange to red based
// on the remaining health percentage.
type HealthBar struct {
	content	fyne.CanvasObject
	fill	*canvas.Rectangle
	label	*canvas.Text
	max		float64
	current	float64
}

// NewHealthBar creates and initializes a new HealthBar with the specified
// maximum health value, defaulting current health to full.
func NewHealthBar(max float64) *HealthBar {
	// 
	bg := canvas.NewRectangle(color.NRGBA{R: 30, G: 30, B: 30, A: 220})
	bg.StrokeColor = color.White
	bg.StrokeWidth = 2

	fill := canvas.NewRectangle(color.NRGBA{R: 40, G: 120, B: 40, A: 255})

	lbl := canvas.NewText("", color.White)
	lbl.TextStyle = fyne.TextStyle{Bold: true}
	lbl.Alignment = fyne.TextAlignCenter

	barArea := container.NewWithoutLayout(bg, fill, lbl)

	hb := &HealthBar{
		fill:		fill,
		label:		lbl,
		max:		max,
		current:	max,
	}

	sizer := canvas.NewRectangle(color.Transparent)
	sizer.SetMinSize(fyne.NewSize(300, 26))

	hb.content = container.NewStack(sizer, barArea)
	hb.SetValue(max)
	return hb
}

// Content returns the underlying Fyne CanvasObject for layout integration.
func (h *HealthBar) Content() fyne.CanvasObject {
	return h.content
}

// SetMax updates the maximum health capacity of the bar.
// Values less than or equal to zero are clamped to 1 to prevent division by zero.
func (h *HealthBar) SetMax(max float64) {
	if max <= 0 {
		max = 1
	}
	h.max = max
}

// SetValue updates the current health value, clamping it between 0 and max,
// and refreshes the bar width, color gradient, and numeric label.
func (h *HealthBar) SetValue(val float64) {
	if val < 0 {
		val = 0
	}
	if val > h.max {
		val = h.max
	}

	ratio := h.current / h.max
	totalW := float32(300)
	totalH := float32(26)

	barArea := h.content.(*fyne.Container).Objects[1].(*fyne.Container)
	bg := barArea.Objects[0].(*canvas.Rectangle)
	bg.Resize(fyne.NewSize(totalW, totalH))

	h.fill.Resize(fyne.NewSize(totalW*float32(ratio), totalH))
	h.fill.FillColor = healthColor(ratio)
	h.fill.Refresh()

	h.label.Text = fmt.Sprintf("%.0f / %.0f", h.current, h.max)
	h.label.Resize(fyne.NewSize(totalW, totalH))
	h.label.Move(fyne.NewPos(0, 3))
	h.label.Refresh()
}

// healthColor computes an RGBA color using linear interpolation based on the
// health ratio: green (100%) to orange (40%) to red (10% and below).
func healthColor(ratio float64) color.NRGBA {
	green := color.NRGBA{R: 40, G: 210, B: 40, A: 255}
	orange := color.NRGBA{R: 245, G: 140, B: 20, A: 255}
	red := color.NRGBA{R: 220, G: 30, B: 30, A: 255}

	lerp := func(a, b uint8, t float64) uint8 {
		return uint8(float64(a) + (float64(b)-float64(a))*t)
	}

	switch {
	case ratio >= 0.40:
		t := (ratio - 0.40) / 0.60
		return color.NRGBA{
			R: lerp(orange.R, green.R, t),
			G: lerp(orange.G, green.G, t),
			B: lerp(orange.B, green.B, t),
			A: 255,
		}
	case ratio >= 0.10:
		t := (ratio - 0.10) / 0.30
		return color.NRGBA{
			R: lerp(orange.R, green.R, t),
			G: lerp(orange.G, green.G, t),
			B: lerp(orange.B, green.B, t),
			A: 255,
		}
	default:
		return red
	}
}
