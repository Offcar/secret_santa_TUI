package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"unicode"

	tea "charm.land/bubbletea/v2"
	rand "math/rand"
)

type pairs struct {
	p1 string
	p2 string
}

type outputCheckedMsg struct {
	containsFiles bool
	err           error
}

type filesGeneratedMsg struct {
	err error
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case outputCheckedMsg:
		if m.quitRequested {
			return m, tea.Quit
		}
		if msg.err != nil {
			moveToFILE_ERROR(&m, msg.err)
		} else if msg.containsFiles {
			moveToOUTPUT_CONFIRM(&m)
		} else {
			moveToFILE_GENERATION(&m)
			return m, fileGenerationProcedure(&m, false)
		}
	case filesGeneratedMsg:
		if msg.err != nil {
			moveToFILE_ERROR(&m, msg.err)
			// Keep cleanup failures visible, even if exit was requested.
			m.quitRequested = false
		} else {
			m.success = true
			return m, tea.Quit
		}
	case tea.KeyPressMsg:
		key := msg.Key().String()
		if key == "ctrl+c" {
			if m.state == OUTPUT_CHECK || m.state == FILE_GENERATION {
				m.quitRequested = true
				return m, nil
			}
			return m, tea.Quit
		}

		switch m.state {
		case MENU:
			menuHandler(&m, key)
		case ADD:
			return m, addHandler(&m, key, msg)
		case EDIT_LIST:
			editListHandler(&m, key)
		case EDIT_NAME:
			return m, editNameHandler(&m, key, msg)
		case DELETE_LIST:
			deleteListHandler(&m, key)
		case DELETE_CONFIRM:
			deleteConfirmHandler(&m, key)
		case FINISH_CONFIRM:
			return m, finishConfirmHandler(&m, key)
		case OUTPUT_CONFIRM:
			return m, outputConfirmHandler(&m, key)
		case FILE_ERROR:
			fileErrorHandler(&m, key)
		}
	}

	return m, nil
}

// ================= HANDLERS =================
func menuHandler(m *model, key string) {
	if len(m.participants) == 0 {
		switch key {
		case "a", "A":
			moveToADD(m)
		}
	} else if len(m.participants) > 0 {
		switch key {
		case "a", "A":
			moveToADD(m)
		case "e", "E":
			moveToEDIT_LIST(m)
		case "r", "R":
			moveToDELETE_LIST(m)
		}

		if len(m.participants) >= 2 && (key == "f" || key == "F") {
			moveToFINISH_CONFIRM(m)
		}
	}
}

func editListHandler(m *model, key string) {
	switch key {
	case "enter":
		moveToEDIT_NAME(m)
	case "esc":
		moveToMENU(m)
	case "j", "J": // DOWN
		if m.cursor == len(m.participants)-1 {
			m.cursor = 0
		} else {
			m.cursor++
		}
	case "k", "K": // UP
		if m.cursor == 0 {
			m.cursor = len(m.participants) - 1
		} else {
			m.cursor--
		}
	}
}

func addHandler(m *model, key string, msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	switch key {
	case "enter":
		if m.input.Value() != "" {
			m.participants = append(m.participants, m.input.Value())
			moveToMENU(m)
		}
	case "esc":
		moveToMENU(m)
	default:
		if key != "space" {
			m.input, cmd = m.input.Update(msg)
		}
	}
	return cmd
}

func editNameHandler(m *model, key string, msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	switch key {
	case "enter":
		if m.input.Value() != "" {
			m.participants[m.cursor] = m.input.Value()
			moveToEDIT_LIST(m)
		}
	case "esc":
		moveToEDIT_LIST(m)
	default:
		if key != "space" {
			m.input, cmd = m.input.Update(msg)
		}
	}
	return cmd
}

func deleteListHandler(m *model, key string) {
	switch key {
	case "enter":
		moveToDELETE_CONFIRM(m)
	case "esc":
		moveToMENU(m)
	case "j", "J": // DOWN
		if m.cursor == len(m.participants)-1 {
			m.cursor = 0
		} else {
			m.cursor++
		}
	case "k", "K": // UP
		if m.cursor == 0 {
			m.cursor = len(m.participants) - 1
		} else {
			m.cursor--
		}
	}
}

func deleteConfirmHandler(m *model, key string) {
	switch key {
	case "enter":
		m.participants = slices.Delete(m.participants, m.cursor, m.cursor+1)
		moveToMENU(m)
	case "esc":
		moveToDELETE_LIST(m)
	}
}

func finishConfirmHandler(m *model, key string) tea.Cmd {
	switch key {
	case "enter":
		moveToOUTPUT_CHECK(m)
		participants := slices.Clone(m.participants)
		outputPath := m.outputPath
		return func() tea.Msg {
			if err := validateParticipants(participants); err != nil {
				return outputCheckedMsg{err: err}
			}
			entries, err := inspectOutput(outputPath)
			return outputCheckedMsg{containsFiles: len(entries) > 0, err: err}
		}
	case "esc":
		moveToMENU(m)
	}
	return nil
}

func outputConfirmHandler(m *model, key string) tea.Cmd {
	switch key {
	case "enter":
		moveToFILE_GENERATION(m)
		return fileGenerationProcedure(m, true)
	case "esc":
		moveToMENU(m)
	}
	return nil
}

func fileErrorHandler(m *model, key string) {
	if key == "esc" || key == "enter" {
		moveToMENU(m)
	}
}

func fileGenerationProcedure(m *model, overwrite bool) tea.Cmd {
	participants := slices.Clone(m.participants)
	outputPath := m.outputPath
	return func() tea.Msg {
		return filesGeneratedMsg{err: generateFiles(outputPath, participants, overwrite)}
	}
}

func validateParticipants(participants []string) error {
	if len(participants) < 2 {
		return fmt.Errorf("at least two participants are required")
	}
	seen := make(map[string]bool, len(participants))
	for _, name := range participants {
		if strings.TrimSpace(name) == "" || strings.TrimRight(name, " .") != name ||
			strings.ContainsAny(name, `/\<>:"|?*`) || strings.ContainsFunc(name, unicode.IsControl) {
			return fmt.Errorf("participant %q cannot be used as a filename", name)
		}
		key := strings.ToLower(name)
		if seen[key] {
			return fmt.Errorf("participant %q would create a duplicate filename", name)
		}
		seen[key] = true
	}
	return nil
}

// Never follow a linked output directory or recursively remove user folders.
func inspectOutput(outputPath string) ([]os.DirEntry, error) {
	if outputPath == "" {
		return nil, fmt.Errorf("output directory is not configured")
	}
	info, err := os.Lstat(outputPath)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("inspect output directory: %w", err)
	}
	if !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return nil, fmt.Errorf("output must be a directory, not a file or symbolic link")
	}
	entries, err := os.ReadDir(outputPath)
	if err != nil {
		return nil, fmt.Errorf("read output directory: %w", err)
	}
	for _, entry := range entries {
		if !entry.Type().IsRegular() {
			return nil, fmt.Errorf("output contains a subdirectory, link or special file: %s", entry.Name())
		}
	}
	return entries, nil
}

func generateFiles(outputPath string, participants []string, overwrite bool) (resultErr error) {
	if err := validateParticipants(participants); err != nil {
		return err
	}
	entries, err := inspectOutput(outputPath)
	if err != nil {
		return err
	}
	if len(entries) > 0 && !overwrite {
		return fmt.Errorf("output now contains files; confirm replacement before continuing")
	}

	players := slices.Clone(participants)
	rand.Shuffle(len(players), func(i, j int) {
		players[i], players[j] = players[j], players[i]
	})
	playerTuples := make([]pairs, 0, len(players))
	for index, player := range players {
		playerTuples = append(playerTuples, pairs{p1: player, p2: players[(index+1)%len(players)]})
	}

	tempDir, err := os.MkdirTemp(filepath.Dir(outputPath), ".secret-santa-")
	if err != nil {
		return fmt.Errorf("prepare output: %w", err)
	}
	defer func() {
		if tempDir != "" {
			if err := os.RemoveAll(tempDir); err != nil {
				resultErr = errors.Join(resultErr, fmt.Errorf("clean temporary directory %s: %w", tempDir, err))
			}
		}
	}()
	for _, pair := range playerTuples {
		filename := filepath.Join(tempDir, pair.p1+".txt")
		content := fmt.Sprintf("%s, debes dar un regalo a %s.\n", pair.p1, pair.p2)
		if err := os.WriteFile(filename, []byte(content), 0600); err != nil {
			return fmt.Errorf("write file for %s: %w", pair.p1, err)
		}
	}

	// Recheck before deleting: the folder may have changed during generation.
	entries, err = inspectOutput(outputPath)
	if err != nil {
		return err
	}
	if len(entries) > 0 && !overwrite {
		return fmt.Errorf("output now contains files; confirm replacement before continuing")
	}
	for _, entry := range entries {
		if err := os.Remove(filepath.Join(outputPath, entry.Name())); err != nil {
			return fmt.Errorf("remove previous file %s: %w", entry.Name(), err)
		}
	}
	if err := os.Remove(outputPath); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("replace output directory: %w", err)
	}
	if err := os.Rename(tempDir, outputPath); err != nil {
		return fmt.Errorf("publish output directory: %w", err)
	}
	tempDir = ""
	return nil
}

// ================= MOVE TO =================
func moveToMENU(m *model) {
	m.cursor = 0
	m.state = MENU
}

func moveToADD(m *model) {
	m.input.Reset()
	m.input.CharLimit = 50
	m.input.SetWidth(50)
	m.input.Prompt = "Write the participant's name: "
	m.input.Placeholder = "Offcar"
	m.input.Focus()
	m.state = ADD
}

func moveToEDIT_LIST(m *model) {
	m.state = EDIT_LIST
}

func moveToEDIT_NAME(m *model) {
	m.input.Reset()
	m.input.CharLimit = 50
	m.input.SetWidth(50)
	m.input.Prompt = "Write the participant's name: "
	m.input.Placeholder = m.participants[m.cursor]
	m.input.Focus()
	m.state = EDIT_NAME
}

func moveToDELETE_LIST(m *model) {
	m.state = DELETE_LIST
}

func moveToDELETE_CONFIRM(m *model) {
	m.state = DELETE_CONFIRM
}

func moveToFINISH_CONFIRM(m *model) {
	m.state = FINISH_CONFIRM
}

func moveToFILE_GENERATION(m *model) {
	m.lastError = nil
	m.state = FILE_GENERATION
}

func moveToOUTPUT_CHECK(m *model) {
	m.lastError = nil
	m.quitRequested = false
	m.state = OUTPUT_CHECK
}

func moveToOUTPUT_CONFIRM(m *model) {
	m.state = OUTPUT_CONFIRM
}

func moveToFILE_ERROR(m *model, err error) {
	m.lastError = err
	m.state = FILE_ERROR
}
