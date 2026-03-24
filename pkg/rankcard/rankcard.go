// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package rankcard

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"

	"github.com/fogleman/gg"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

// Data holds all information needed to render a rank card.
type Data struct {
	Username  string
	AvatarPNG []byte // raw PNG bytes (nil for default avatar)
	Level     int
	CurrentXP int
	NeededXP  int
	TotalXP   int
	Rank      int
}

const (
	cardWidth  = 934
	cardHeight = 282
	padding    = 24
	avatarSize = 180
)

var (
	colorBG       = color.RGBA{R: 30, G: 33, B: 36, A: 255}
	colorCard     = color.RGBA{R: 43, G: 47, B: 53, A: 255}
	colorAccent   = color.RGBA{R: 88, G: 101, B: 242, A: 255} // Discord blurple
	colorBarBG    = color.RGBA{R: 72, G: 75, B: 80, A: 255}
	colorText     = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	colorTextDim  = color.RGBA{R: 185, G: 187, B: 190, A: 255}
	colorAvatarBG = color.RGBA{R: 54, G: 57, B: 63, A: 255}
)

// Generate creates a rank card PNG image and returns the bytes.
func Generate(data Data) ([]byte, error) {
	dc := gg.NewContext(cardWidth, cardHeight)

	// Background with rounded corners
	dc.SetColor(colorBG)
	dc.Clear()
	drawRoundedRect(dc, 0, 0, cardWidth, cardHeight, 20)
	dc.SetColor(colorCard)
	dc.Fill()

	// Load font faces
	fontBold28, err := loadFontFace(28)
	if err != nil {
		return nil, fmt.Errorf("failed to load font 28: %w", err)
	}
	fontBold22, err := loadFontFace(22)
	if err != nil {
		return nil, fmt.Errorf("failed to load font 22: %w", err)
	}
	fontBold18, err := loadFontFace(18)
	if err != nil {
		return nil, fmt.Errorf("failed to load font 18: %w", err)
	}
	fontBold40, err := loadFontFace(40)
	if err != nil {
		return nil, fmt.Errorf("failed to load font 40: %w", err)
	}

	// Avatar
	avatarX := padding + float64(avatarSize)/2 + 10
	avatarY := float64(cardHeight) / 2
	avatarRadius := float64(avatarSize) / 2

	if data.AvatarPNG != nil {
		avatarImg, err := png.Decode(bytes.NewReader(data.AvatarPNG))
		if err == nil {
			drawCircularImage(dc, avatarImg, avatarX, avatarY, avatarRadius)
		} else {
			drawDefaultAvatar(dc, avatarX, avatarY, avatarRadius)
		}
	} else {
		drawDefaultAvatar(dc, avatarX, avatarY, avatarRadius)
	}

	// Avatar border ring
	dc.SetColor(colorAccent)
	dc.SetLineWidth(4)
	dc.DrawCircle(avatarX, avatarY, avatarRadius+2)
	dc.Stroke()

	// Text area starts after avatar
	textLeft := avatarX + avatarRadius + 30
	textRight := float64(cardWidth) - padding - 10

	// Username (top-left of text area)
	dc.SetFontFace(fontBold28)
	dc.SetColor(colorText)
	name := truncateString(data.Username, 20)
	dc.DrawStringAnchored(name, textLeft, 60, 0, 0.5)

	// Rank badge (top-right)
	dc.SetFontFace(fontBold40)
	dc.SetColor(colorAccent)
	rankText := fmt.Sprintf("#%d", data.Rank)
	dc.DrawStringAnchored(rankText, textRight, 55, 1, 0.5)

	// Level (right of username area)
	dc.SetFontFace(fontBold22)
	dc.SetColor(colorTextDim)
	dc.DrawStringAnchored("LEVEL", textRight-90, 100, 0, 0.5)
	dc.SetFontFace(fontBold40)
	dc.SetColor(colorText)
	dc.DrawStringAnchored(fmt.Sprintf("%d", data.Level), textRight, 100, 1, 0.5)

	// Progress bar
	barY := 170.0
	barHeight := 30.0
	barWidth := textRight - textLeft

	// Bar background
	drawRoundedRect(dc, textLeft, barY, barWidth, barHeight, barHeight/2)
	dc.SetColor(colorBarBG)
	dc.Fill()

	// Bar fill
	progress := 0.0
	if data.NeededXP > 0 {
		progress = float64(data.CurrentXP) / float64(data.NeededXP)
		if progress > 1 {
			progress = 1
		}
	}
	if progress > 0.01 {
		fillWidth := barWidth * progress
		if fillWidth < barHeight {
			fillWidth = barHeight // Minimum width for rounded ends
		}
		drawRoundedRect(dc, textLeft, barY, fillWidth, barHeight, barHeight/2)
		dc.SetColor(colorAccent)
		dc.Fill()
	}

	// XP text (below progress bar)
	dc.SetFontFace(fontBold18)
	dc.SetColor(colorTextDim)
	xpText := fmt.Sprintf("%s / %s XP", formatNumber(data.CurrentXP), formatNumber(data.NeededXP))
	dc.DrawStringAnchored(xpText, textRight, barY+barHeight+25, 1, 0.5)

	// Total XP
	dc.SetColor(colorTextDim)
	totalText := fmt.Sprintf("Total: %s XP", formatNumber(data.TotalXP))
	dc.DrawStringAnchored(totalText, textLeft, barY+barHeight+25, 0, 0.5)

	// Encode to PNG
	var buf bytes.Buffer
	if err := png.Encode(&buf, dc.Image()); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func loadFontFace(size float64) (font.Face, error) {
	f, err := opentype.Parse(fontData)
	if err != nil {
		return nil, err
	}
	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, err
	}
	return face, nil
}

func drawRoundedRect(dc *gg.Context, x, y, w, h, r float64) {
	dc.NewSubPath()
	dc.DrawArc(x+r, y+r, r, math.Pi, 1.5*math.Pi)
	dc.LineTo(x+w-r, y)
	dc.DrawArc(x+w-r, y+r, r, 1.5*math.Pi, 2*math.Pi)
	dc.LineTo(x+w, y+h-r)
	dc.DrawArc(x+w-r, y+h-r, r, 0, 0.5*math.Pi)
	dc.LineTo(x+r, y+h)
	dc.DrawArc(x+r, y+h-r, r, 0.5*math.Pi, math.Pi)
	dc.ClosePath()
}

func drawCircularImage(dc *gg.Context, img image.Image, cx, cy, radius float64) {
	diameter := int(radius * 2)
	resized := gg.NewContext(diameter, diameter)
	resized.DrawImage(img, 0, 0)

	dc.Push()
	dc.DrawCircle(cx, cy, radius)
	dc.Clip()

	bounds := img.Bounds()
	sx := radius * 2 / float64(bounds.Dx())
	sy := radius * 2 / float64(bounds.Dy())
	scale := sx
	if sy > sx {
		scale = sy
	}

	dc.Push()
	dc.Translate(cx, cy)
	dc.Scale(scale, scale)
	dc.DrawImageAnchored(img, 0, 0, 0.5, 0.5)
	dc.Pop()

	dc.Pop()
	dc.ResetClip()
}

func drawDefaultAvatar(dc *gg.Context, cx, cy, radius float64) {
	dc.SetColor(colorAvatarBG)
	dc.DrawCircle(cx, cy, radius)
	dc.Fill()

	dc.SetColor(colorTextDim)
	dc.DrawCircle(cx, cy-radius*0.15, radius*0.3)
	dc.Fill()

	drawRoundedRect(dc, cx-radius*0.4, cy+radius*0.2, radius*0.8, radius*0.5, radius*0.2)
	dc.Fill()
}

func truncateString(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen-1]) + "…"
}

func formatNumber(n int) string {
	if n < 1000 {
		return fmt.Sprintf("%d", n)
	}
	if n < 1000000 {
		return fmt.Sprintf("%d,%03d", n/1000, n%1000)
	}
	return fmt.Sprintf("%d,%03d,%03d", n/1000000, (n%1000000)/1000, n%1000)
}
