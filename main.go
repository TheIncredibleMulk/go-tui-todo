package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// <------------------------------------------------------------------------------------------------------------------------>
// TODO 'esc' is not escaping typing mode
// TODO entering typing mode carries over the entry character ex. '[' or 'i' show up in the typing box
// <------------------------------------------------------------------------------------------------------------------------>

type errMsg error

// model
type model struct {
	lastPress string
	textInput textinput.Model // input from user
	// cursorMode textinput.CursorMode // the mode our cursor is in, ex. blink, static, none
	err        error            // errors
	dante      string           // trying to fix the confusion
	insertMode bool             // insert mode chooses our focus state
	choices    []string         // items on the to-do list
	cursor     int              // which to-do list item our cursor is pointing at
	selected   map[int]struct{} // which to-do items are selected
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Midway on our life's journey, I found myself"
	ti.CharLimit = 255
	ti.Width = 20
	return model{
		// our text input
		textInput: ti,
		err:       nil,
		// Our shopping list is a grocery listm.textInput.F
		choices: []string{"Buy carrots", "Buy celery", "Buy kohlrabi"},

		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

// update
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.textInput.Blur()
	if m.insertMode {
		m.textInput.Focus()
	}
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:
		m.lastPress = msg.String()
		switch msg.String() {
		// These keys should exit the program.
		case "ctrl+c":
			return m, tea.Quit

		// jump into input mode by pushing '[' or 'i'
		case "[", "i":
			m.insertMode = true

		// escape should exit insert mode
		case "esc":
			m.insertMode = false
		}
		if !m.insertMode {
			switch msg.String() {
			// The "up" and "k" keys move the cursor up
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}

			// The "down" and "j" keys move the cursor down
			case "down", "j":
				if m.cursor < len(m.choices)-1 {
					m.cursor++
				}

			// The "enter" key and the spacebar (a literal space) toggleesc
			// the selected state for the item that the cursor is pointing at.
			case "enter", " ":
				_, ok := m.selected[m.cursor]
				if ok {
					delete(m.selected, m.cursor)
				} else {
					m.selected[m.cursor] = struct{}{}
				}
			// These keys should exit the program.
			case "ctrl+c", "q":
				return m, tea.Quit

			}
		}

		if m.insertMode {
			switch msg.String() {
			// standard function
			// switch for escaping the input

			// case "]", "escape":
			// 	return m, tea.Quit
			case "enter":
				_, cmd := m.textInput.Update(msg)
				return m, cmd

			// These keys should exit the program.
			case "ctrl+c", "q":
				return m, tea.Quit

			case "]":
				return m, tea.Quit
			}
		}
	}

	// Update the screen with the inputs from the text input box
	m.textInput, cmd = m.textInput.Update(msg)
	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, cmd
}

// view
func (m model) View() string {
	// The header
	s := "What should we buy at the market?\n\n"
	s += fmt.Sprintf("%v, %v Speak! %s\n\n", m.insertMode, m.textInput.Focused(), m.textInput.View())
	s += "Dante's Data: " + m.dante + "\n\n"
	s += "Last Key Press: " + m.lastPress + "\n\n"

	// Iterate over our choices
	for i, choice := range m.choices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := m.selected[i]; ok {
			checked = "x" // selected!
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
