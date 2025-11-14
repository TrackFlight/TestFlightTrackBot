package graphics

import (
	"bytes"
	"image"
	"image/color"

	"github.com/TrackFlight/TestFlightTrackBot/internal/utils"
	"github.com/fogleman/gg"
)

type TemporaryText struct {
	rect *utils.Rect
	text string
}

func DrawText(dc *gg.Context, tr *utils.TextRect) {
	face := tr.Font()
	dc.SetFontFace(face)

	m := face.Metrics()
	ascent := float64(m.Ascent) / 64.0
	descent := float64(m.Descent) / 64.0
	lineHeight := ascent + descent
	baselineY := float64(tr.Y()) + (float64(tr.Height())+lineHeight)/2 - descent
	dc.DrawString(tr.String(), float64(tr.X()), baselineY)
}

func DrawRoundedImage(dc *gg.Context, src []byte, x, y, w, h, r float64) error {
	img, _, err := image.Decode(bytes.NewReader(src))
	if err != nil {
		return err
	}
	temp := gg.NewContext(int(w), int(h))
	temp.DrawRoundedRectangle(0, 0, w, h, r)
	temp.Clip()
	temp.DrawImage(img, 0, 0)
	dc.DrawImage(temp.Image(), int(x), int(y))
	return nil
}

func ColorDodge(base, blend image.Image) *image.RGBA {
	b := base.Bounds()
	out := image.NewRGBA(b)

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {

			br, bg, bb, ba := base.At(x, y).RGBA()
			tr, tg, tb, ta := blend.At(x, y).RGBA()

			rb := float64(br) / 65535.0
			gb := float64(bg) / 65535.0
			bb2 := float64(bb) / 65535.0

			rt := float64(tr) / 65535.0
			gt := float64(tg) / 65535.0
			bt := float64(tb) / 65535.0

			alphaOverlay := float64(ta) / 65535.0
			alphaBase := float64(ba) / 65535.0

			dr := dodge(rb, rt)
			dg := dodge(gb, gt)
			db := dodge(bb2, bt)

			cr := lerp(rb, dr, alphaOverlay)
			cg := lerp(gb, dg, alphaOverlay)
			cb := lerp(bb2, db, alphaOverlay)

			a := alphaOverlay + alphaBase*(1-alphaOverlay)

			out.Set(x, y, color.NRGBA{
				R: uint8(cr * 255),
				G: uint8(cg * 255),
				B: uint8(cb * 255),
				A: uint8(a * 255),
			})
		}
	}
	return out
}

func dodge(base, blend float64) float64 {
	if blend >= 1.0 {
		return 1.0
	}
	return clamp(base / (1.0 - blend))
}

func lerp(a, b, t float64) float64 { return a + (b-a)*t }

func clamp(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}
