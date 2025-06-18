package main

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	headerStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"})
	appStyle    = lipgloss.NewStyle().Padding(1, 2)
)

// Level of Menu
type screen int

const (
	screenMainMenu screen = iota
	screenSubMenu
)

type item string

func (i item) Title() string       { return string(i) }
func (i item) Description() string { return "" }
func (i item) FilterValue() string { return string(i) }

type listKeyMap struct {
	addLLMs key.Binding
	quit    key.Binding
	up      key.Binding
	down    key.Binding
	enter   key.Binding
}

func newListKeymap() *listKeyMap {
	return &listKeyMap{
		addLLMs: key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "add item")),
		quit:    key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
		enter:   key.NewBinding(key.WithKeys("enter", "return"), key.WithHelp("↵", "select item")),
		up:      key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up")),
		down:    key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down")),
	}
}

type model struct {
	list          list.Model
	keys          *listKeyMap
	width         int
	height        int
	pwd           string
	currentScreen screen
}

func initialModel(dir string) model {
	l := newMainMenuList()

	return model{
		list:          l,
		keys:          newListKeymap(),
		pwd:           dir,
		currentScreen: screenMainMenu,
	}
}

func newMainMenuList() list.Model {
	items := []list.Item{
		item("Start prompting and building (Coding)"),
		item("Add API keys for new agents (LLMs)"),
		item("Exit (See ya)"),
	}
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = headerStyle.Render("CozyCode")
	return l
}

func newSubMenuList2() list.Model {
	items := []list.Item{
		item("Groq LLM"),
	}
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = headerStyle.Render("Sub Menu")
	return l
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetSize(msg.Width-4, msg.Height-6)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.enter):

			switch m.currentScreen {
			case screenMainMenu:
				switch m.list.SelectedItem().(item) {
				case item("Start prompting and building (Coding)"):
					m.list = newSubMenuList2()
					m.list.SetSize(m.width-4, m.height-6)
					m.currentScreen = screenSubMenu

				case item("Add API keys for new agents (LLMs)"):
					m.list = newSubMenuList2()
					m.list.SetSize(m.width-4, m.height-6)
					m.currentScreen = screenSubMenu

				case item("Exit (See ya)"):
					return m, tea.Quit

				}

			case screenSubMenu:
				switch m.list.SelectedItem().(item) {
				case item("Back"):
					m.list = newMainMenuList()
					m.list.SetSize(m.width-4, m.height-6)
					m.currentScreen = screenMainMenu
				}

			}

		}
	}

	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return appStyle.Render(m.list.View())
}

func (m model) Init() tea.Cmd {
	return nil
}

func runTUI() error {

	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	p := tea.NewProgram(initialModel(dir), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("unable to run tui program: %w", err)
	}

	return nil
}

func main() {

	closer, err := setupLog()
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}

	if err := runTUI(); err != nil {
		_ = closer()
		log.Fatal(err)
		os.Exit(-1)
	}

	_ = closer()
}
