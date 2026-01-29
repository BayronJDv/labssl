package bubbletea

import (
	"github.com/BayronJDv/labssl/bubbletea/analyze"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Constantes para identificar las vistas
const (
	menuView   = 0
	inputView  = 1
	listView   = 2
	configView = 3
	helpView   = 4
)

// configs mantiene las configuraciones de la aplicación
type configs struct {
	maxAge   int
	ispublic string
	allopc   string
	startNew string
}

// Model mantiene el estado de toda la aplicación
type Model struct {
	currentView  int
	menuItems    []string
	cursor       int
	textInput    textinput.Model
	notification string
	configs      configs
	reports      map[string]analyze.SSLLabsResponse
	stringreport string
}

// InitialModel crea el estado inicial
func InitialModel() Model {
	ti := textinput.New()
	ti.Placeholder = "example.com"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return Model{
		currentView:  menuView,
		menuItems:    []string{"Analyze a domain", "View details", "config", "help"},
		textInput:    ti,
		notification: "✉️ notifications: ",
		configs: configs{
			startNew: "on",
			maxAge:   24,
			ispublic: "off",
			allopc:   "done",
		},
		reports: make(map[string]analyze.SSLLabsResponse),
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}
