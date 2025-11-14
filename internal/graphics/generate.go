package graphics

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	"image/png"

	"github.com/TrackFlight/TestFlightTrackBot/internal/utils"
	"github.com/disintegration/gift"
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
)

type BannerType string

const (
	Top5       BannerType = "top5"
	Rising     BannerType = "rising"
	HiddenGems BannerType = "hidden_gems"
	Reopened   BannerType = "reopened_slots"
)

//go:embed assets/*
var assetsFS embed.FS

func GenerateBanner(t BannerType, title, subtitle string, images [][]byte) ([]byte, error) {
	data, err := assetsFS.ReadFile(fmt.Sprintf("assets/backgrounds/%s.png", t))
	if err != nil {
		return nil, err
	}

	src, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	data, err = assetsFS.ReadFile("assets/fonts/SFProText-SemiBold.ttf")
	if err != nil {
		return nil, err
	}
	semiBoldFt, err := truetype.Parse(data)
	if err != nil {
		return nil, err
	}

	data, err = assetsFS.ReadFile("assets/fonts/SFProText-Bold.ttf")
	if err != nil {
		return nil, err
	}
	boldFt, err := truetype.Parse(data)
	if err != nil {
		return nil, err
	}

	dc := gg.NewContextForImage(src)

	rect := utils.NewRect(
		dc.Width(),
		dc.Height(),
	)
	rect.SetAlignment(utils.AlignVerticallyCenter | utils.AlignHorizontallyCenter)

	tr1 := rect.AddTextChild(dc, subtitle, semiBoldFt, 148)
	tr2 := rect.AddTextChild(dc, title, boldFt, 365)
	tr2.SetMargin(0, -10, 0, -10)

	rectAppsContainer := rect.AddLayoutChild(
		utils.SizeFitParent,
		utils.SizeAuto,
	)
	rectAppsContainer.SetInset(80)
	rectAppsContainer.SetMargin(0, 62, 0, 0)
	rectAppsContainer.SetAlignment(utils.AlignVerticallyCenter | utils.AlignHorizontallyCenter)

	var rectApps *utils.Rect
	for i := 0; i < len(images); i++ {
		if rectApps == nil || len(rectApps.Children()) >= 3 {
			rectApps = rectAppsContainer.AddLayoutChild(
				utils.SizeFitParent,
				utils.SizeAuto,
			)
			rectApps.SetInset(90)
			rectApps.SetOrientation(utils.OrientationHorizontal)
			rectApps.SetAlignment(utils.AlignHorizontallyCenter)
		}
		rectApps.AddRectChild(480, 480)
	}

	dc.SetRGBA(1, 1, 1, 0.9)
	DrawText(dc, tr1)
	DrawText(dc, tr2)

	dcTemp := gg.NewContext(dc.Width(), dc.Height())
	for x, rectAppsChild := range rectAppsContainer.Children() {
		for y, child := range rectAppsChild.Children() {
			index := x*3 + y
			imgData := images[index]
			err = DrawRoundedImage(dcTemp, imgData, float64(child.X()), float64(child.Y()), float64(child.Width()), float64(child.Height()), 140)
			if err != nil {
				return nil, err
			}
		}
	}
	g := gift.New(gift.GaussianBlur(200))
	dcTempBlurred := image.NewRGBA(src.Bounds())
	g.Draw(dcTempBlurred, dcTemp.Image())

	blended := ColorDodge(dc.Image(), dcTempBlurred)
	dc = gg.NewContextForRGBA(blended)

	dc.DrawImage(dcTemp.Image(), 0, 0)

	var buf bytes.Buffer
	err = png.Encode(&buf, dc.Image())
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
