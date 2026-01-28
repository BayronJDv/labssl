package bubbletea

import (
	"fmt"

	"github.com/BayronJDv/labssl/analyze"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if km, ok := msg.(tea.KeyMsg); ok {
		switch km.String() {
		case "q":
			return m, tea.Quit
		}
	}
	switch m.currentView {
	case menuView:
		return m.updateMenuView(msg)
	case inputView:
		return m.updateInputView(msg)
	case listView:
		return m.updateListView(msg)
	case configView:
		return m.updateConfigView(msg)
	case helpView:
		return m.updateHelpView(msg)
	default:
		return m, cmd
	}

}

func (m Model) updateMenuView(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.menuItems)-1 {
				m.cursor++
			}
		case "q", "esc":
			return m, tea.Quit
		case "enter", " ":
			switch m.cursor {
			case 0:
				m.currentView = inputView
			case 1:
				m.currentView = listView
			case 2:
				m.currentView = configView
			case 3:
				m.currentView = helpView
			}
		}

	case analyze.StatusMsg:
		if msg == 200 {
			m.notification += "\n an analysis has been started, wait for the notification, then go to 'View Results' to see it."
		} else {
			m.notification += fmt.Sprintf("Received unexpected status code: %d", msg)
		}
	case analyze.SuccessMsg:
		m.isloading = false
		m.notification = "✉️ notifications: "
		m.notification += fmt.Sprintf("%s \n check full results at view from cache ", string(msg))

	}
	return m, cmd

}

func (m Model) updateInputView(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			domain := m.textInput.Value()
			if domain != "" {
				m.isloading = true
				m.notification += fmt.Sprintf("Started analysis for %s", domain)
				cmd2 := analyze.CheckSomeUrl(m.configs.maxAge, domain, m.configs.ispublic, m.configs.startNew, m.configs.allopc)
				m.textInput.SetValue("")
				m.currentView = menuView
				return m, tea.Batch(cmd, cmd2)
			} else {
				m.notification += "Please enter a valid domain."
			}
		case "esc":
			m.currentView = menuView
			m.notification = "✉️ notifications: "
		}
	}
	return m, cmd
}

func (m Model) updateListView(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			domain := m.textInput.Value()
			if domain != "" {
				m.isloading = true
				m.notification += fmt.Sprintf("lokking for cached analysis for %s", domain)
				cmd2 := analyze.Checkfromcache(domain)
				m.textInput.SetValue("")
				m.currentView = menuView
				return m, tea.Batch(cmd, cmd2)
			} else {
				m.notification += "Please enter a valid domain."
			}
		case "esc":
			m.currentView = menuView
			m.notification = "✉️ notifications: "
		}
	}
	return m, cmd
}

func (m Model) updateConfigView(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok && msg.String() == "esc" {
		m.currentView = menuView
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			m.configs.maxAge++
		case "down":
			if m.configs.maxAge > 0 {
				m.configs.maxAge--
			}
		case "right":
			if m.configs.ispublic == "on" {
				m.configs.ispublic = "off"
			} else {
				m.configs.ispublic = "on"
			}
		case "left":
			switch m.configs.allopc {
			case "on":
				m.configs.allopc = "off"
			case "off":
				m.configs.allopc = "done"
			default:
				m.configs.allopc = "on"
			}
		}
	}

	return m, nil
}

func (m Model) updateHelpView(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok && msg.String() == "esc" {
		m.currentView = menuView
	}
	return m, nil
}
