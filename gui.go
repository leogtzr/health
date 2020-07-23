package main

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/marcusolsson/tui-go"
)

func newInputCommandEntry() *tui.Entry {
	inputCommand := tui.NewEntry()
	inputCommand.SetFocused(true)
	inputCommand.SetSizePolicy(tui.Expanding, tui.Maximum)
	inputCommand.SetEchoMode(tui.EchoModeNormal)

	inputCommandBox := tui.NewHBox(inputCommand)
	inputCommandBox.SetBorder(true)
	inputCommandBox.SetSizePolicy(tui.Expanding, tui.Maximum)
	return inputCommand
}

func newInputCommandBox(input *tui.Entry) *tui.Box {
	inputCommandBox := tui.NewHBox(input)
	inputCommandBox.SetBorder(true)
	inputCommandBox.SetSizePolicy(tui.Expanding, tui.Maximum)
	return inputCommandBox
}

func addCloseApplicationKeyBinding(ui tui.UI) {
	ui.SetKeybinding(closeApplicationKeyBindingAlternative1, func() {
		ui.Quit()
	})
}

func printServiceStatusInScreen(textEdit *tui.TextEdit, ui tui.UI, services *[]Service) {

	var sb strings.Builder
	now := time.Now()
	date := now.Format("2006-01-02 15:04:05")
	sb.WriteString(fmt.Sprintf("%s\n\n", date))

	byInstance := groupByInstance(services)

	for _, k := range sortKeys(byInstance) {
		svcs := (*byInstance)[k]
		sb.WriteString(fmt.Sprintf("==> %s <==\n", k))
		for _, svc := range svcs {
			sb.WriteString(fmt.Sprintf("%s ", svc.ShortString()))
		}
		sb.WriteString("\n")
	}

	textEdit.SetText(sb.String())
	ui.Repaint()
}

func sortKeys(byInstance *map[string][]Service) []string {
	keys := make([]string, 0)
	for k := range *byInstance {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
