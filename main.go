package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/viewport"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
)

type model struct {
	viewport viewport.Model
	ready    bool
	width    int
	height   int
	content  string
}

var custom_theme = []byte(`{
	"document": {
		"margin": 2
	},
	"heading": {
		"color": "#00FFCC",
		"bold": true,
		"upper": true
	},
	"h2": {
		"color": "#8888FF",
		"bold": true
	},
	"block_quote": {
		"color": "#FFB86C",
		"italic": true,
		"indent": 4
	},
	"code": {
		"color": "#FF5555"
	},
	"code_block": {
		"margin": 2,
		"theme": "dracula"
	}
}`)

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		header_height := 1
		footer_height := 2
		vertical_margin := header_height + footer_height

		if !m.ready {
			//when init
			m.viewport = viewport.New(msg.Width, msg.Height-vertical_margin)
			m.viewport.YPosition = header_height
			m.ready = true
		} else {
			//resizign when required
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - vertical_margin
		}

		rendered, err := m.render_markdown()
		if err == nil {
			m.viewport.SetContent(rendered)
		}
	}

	//keyaction for scrolling directly to viewport
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if !m.ready {
		return "not initialized"
	}

	header := fmt.Sprintf("mdv [%d%%]", int(m.viewport.ScrollPercent()*100))

	footer := fmt.Sprintf("Press 'q' to quit")

	return header + "\n" + m.viewport.View() + "\n" + footer
}

func (m model) render_markdown() (string, error) {
	r, err := glamour.NewTermRenderer(glamour.WithStylesFromJSONBytes(custom_theme), glamour.WithWordWrap(m.width-6))
	if err != nil {
		return "", err
	}

	return r.Render(m.content)
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: reader <file.md>")
		os.Exit(1)
	}

	md, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println("something went wrong when reading the file: ", err)
		os.Exit(1)
	}

	p := tea.NewProgram(model{content: string(md)}, tea.WithAltScreen(), tea.WithMouseCellMotion())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
