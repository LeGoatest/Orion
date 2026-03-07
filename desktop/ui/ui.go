package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// OrionUI manages the desktop interface
type OrionUI struct {
	app    fyne.App
	window fyne.Window
}

// NewOrionUI creates a new OrionUI
func NewOrionUI() *OrionUI {
	a := app.New()
	w := a.NewWindow("Orion Cognitive Runtime")
	w.Resize(fyne.NewSize(1024, 768))

	return &OrionUI{
		app:    a,
		window: w,
	}
}

// Start launches the Fyne application loop
func (ui *OrionUI) Start() {
	// Setup UI layout
	sidebar := ui.createSidebar()
	mainArea := ui.createMainArea()

	content := container.NewHSplit(sidebar, mainArea)
	content.Offset = 0.2

	ui.window.SetContent(content)
	ui.window.ShowAndRun()
}

func (ui *OrionUI) createSidebar() fyne.CanvasObject {
	list := widget.NewList(
		func() int { return 3 }, // Placeholder for workspaces
		func() fyne.CanvasObject { return widget.NewLabel("Workspace") },
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(fmt.Sprintf("Workspace %d", id+1))
		},
	)

	header := widget.NewLabelWithStyle("WORKSPACES", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	return container.NewBorder(header, nil, nil, nil, list)
}

func (ui *OrionUI) createMainArea() fyne.CanvasObject {
	tabs := container.NewAppTabs(
		container.NewTabItem("Agent Workspace", widget.NewLabel("Agent interaction view")),
		container.NewTabItem("Memory Explorer", widget.NewLabel("Knowledge graph visualization")),
		container.NewTabItem("Goal Timeline", widget.NewLabel("Execution history")),
		container.NewTabItem("System Logs", widget.NewMultiLineEntry()),
	)
	return tabs
}
