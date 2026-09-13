package main

import (
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"fmt"
	"os"
	"path/filepath"
)

const outputDir = "Output"

type state = int

const (
	MENU state = iota
	ADD
	EDIT_NAME
	EDIT_LIST
	DELETE_LIST
	DELETE_CONFIRM
	FINISH_CONFIRM
	OUTPUT_CHECK
	OUTPUT_CONFIRM
	FILE_GENERATION
	FILE_ERROR
)

type participant struct {
	name string
}

type model struct {
	input         textinput.Model
	cursor        int
	state         state
	participants  []string
	outputPath    string
	lastError     error
	success       bool
	quitRequested bool
}

func InitModel() model {
	m := model{
		state: MENU,
		input: textinput.New(),
	}
	if outputDir == "" || outputDir == "." || outputDir == ".." || filepath.IsAbs(outputDir) || filepath.Base(outputDir) != outputDir {
		m.lastError = fmt.Errorf("output directory must be a folder name")
	} else {
		m.outputPath, m.lastError = filepath.Abs(outputDir)
	}
	if m.lastError != nil {
		moveToFILE_ERROR(&m, m.lastError)
	}
	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

func main() {
	p := tea.NewProgram(InitModel())
	m, err := p.Run()

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if finalModel, ok := m.(model); ok && finalModel.success {
		fmt.Println("Files created successfully, happy Secret Santa!")
	}
}
