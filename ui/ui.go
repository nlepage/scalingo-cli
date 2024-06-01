package ui

import (
	"context"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-scalingo/v7"
)

func Start(ctx context.Context) error {
	c, err := config.ScalingoAuthClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	p := tea.NewProgram(&model{
		context: ctx,
		client:  c,
	}, tea.WithAltScreen(), tea.WithContext(ctx))
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}

type model struct {
	context       context.Context
	client        *scalingo.Client
	regions       []scalingo.Region
	regionClients map[string]*scalingo.Client
}

func (m *model) Init() tea.Cmd {
	return tea.Batch(
		tea.SetWindowTitle("Scalingo"),
		tea.Sequence(
			fetchRegions(m.context),
		),
	)
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC || msg.Type == tea.KeyCtrlD {
			return m, tea.Quit
		}

	case []scalingo.Region:
		m.regions = msg
		// FIXME regionClients

	}

	return m, nil
}

func (m *model) View() string {
	s := "Regions:\n"
	for _, r := range m.regions {
		s += " - " + r.Name + " - " + r.DisplayName + "\n"
	}

	return s
}

func fetchRegions(ctx context.Context) tea.Cmd {
	return func() tea.Msg {
		regionCache, err := config.EnsureRegionsCache(ctx, config.C, config.GetRegionOpts{
			Token: os.Getenv("SCALINGO_API_TOKEN"),
		})
		if err != nil {
			// FIXME how to manage errors?
			log.Fatal("fail to read regions cache")
		}

		return regionCache.Regions
	}
}
