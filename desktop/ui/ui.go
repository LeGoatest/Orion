package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type UI struct {
	app fyne.App
	win fyne.Window
}

func NewUI() *UI {
	a := app.New()
	w := a.NewWindow("Orion Cognitive Runtime")
	w.Resize(fyne.NewSize(800, 600))

	// Agent Workspace View
	agentList := widget.NewList(
		func() int { return 6 },
		func() fyne.CanvasObject { return widget.NewLabel("Agent") },
		func(i widget.ListItemID, o fyne.CanvasObject) { o.(*widget.Label).SetText("Agent") },
	)

	// Chat Interface
	chatLog := widget.NewEntry()
	chatLog.MultiLine = true
	chatLog.SetText("Orion Kernel Booted.\nReady for commands...")

	inputField := widget.NewEntry()
	inputField.SetPlaceHolder("Enter command...")

	sendBtn := widget.NewButton("Send", func() {
		// Communicate with kernel through internal/api
		chatLog.SetText(chatLog.Text + "\nUser: " + inputField.Text)
		inputField.SetText("")
	})

	mainContent := container.NewHSplit(
		container.NewVBox(widget.NewLabel("Agents"), agentList),
		container.NewBorder(nil, container.NewHBox(inputField, sendBtn), nil, nil, chatLog),
	)
	mainContent.Offset = 0.2

	w.SetContent(mainContent)
	return &UI{app: a, win: w}
}

func (u *UI) Run() {
	u.win.ShowAndRun()
}
