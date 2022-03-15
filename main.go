package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"

	tea "github.com/charmbracelet/bubbletea"
)

// <------------------------------------------------------------------------------------------------------------------------>
// TODO: There's a weird thing where the first item in the list is pushed over an extra tab or space or something
// ? https://github.com/gizak/termui more tui prettyness to implement
// <------------------------------------------------------------------------------------------------------------------------>

/*

I really like these colors for TUI's :D
   softblack: #222222;
   ansiBlack: #000000;
   ansiRed: #cd3131;
   ansiGreen: #0dbc79;
   ansiYellow: #e5e510;
   ansiBlue: #2472c8;
   ansiMagenta: #bc3fbc;
   ansiCyan: #11a8cd;
   ansiWhite: #e5e5e5;
   ansiBrightBlack: #666666;
   ansiBrightRed: #f14c4c;
   ansiBrightGreen: #23d18b;
   ansiBrightYellow: #f5f543;
   ansiBrightBlue: #3b8eea;
   ansiBrightMagenta: #d670d6;
   ansiBrightCyan: #29b8db;
   ansiBrightWhite: #e5e5e5;
   activeCodeBorder: #3794ff;
   inactiveCodeBorder: #212121;
*/

// this is a bunch of settings about how the tui looks

var (
	titleStyle      = lipgloss.NewStyle().Bold(false).Foreground(lipgloss.Color("#23d18b"))
	readmeStyle     = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#23d18b"))
	foregroundStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	backgroundStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle     = foregroundStyle.Copy()
	noStyle         = lipgloss.NewStyle()
	helpStyle       = backgroundStyle.Copy()
)

type errMsg error

// model
type model struct {
	lastPress  string
	textInput  textinput.Model  // input from user
	err        error            // errors
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

		// Our shopping list is a grocery list
		choices: []string{},

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
	if m.insertMode {
		m.textInput.Focus()
	} else {
		m.textInput.Blur()
	}
	switch msg := msg.(type) {

	case errMsg:
		m.err = msg
		return m, nil

	// Is it a key press?
	case tea.KeyMsg:
		m.lastPress = msg.String()
		switch msg.String() {
		// These keys should exit the program.
		case "ctrl+c":
			return m, tea.Quit

		// jump into input mode by pushing '[' or 'i'
		case "[", "i":
			if !m.insertMode {
				m.textInput.SetValue("")
			}
			m.insertMode = true

		// escape should exit insert mode
		case "esc":
			m.insertMode = false
		}

		// Standard List Selection Mode
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

			// The "enter" key and the spacebar (a literal space) toggles
			// the selected state for the item that the cursor is pointing at.
			case "enter", " ":
				_, ok := m.selected[m.cursor]
				if ok {
					delete(m.selected, m.cursor)
				} else {
					m.selected[m.cursor] = struct{}{}
				}

			//The "d" or "delete" key will remove an item from the choices array
			case "d", tea.KeyDelete.String():
				if len(m.choices) > 0 {
					m.choices = append(m.choices[:m.cursor], m.choices[m.cursor+1:]...)
				}

			// These keys should exit the program.
			case "ctrl+c", "q":
				return m, tea.Quit

			}
		}

		// Input/Insert Mode
		// Text input mode for storage into "Dante's Data"
		if m.insertMode {
			switch msg.String() {
			case "enter":
				_, cmd := m.textInput.Update(msg)
				m.choices = append(m.choices, m.textInput.Value())
				m.insertMode = false
				return m, cmd

			// These keys should exit the program.
			case "ctrl+c":
				return m, tea.Quit

			// These keys should return us to todolist mode
			case "]":
				m.insertMode = false
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
	s := titleStyle.Render("\n<--------------------------TheIncredibleMulk's Todo List!-------------------------->\n")
	if len(m.choices) <= 0 {
		s += readmeStyle.Render("\nNo items added yet!\nPress i or [ to add an item to the list.\n")
	}

	if m.insertMode {
		s += foregroundStyle.Render(fmt.Sprintf("\nEnter what you'd like to add to the list.  %s\n", m.textInput.View()))
	} else {

	}

	s += foregroundStyle.Render("\nTodo Items\n\n")

	// Iterate over our choices
	if len(m.choices) > 0 {
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
			s += fmt.Sprintf("\n%s [%s] %s", cursor, checked, choice)
		}
	}

	// The footer
	if m.insertMode {
		s += helpStyle.Render("\nPress esc or ] return to selection mode\nPress ctrl+c to quit.\n")
	} else {
		s += helpStyle.Render("\nPress i or [ to switch to input mode\nPress ctrl+c or q to quit.\n")
	}

	// s += helpStyle.Render(fmt.Sprintf("\nm.choices: %#v:%+v", m.choices, m.choices))
	// s += helpStyle.Render(fmt.Sprintf("\nSelected: %+v", m.selected))
	s += helpStyle.Render("\nLast Key Press: " + m.lastPress)

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
