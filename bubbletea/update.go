package bubbletea

import (
	"fmt"

	"github.com/BayronJDv/labssl/bubbletea/analyze"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if km, ok := msg.(tea.KeyMsg); ok {
		switch km.String() {
		case "ctrl+q":
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

	case analyze.AResponse:
		if msg.Typeofres == "fromcache" {
			m.notification = "✉️ notifications: "
			m.notification += fmt.Sprintf("Cached analysis found for %s. Grade: %s \n", msg.Report.Host, analyze.Resumegrades(msg.Report))
			m.reports[msg.Report.Host] = msg.Report
		}
		if msg.Typeofres == "newanalysis" {
			m.notification = "✉️ notifications: "
			m.notification += fmt.Sprintf("New analysis started for %s. You will be notified when it's complete. \n", msg.Report.Host)
			m.reports[msg.Report.Host] = msg.Report

		}
		if msg.Typeofres == "fromnewanalysis" {
			m.notification = "✉️ notifications: "
			m.notification += fmt.Sprintf("The new Analysis for %s is complete. Grades: %s \n", msg.Report.Host, analyze.Resumegrades(msg.Report))
			m.reports[msg.Report.Host] = msg.Report
		}
		if msg.Typeofres == "waiting for completion" {
			m.notification += "\n analysis is still in progress, wait for the notification. \n"
		}
	case analyze.ErrMsg:
		m.notification = "✉️ notifications: "
		m.notification += fmt.Sprintf("Error: %s", msg.Err)
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
				if report, ok := m.reports[domain]; ok {
					m.stringreport = analyze.Viewfullreport(report)
				} else {
					m.notification += fmt.Sprintf("\nNo analysis found for %s. Please analyze it first.\n", domain)
				}
				m.textInput.SetValue("")
				return m, cmd

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
		case "n":
			if m.configs.startNew == "on" {
				m.configs.startNew = "off"
			} else {
				m.configs.startNew = "on"
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
