package service

import (
	"demofine/data"
	"demofine/internal/utils"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (s *Service) WelcomeScreen(_ fyne.Window) fyne.CanvasObject {
	logo := canvas.NewImageFromResource(data.FyneLogoTransparent)
	logo.FillMode = canvas.ImageFillContain
	if fyne.CurrentDevice().IsMobile() {
		logo.SetMinSize(fyne.NewSize(192, 192))
	} else {
		logo.SetMinSize(fyne.NewSize(256, 256))
	}

	return container.NewCenter(container.NewVBox(
		widget.NewLabelWithStyle("Welcome to the Fyne toolkit demo app", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		logo,
		container.NewHBox(
			widget.NewHyperlink("fyne.io", utils.ParseURL("https://fyne.io/")),
			widget.NewLabel("-"),
			widget.NewHyperlink("documentation", utils.ParseURL("https://developer.fyne.io/")),
			widget.NewLabel("-"),
			widget.NewHyperlink("sponsor", utils.ParseURL("https://fyne.io/sponsor/")),
		),
		widget.NewLabel(""),
	))
}
