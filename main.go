package main

import (
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"fmt"
)

// ========= DECLARACIONES

type state = int

const (
	MENU state = iota
	ADD
	EDIT
	DELETE
	CONFIRM
	GENERATE_FILES
)

type participant struct {
	name string
}

type model struct {
	input        textinput.Model
	state        state
	participants []string
}

func InitModel() model {
	return model{
		state:        MENU,
		// participants: []participant{},
		participants: []string{},
		input: textinput.New(),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func main() {
	textinput.New()
	p := tea.NewProgram(InitModel())
	m, err := p.Run()

	if err != nil {
		fmt.Println(m)
		panic(err)
	}
}
