package ui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-scalingo/v7"
)

func (m *model) FetchRegions() tea.Cmd {
	return func() tea.Msg {
		regionCache, err := config.EnsureRegionsCache(m.ctx, config.C, config.GetRegionOpts{
			Token: os.Getenv("SCALINGO_API_TOKEN"),
		})
		if err != nil {
			return fmt.Errorf("failed to read regions cache: %w", err)
		}

		return regionCache.Regions
	}
}

func (m *model) SetRegions(regions []scalingo.Region) (tea.Model, tea.Cmd) {
	m.regions = regions

	return m, func() tea.Msg {
		regionClients := make(map[string]*scalingo.Client, len(regions))

		for _, r := range regions {
			c, err := config.ScalingoClientForRegion(m.ctx, r.Name)
			if err != nil {
				return fmt.Errorf("failed to create client for region %s: %w", r.Name, err)
			}

			regionClients[r.Name] = c
		}

		return regionClients
	}
}
