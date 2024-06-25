package ui

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	scalingo "github.com/Scalingo/go-scalingo/v7"
)

type app scalingo.App

func (a *app) Title() string {
	return a.Name
}

func (a *app) Description() string {
	return a.Region + " - " + string(a.Status)
}

func (a *app) FilterValue() string {
	return a.Name
}

type applicationsView struct {
	ctx           context.Context
	apps          list.Model
	regionClients map[string]*scalingo.Client
}

var _ view = &applicationsView{}

func (v *applicationsView) Update(msg tea.Msg) (cmd tea.Cmd) {

	switch msg := msg.(type) {

	case initViewMsg:
		return tea.Batch(
			tea.SetWindowTitle("Scalingo - Apps"),
			v.FetchApps(),
		)

	case []*scalingo.App:
		items := v.apps.Items()
		for _, sa := range msg {
			items = append(items, (*app)(sa))
		}
		slices.SortFunc(items, func(item1, item2 list.Item) int {
			return strings.Compare(item1.(*app).Name, item2.(*app).Name)
		})
		return v.apps.SetItems(items)

	case error:
		return v.apps.NewStatusMessage(msg.Error())

	}

	v.apps, cmd = v.apps.Update(msg)
	return
}

func (v *applicationsView) View() string {
	return v.apps.View()
}

func (v *applicationsView) SetSize(w int, h int) {
	v.apps.SetSize(w, h)
}

func (v *applicationsView) SetRegionClients(regionClients map[string]*scalingo.Client) {
	v.regionClients = regionClients
}

func (v *applicationsView) FetchApps() tea.Cmd {
	cmds := make([]tea.Cmd, 0, len(v.regionClients))

	for r, c := range v.regionClients {
		cmds = append(cmds, v.FetchRegionApps(r, c))
	}

	return tea.Sequence(
		v.apps.StartSpinner(),
		tea.Batch(cmds...),
		func() tea.Msg {
			v.apps.StopSpinner()
			return nil
		},
	)
}

func (v *applicationsView) FetchRegionApps(r string, c *scalingo.Client) tea.Cmd {
	return func() tea.Msg {
		apps, err := c.AppsList(v.ctx)
		if err != nil {
			return fmt.Errorf("failed to load apps for %s: %w", r, err)
		}

		return apps
	}
}

func createApplicationsView(ctx context.Context) *applicationsView {
	apps := list.New(nil, list.NewDefaultDelegate(), 0, 0)
	apps.Title = "Apps"
	apps.StatusMessageLifetime = 5 * time.Second

	return &applicationsView{
		ctx:  ctx,
		apps: apps,
	}
}
