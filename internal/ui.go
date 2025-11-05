package runner

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// UI Configuration Constants
const (
	// Debug mode
	debugMode = false // Set to true to show viewport/scroll debug info

	// Viewport sizing
	minViewportHeight     = 5   // Minimum height for the viewport in small terminals
	maxViewportHeight     = 0   // Maximum height (0 = no limit, use full terminal)
	headerFooterLines     = 6   // Lines reserved for title (2), filter (2), help (2)
	linesPerScriptOption  = 2   // Lines each script takes (name, command+metadata)

	// Text input sizing
	filterCharLimit       = 100 // Maximum characters in filter input
	filterPromptWidth     = 3   // Width of "/ " prompt plus spacing

	// Command truncation
	commandMaxWidthBuffer = 5   // Reserve this many chars from right edge for "..."

	// Initial dimensions (will be overridden by terminal size)
	initialViewportWidth  = 80
	initialViewportHeight = 10
)

// Color palette for easy customization
type colorPalette struct {
	Black   string
	Red     string
	Green   string
	Yellow  string
	Blue    string
	Magenta string
	Cyan    string
	White   string
	Gray    string
}

var colors = colorPalette{
	Black:   "#000000",
	Red:     "#E88388",
	Green:   "#A8CC8C",
	Yellow:  "#DBAB79",
	Blue:    "#71BEF2",
	Magenta: "#D290E4",
	Cyan:    "#66C2CD",
	White:   "#FFFFFF",
	Gray:    "#B9BFCA",
}

// https://github.com/charmbracelet/lipgloss/blob/7d1b622c64d1a68cdc94b30864ae5ec3e6abc2dd/examples/ssh/main.go#L38
var (
	// Styles
	scriptNameStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(colors.Cyan))

	commandStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("250")) // Lighter gray for better readability

	metadataStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")) // Darker gray for stars/run data

	cursorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colors.Magenta))

	// Selected line backgrounds (Option 1: Dark Gray)
	selectedScriptNameBgStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#2A2A2A"))

	selectedCommandBgStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#2A2A2A"))

	promptStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(colors.Cyan))

	defaultAnswerStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color(colors.Green))
)

func FormatTimeAgo(t time.Time) string {
	duration := time.Since(t)

	switch {
	case duration < time.Minute:
		return "just now"
	case duration < time.Hour:
		mins := int(duration.Minutes())
		if mins == 1 {
			return "1 min ago"
		}
		return fmt.Sprintf("%d mins ago", mins)
	case duration < 24*time.Hour:
		hours := int(duration.Hours())
		if hours == 1 {
			return "1h ago"
		}
		return fmt.Sprintf("%dh ago", hours)
	case duration < 7*24*time.Hour:
		days := int(duration.Hours() / 24)
		if days == 1 {
			return "1d ago"
		}
		return fmt.Sprintf("%dd ago", days)
	case duration < 30*24*time.Hour:
		weeks := int(duration.Hours() / 24 / 7)
		if weeks == 1 {
			return "1w ago"
		}
		return fmt.Sprintf("%dw ago", weeks)
	default:
		months := int(duration.Hours() / 24 / 30)
		if months == 1 {
			return "1mo ago"
		}
		return fmt.Sprintf("%dmo ago", months)
	}
}

func FormatScriptOption(scored ScoredScript) string {
	return FormatScriptOptionWithWidth(scored, 0)
}

func FormatScriptOptionWithWidth(scored ScoredScript, maxWidth int) string {
	// Format: "script-name → command [★★★★☆ 24 runs, 2h ago]"
	scriptName := scriptNameStyle.Render(scored.Script.Name)

	// Prepare metadata with source indicator
	var metadata string
	var sourceIndicator string

	// Format source indicator with color
	if scored.Script.Source == "make" {
		sourceIndicator = lipgloss.NewStyle().Foreground(lipgloss.Color(colors.Green)).Render(scored.Script.Source)
	} else if scored.Script.Source != "" {
		sourceIndicator = metadataStyle.Render(scored.Script.Source)
	}

	if scored.LastUsed != nil {
		timeAgo := FormatTimeAgo(*scored.LastUsed)
		if sourceIndicator != "" {
			metadata = fmt.Sprintf("[%s %s%s]",
				sourceIndicator,
				metadataStyle.Render(fmt.Sprintf("%d runs, ", scored.UseCount)),
				metadataStyle.Render(timeAgo))
		} else {
			metadata = metadataStyle.Render(fmt.Sprintf("[%d runs, %s]", scored.UseCount, timeAgo))
		}
	} else {
		if sourceIndicator != "" {
			metadata = fmt.Sprintf("[%s %s]", sourceIndicator, metadataStyle.Render("0 runs"))
		} else {
			metadata = metadataStyle.Render("[0 runs]")
		}
	}

	// Calculate available width for command (accounting for prefix, metadata, buffer)
	commandText := scored.Script.Command
	if maxWidth > 0 {
		// Account for: "  " (2 chars) + metadata + buffer
		metadataWidth := lipgloss.Width(metadata)
		availableWidth := maxWidth - 2 - metadataWidth - commandMaxWidthBuffer

		// Truncate command if needed
		if len(commandText) > availableWidth && availableWidth > 3 {
			commandText = commandText[:availableWidth-3] + "..."
		}
	}

	command := commandStyle.Render(fmt.Sprintf(" %s", commandText))

	return fmt.Sprintf("%s\n  %s %s", scriptName, command, metadata)
}

func PromptForDefault(scored ScoredScript) (bool, error) {
	fmt.Println()
	fmt.Println(promptStyle.Render("Run the most recent script?"))
	fmt.Println()
	fmt.Println(FormatScriptOption(scored))
	fmt.Println()
	fmt.Print(promptStyle.Render("Run this script? ") + defaultAnswerStyle.Render("[Y/n]") + ": ")

	var confirmed bool = true

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("").
				Value(&confirmed).
				Affirmative("Yes").
				Negative("No"),
		),
	).WithShowHelp(false).WithShowErrors(false)

	err := form.Run()
	if err != nil {
		return false, err
	}

	return confirmed, nil
}

func ShowScriptSelection(scoredScripts []ScoredScript, initialFilter string) (*ScoredScript, error) {
	// Use the custom filterable selector for all cases now (provides dynamic sizing)
	return ShowScriptSelectionWithFilter(scoredScripts, initialFilter)
}

func PrintScriptsList(scoredScripts []ScoredScript, packageManager string) {
	fmt.Println()
	fmt.Println(promptStyle.Render("Available npm scripts (sorted by frecency):"))
	fmt.Println()

	for _, scored := range scoredScripts {
		fmt.Println(FormatScriptOption(scored))
		fmt.Printf("  %s\n", commandStyle.Render(fmt.Sprintf("Run with: %s run %s", packageManager, scored.Script.Name)))
		fmt.Println()
	}
}

// filterableSelector is a custom Bubble Tea model for script selection with editable filter
type filterableSelector struct {
	filter          textinput.Model
	viewport        viewport.Model
	allScripts      []ScoredScript
	filteredScripts []ScoredScript
	selected        int
	result          *ScoredScript
	width           int
	height          int
	quitting        bool
}

// Init initializes the filterableSelector model
func (m *filterableSelector) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles keyboard input and updates the model state
func (m *filterableSelector) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			if len(m.filteredScripts) > 0 && m.selected < len(m.filteredScripts) {
				m.result = &m.filteredScripts[m.selected]
				m.quitting = true
				return m, tea.Quit
			}

		case "esc", "ctrl+u", "alt+backspace", "ctrl+backspace":
			// Clear filter and show all scripts
			// Multiple shortcuts for different terminal/platform preferences:
			// - esc: universal
			// - ctrl+u: standard terminal "clear line"
			// - alt+backspace: may work as cmd+backspace on Mac
			// - ctrl+backspace: works on some terminals
			// - ctrl+w: standard "delete word backward"
			m.filter.SetValue("")
			m.filteredScripts = m.allScripts
			m.selected = 0
			m.viewport.GotoTop()

		case "up":
			if m.selected > 0 {
				m.selected--
			} else {
				// Wrap to bottom
				m.selected = len(m.filteredScripts) - 1
			}
			m.updateViewport()

		case "down":
			if m.selected < len(m.filteredScripts)-1 {
				m.selected++
			} else {
				// Wrap to top
				m.selected = 0
			}
			m.updateViewport()

		default:
			// Update the text input (handles typing, backspace, etc.)
			m.filter, cmd = m.filter.Update(msg)
			// Update filtered list based on new filter value
			m.filterScripts()
			// Reset selection to top when filter changes
			if m.selected >= len(m.filteredScripts) {
				m.selected = 0
			}
			m.updateViewport()
			return m, cmd
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Update text input width to match terminal (leave space for prompt)
		m.filter.Width = msg.Width - filterPromptWidth

		// Calculate viewport dimensions using configured constants
		viewportHeight := max(msg.Height-headerFooterLines, minViewportHeight)

		// Apply max height if configured (0 = no limit)
		if maxViewportHeight > 0 && viewportHeight > maxViewportHeight {
			viewportHeight = maxViewportHeight
		}

		m.viewport.Width = msg.Width
		m.viewport.Height = viewportHeight
	}

	return m, cmd
}

// updateViewport scrolls to keep the selected item visible
func (m *filterableSelector) updateViewport() {
	selectedLine := m.selected * linesPerScriptOption
	selectedLastLine := selectedLine + linesPerScriptOption - 1 // Last line of selected item

	// Get viewport bounds
	viewportTop := m.viewport.YOffset
	viewportBottom := m.viewport.YOffset + m.viewport.Height - 1 // Last visible line (inclusive)

	// Scroll up if selection starts above viewport
	if selectedLine < viewportTop {
		m.viewport.YOffset = selectedLine
	}

	// Scroll down if selection extends to or beyond viewport bottom
	if selectedLastLine >= viewportBottom {
		// Position the selected item so its last line is at the bottom of viewport
		// But ensure we don't go past the end of content
		m.viewport.YOffset = max(0, selectedLastLine-m.viewport.Height+1)
	}
}

// filterScripts updates the filtered list based on the current filter value
func (m *filterableSelector) filterScripts() {
	filterValue := strings.TrimSpace(m.filter.Value())

	if filterValue == "" {
		m.filteredScripts = m.allScripts
		return
	}

	// Use the fuzzy SearchScripts function for consistent behavior
	m.filteredScripts = SearchScripts(m.allScripts, filterValue)
}

// View renders the UI
func (m *filterableSelector) View() string {
	if m.quitting {
		return ""
	}

	var s strings.Builder

	// Title with optional debug info
	title := promptStyle.Render("📦 Search scripts (type to filter)")
	if debugMode {
		selectedLine := m.selected * linesPerScriptOption
		debugInfo := metadataStyle.Render(fmt.Sprintf(" [Term: %dx%d, VP: %dx%d (offset:%d), Sel: %d (line:%d), Scripts: %d]",
			m.width, m.height, m.viewport.Width, m.viewport.Height, m.viewport.YOffset, m.selected, selectedLine, len(m.filteredScripts)))
		title += debugInfo
	}
	s.WriteString(title + "\n\n")

	// Filter input
	s.WriteString(m.filter.View() + "\n\n")

	// Build options view
	var optionsView strings.Builder
	cursor := cursorStyle.Render("❯ ")
	blank := strings.Repeat(" ", lipgloss.Width(cursor))

	if len(m.filteredScripts) == 0 {
		optionsView.WriteString(metadataStyle.Render("No matching scripts found") + "\n")
	} else {
		for i, scored := range m.filteredScripts {
			// Add cursor for selected item
			prefix := blank
			if i == m.selected {
				prefix = cursor
			}

			// Format the option with width constraint to prevent wrapping
			formatted := FormatScriptOptionWithWidth(scored, m.width)

			// Add prefix to the first line (script name)
			lines := strings.Split(formatted, "\n")
			if len(lines) > 0 {
				scriptNameLine := lines[0]
				// Apply full-width background to selected item's script name line
				if i == m.selected {
					scriptNameLine = selectedScriptNameBgStyle.Width(m.width).Render(scriptNameLine)
				}
				optionsView.WriteString(prefix + scriptNameLine + "\n")
				// Add remaining lines with proper indentation (command + metadata)
				for _, line := range lines[1:] {
					// Apply background to selected item's command line (starting after indentation)
					if i == m.selected {
						// Find where actual content starts (after leading spaces)
						trimmed := strings.TrimLeft(line, " ")
						leadingSpaces := len(line) - len(trimmed)
						indent := line[:leadingSpaces]
						content := line[leadingSpaces:]
						// Apply background from first letter to end of line
						line = indent + selectedCommandBgStyle.Width(m.width-leadingSpaces).Render(content)
					}
					optionsView.WriteString(line + "\n")
				}
			}

			// No blank line between options for compact view
		}
	}

	// Set viewport content
	content := optionsView.String()
	m.viewport.SetContent(content)

	// Ensure viewport is scrolled to show selected item
	m.updateViewport()

	// Render viewport
	s.WriteString(m.viewport.View())

	// Optional line count debug
	if debugMode {
		actualLines := strings.Count(content, "\n")
		lineDebug := metadataStyle.Render(fmt.Sprintf("\nContent lines: %d, Expected: %d", actualLines, len(m.filteredScripts)*linesPerScriptOption))
		s.WriteString(lineDebug)
	}

	// Help text
	help := metadataStyle.Render("\n↑/↓: navigate • enter: select • esc: clear • q: quit")
	s.WriteString(help)

	return s.String()
}

// ShowScriptSelectionWithFilter shows an interactive script selector with pre-populated filter
func ShowScriptSelectionWithFilter(scoredScripts []ScoredScript, initialFilter string) (*ScoredScript, error) {
	if len(scoredScripts) == 0 {
		return nil, fmt.Errorf("no scripts available")
	}

	// Initialize text input for filter
	ti := textinput.New()
	ti.Placeholder = "Type to filter..."
	ti.Focus()
	ti.CharLimit = filterCharLimit
	ti.Width = initialViewportWidth - filterPromptWidth
	ti.Prompt = "/ "

	// Set initial filter value
	ti.SetValue(initialFilter)

	// Initialize viewport with minimal size - will be resized when window size is detected
	vp := viewport.New(initialViewportWidth, initialViewportHeight)
	vp.Style = lipgloss.NewStyle()

	// Create model (use pointer for tea.Model interface)
	model := &filterableSelector{
		filter:     ti,
		viewport:   vp,
		allScripts: scoredScripts,
		selected:   0,
		width:      0, // Will be set by WindowSizeMsg
		height:     0, // Will be set by WindowSizeMsg
	}

	// Apply initial filter
	model.filterScripts()

	// Run the program with alt screen to get full terminal dimensions
	p := tea.NewProgram(model, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		return nil, err
	}

	// Extract result
	if m, ok := finalModel.(*filterableSelector); ok {
		return m.result, nil
	}

	return nil, fmt.Errorf("no script selected")
}
