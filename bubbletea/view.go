package bubbletea

import (
	"fmt"

	"github.com/BayronJDv/labssl/style"
)

func (m Model) View() string {
	aid := "press esc to go back   -    press ctrl + q to quit\n"
	switch m.currentView {
	case menuView:
		return m.viewMenu() + "\n" + m.notification + "\n" + aid
	case inputView:
		return m.viewInput() + "\n" + m.notification + "\n" + aid
	case listView:
		return m.viewList() + "\n" + m.notification + "\n" + aid
	case configView:
		return m.viewConfig() + "\n" + m.notification + "\n" + aid
	case helpView:
		return m.viewHelp() + "\n" + m.notification + "\n" + aid
	default:
		return "Unknown view"
	}
}

func (m Model) viewMenu() string {
	// Usamos style.Purple en lugar de la constante local
	s := style.Purple + style.Bold + "welcome to labssl\n\n" + style.Reset
	s += "it allows you to analyze SSL/TLS configurations of domains using the SSLLabs API.\n"
	s += "Navigate using the arrow keys and press enter to select an option.\n"
	s += style.Purple + style.Bold + "Select an option:" + style.Reset + "\n\n"
	for i, item := range m.menuItems {
		if m.cursor == i {
			s += "> " + item + "\n"
		} else {
			s += "  " + item + "\n"
		}
	}
	return s + "\n"
}

func (m Model) viewInput() string {
	return fmt.Sprintf(
		"%s\n\n%s\n",
		"Enter a domain to analyze:",
		m.textInput.View(),
	)
}

func (m Model) viewList() string {
	return fmt.Sprintf(
		"%s\n\n%s\n\n%s",
		"Enter a domain to see a full report:",
		m.textInput.View(),
		"current report for the analyzed domain:\n"+
			m.stringreport+"\n",
	)
}

func (m Model) viewConfig() string {

	s := style.Purple + style.Bold + "Configuration settings\n" + style.Reset
	s += "visit help section if you need more information.\n\n"
	s += fmt.Sprintf("StartNew: %s   press n to change it\n", m.configs.startNew)
	s += fmt.Sprintf("Max Age: %d   use up/down arrows to change it\n", m.configs.maxAge)
	s += fmt.Sprintf("Is Public: %s   use right arrow to change it \n", m.configs.ispublic)
	s += fmt.Sprintf("Allow OPC: %s  use left arrow to change it\n", m.configs.allopc)
	s += "\n"
	return s
}

func (m Model) viewHelp() string {
	s := "\n"
	s += style.Purple + style.Bold + "Help information.\n" + style.Reset

	s += "here you have the description of all options:\n\n"

	s += "1. Analyze a domain: Enter a domain to analyze its SSL/TLS configuration.\n"
	s += "   depending the config it will start a new analysis or use cached results.\n"
	s += "2. View details: View the details from the analysis youve launched .\n"
	s += "   it shows the analysis in memory \n"
	s += "3. Config: Configure application settings.\n"
	s += "you can change this options :\n"
	s += "   - startsNew: if on it will force a new analysis even if there are cached results.\n"
	s += "   - Max Age: Set the maximum age for cached results (in hours).\n"
	s += "   - Is Public: off by default. if on the analysis will be published at public dashboard \n"
	s += "   - Allopc: controls how much information the endpoint returns.\n"
	s += "     'off returns a resume .\n"
	s += "     'on' returns more detailed information during the analysis process.\n"
	s += "     'done' returns all the information when the analysis is complete. (default)\n"
	s += "\n"

	return s
}
