package main

import (
	tea "charm.land/bubbletea/v2"
	"fmt"
)

func (m model) View() tea.View {
	if m.success {
		return tea.NewView("")
	}
	s := "  -  SECRET SANTA  -  \n"
	s += "----------------------\n\n"

	switch m.state {
	case MENU:
		s += renderMENU(&m)
	case ADD:
		s += renderADD(&m)
	case EDIT_LIST:
		s += renderEDIT_LIST(&m)
	case EDIT_NAME:
		s += renderEDIT_NAME(&m)
	case DELETE_LIST:
		s += renderDELETE_LIST(&m)
	case DELETE_CONFIRM:
		s += renderDELETE_CONFIRM(&m)
	case FINISH_CONFIRM:
		s += renderFINISH_CONFIRM(&m)
	case OUTPUT_CHECK:
		s += renderOUTPUT_CHECK(&m)
	case OUTPUT_CONFIRM:
		s += renderOUTPUT_CONFIRM(&m)
	case FILE_GENERATION:
		s += renderFILE_GENERATION(&m)
	case FILE_ERROR:
		s += renderFILE_ERROR(&m)
	}

	return tea.NewView(s)
}

func renderOUTPUT_CHECK(m *model) string {
	if m.quitRequested {
		return "  Waiting for directory check before exiting...\n"
	}
	return "  Checking output directory...\n"
}

func renderOUTPUT_CONFIRM(m *model) string {
	s := fmt.Sprintf("  The folder %s already contains files.\n", outputDir)
	s += "  Continuing will delete ALL existing files in this folder.\n\n"
	s += fmt.Sprintf("  %s\n\n", m.outputPath)
	s += "[Enter] Delete files and continue\n"
	s += "[Esc] Cancel\n"
	return s
}

func renderFILE_GENERATION(m *model) string {
	if m.quitRequested {
		return "  Waiting for file generation and cleanup before exiting...\n"
	}
	return "  Creating Secret Santa files...\n"
}

func renderFILE_ERROR(m *model) string {
	s := "  COULD NOT CREATE SECRET SANTA FILES\n\n"
	s += fmt.Sprintf("  %v\n\n", m.lastError)
	s += "[Esc] Back to participants\n"
	return s
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
		s += "[R]emove participant\n\n"

		if len(m.participants) >= 2 {
			s += "[F]inish and create SECRET_SANTA files\n"
		}
	}
	return s
}

func renderADD(m *model) string {
	s := "  ADD A PARTICIPANT\n"
	s += "  " + m.input.View() + "\n\n"
	s += "[Esc] to cancel"

	return s
}

func renderEDIT_LIST(m *model) string {
	s := "  CHOOSE A PARTICIPANT NAME TO EDIT:\n"
	for i, n := range m.participants {
		if m.cursor == i {
			s += fmt.Sprintf("  - > (%d) %s \n", i, n)
		} else {
			s += fmt.Sprintf("  - (%d) %s \n", i, n)
		}
	}
	s += "\n\n[Esc] to go back"
	return s
}

func renderDELETE_LIST(m *model) string {
	s := "  CHOOSE A PARTICIPANT NAME TO DELETE:\n"
	for i, n := range m.participants {
		if m.cursor == i {
			s += fmt.Sprintf("  - > (%d) %s \n", i, n)
		} else {
			s += fmt.Sprintf("  - (%d) %s \n", i, n)
		}
	}
	s += "\n\n[Enter] TO DELETE"
	s += "\n\n[Esc] to go cancel"
	return s
}

func renderDELETE_CONFIRM(m *model) string {
	s := "    ¿Are you sure you want to remove this participant?\n"
	s += fmt.Sprintf("          (%d) - %s \n\n", m.cursor, m.participants[m.cursor])
	s += "\n\n[Enter] TO CONFIRM DELETION\n"
	s += "[Esc] to go cancel"
	return s
}

func renderEDIT_NAME(m *model) string {
	s := "  EDIT PARTICIPANT NAME\n"
	s += "  " + m.input.View() + "\n\n"
	s += "[Esc] to cancel"
	return s
}

func renderFINISH_CONFIRM(m *model) string {
	s := " ¿CREATE SECRET SANTA FILES WITH THIS PARTICIPANT LIST?\n"

	for i, n := range m.participants {
		s += fmt.Sprintf("    (%d) - %s\n", i, n)
	}

	s += "\n\n"
	s += "[Enter] to confirm\n"
	s += "[Esc] to cancel\n"

	return s
}
