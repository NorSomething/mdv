package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/viewport"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
)

const longMarkdown = ``

Welcome to your **bespoke** terminal reading experience! This layout uses a fully customized stylesheet.

## 🛠 Features Implemented:
1. **Interactive Viewport**: Use your arrow keys, Page Up/Down, or mouse wheel to scroll.
2. **Bespoke Theme**: Changed headings to a custom color, added a left border margin, and customized the blockquote style.
3. **Dynamic Reflow**: Automatically updates line-wrapping when you resize your window.

### Sample Code Blocks
` + "```go" + `
package main

import "fmt"

func main() {
    fmt.Println("Syntax highlighting matches your theme!")
}
` + "```" + `

> "The details are not the details. They make the design."
> — Charles Eames

### Scroll down to read more...
* Item A
* Item B
* Item C
* Item D
* Item E
* Item F
* Item G
* Item H
* Item I
* Item J

---
*End of Document*
`

var custom_theme = []byte(`{
	"document": {
		"margin": 2
	},
	"heading": {
		"color": "#00FFCC",
		"bold": true,
		"upper": true
	},
	"heading_2": {
		"color": "#8888FF",
		"bold": true
	},
	"blockquote": {
		"color": "#FFB86C",
		"italic": true,
		"indent": 4
	},
	"code": {
		"color": "#FF5555"
	},
	"code_block": {
		"margin": 2
	}
}`)

type model struct {
	viewport viewport.Model
	ready    bool
	width    int
	height   int
}

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

	header := fmt.Sprintf("Reader [%d%%]", int(m.viewport.ScrollPercent()*100))

	footer := fmt.Sprintf("Press 'q' to quit")

	return header + "\n" + m.viewport.View() + "\n" + footer
}

func (m model) render_markdown() (string, error) {
	r, err := glamour.NewTermRenderer(glamour.WithStylesFromJSONBytes(custom_theme), glamour.WithWordWrap(m.width-6))
	if err != nil {
		return "", err
	}

	return r.Render(longMarkdown)
}

func main() {
	p := tea.NewProgram(model{}, tea.WithAltScreen(), tea.WithMouseCellMotion())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
