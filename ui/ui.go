package ui

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/Scalingo/go-scalingo/v7"
)

func Start(ctx context.Context) error {
	p := tea.NewProgram(&model{
		ctx:          ctx,
		applications: createApplicationsView(ctx),
		currentView:  &loadingView{},
	}, tea.WithAltScreen(), tea.WithContext(ctx))

	if _, err := p.Run(); err != nil {
		return err
	}

	return nil
}

type model struct {
	ctx          context.Context
	regions      []scalingo.Region
	currentView  view
	applications view
}

func (m *model) Init() tea.Cmd {
	return tea.Batch(
		tea.SetWindowTitle("Scalingo"),
		m.FetchRegions(),
	)
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.applications.SetSize(msg.Width, msg.Height)

	case []scalingo.Region:
		return m.SetRegions(msg)

	case map[string]*scalingo.Client:
		m.applications.SetRegionClients(msg)
		return m, setCurrentView(m.applications)

	case view:
		m.currentView = msg

	}

	return m, m.currentView.Update(msg)
}

func (m *model) View() string {
	return m.currentView.View()
}

type view interface {
	Update(tea.Msg) tea.Cmd
	View() string
	SetSize(w, h int)
	SetRegionClients(map[string]*scalingo.Client)
}

type initViewMsg struct{}

func setCurrentView(v view) tea.Cmd {
	return tea.Sequence(
		func() tea.Msg {
			return v
		},
		func() tea.Msg {
			return initViewMsg{}
		},
	)
}
