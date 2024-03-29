package service

import (
	"demofine/data"
	"demofine/internal/utils"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (s *Service) WelcomeScreen(_ fyne.Window, _ string) fyne.CanvasObject {
	logo := canvas.NewImageFromResource(data.RSHULogoTransparent)
	logo.FillMode = canvas.ImageFillContain
	if fyne.CurrentDevice().IsMobile() {
		logo.SetMinSize(fyne.NewSize(192, 192))
	} else {
		logo.SetMinSize(fyne.NewSize(256, 256))
	}

	return container.NewCenter(container.NewVBox(
		container.NewHBox(
			widget.NewHyperlink("Документация", utils.ParseURL("https://github.com/pasha1coil/guiRSHUdata")),
		),
		widget.NewLabel(""),
	))
}
