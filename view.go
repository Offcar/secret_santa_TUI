package main

import (
	tea "charm.land/bubbletea/v2"
	"fmt"
)

func (m model) View() tea.View {
	s := "  -  SECRET SANTA  -  \n"
	s += "----------------------\n\n"

	switch m.state{
	case MENU:
		s += renderMENU(&m)
	case ADD:
		s += renderADD(&m)
	}

	return tea.NewView(s)
}

func renderMENU(m *model) string {
		s := ""
		if len(m.participants) == 0 {
			s += "    NO PARTICIPANTS\n"
		} else {
			// Agregar ... arriba

			// if() mostrar ... como elemento 0
			// Mejor sería un builder, brigido igual, por ahora el modo ineficiente
			for i, p := range m.participants {
				s += fmt.Sprintf("  - (%d) %s\n", i, p)
			}
			// if () mostrar ... en elemento final m.participants + 2
		}
		s += "\n\n\n"
		s += "[A]dd participant\n\n"
		if len(m.participants) > 0 {
			s += "[E]dit participant\n"
			s += "[R]emove participant\n"
			s += "[C]reate secret_santa list\n"
		}
		return s
}

func renderADD(m *model) string {
	s := "  ADD A PARTICIPANT\n"
	s += "  " + m.input.View() + "\n\n"
	s += "[Esc] to cancel"

	return s
}
