package ui

import (
	tea "github.com/charmbracelet/bubbletea"

	scalingo "github.com/Scalingo/go-scalingo/v7"
)

type loadingView struct{}

var _ view = &loadingView{}

func (l *loadingView) SetSize(w, h int) {}

func (l *loadingView) Update(msg tea.Msg) tea.Cmd {

	switch msg := msg.(type) {

	case error:
		return tea.Sequence(
			tea.ExitAltScreen,
			tea.Printf("error: %s", msg),
			tea.Quit,
		)

	}

	return nil
}

func (l *loadingView) View() string {
	return "Loading..."
}

func (l *loadingView) SetRegionClients(map[string]*scalingo.Client) {}
