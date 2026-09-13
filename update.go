package main

import (
	tea "charm.land/bubbletea/v2"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		key := msg.Key().String()
		if key == "ctrl+c" {
			return m, tea.Quit
		}

		switch m.state {
		case MENU:
			menuHandler(&m, key)
		case ADD:
			addHandler(&m, key, msg)
		}
	}

	return m, nil
}

func menuHandler(m *model, key string) {
	if len(m.participants) > 0 {
		switch key {
		case "a", "A":
			moveToAdd(m)
		}
	} else {
		switch key {
		case "a", "A":
			moveToAdd(m)
		case "d", "D":
			// moveToEdit(m)
		}
	}
}

func addHandler(m *model, key string, msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	switch key {
	case "enter":
		if m.input.Value() != "" {
			m.participants = append(m.participants, m.input.Value())
			moveToMenu(m)
		}
	case "esc":
		moveToMenu(m)
	default:
		if key != "space" {
			m.input, cmd = m.input.Update(msg)
		}
	}
	return cmd
}

func moveToMenu(m *model) {
	m.input.Reset()
	m.state = MENU
}

func moveToAdd(m *model) {
	m.input.Reset()
	m.input.Prompt = "Write the participant's name: "
	m.input.Placeholder = "Andrew ;)"
	m.input.CharLimit = 50
	m.input.Focus()
	m.state = ADD
}

func moveToEdit()   {}
func moveToDelete() {}
