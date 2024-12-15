package cmd

import (
	"context"

	"github.com/bootengine/boot/internal/helper"
	"github.com/bootengine/boot/internal/usecase"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all installed modules.",
	Long:  `list all installed modules, displaying their name, type, and location in the filesystem.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return helper.WithModuleUsecase(func(ctx context.Context, use *usecase.ModuleUsecase) error {
			list, err := use.ListModules(ctx)
			if err != nil {
				return err
			}
			columns := []table.Column{
				{Title: "name", Width: 10},
				{Title: "location", Width: 60},
				{Title: "type", Width: 10},
			}

			rows := func() []table.Row {
				res := make([]table.Row, len(list))
				for i, module := range list {
					res[i] = table.Row{module.Name, module.Path, string(module.Type)}
				}
				return res
			}()

			t := table.New(
				table.WithColumns(columns),
				table.WithRows(rows),
			)

			s := table.DefaultStyles()
			s.Header = s.Header.
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("240")).
				BorderBottom(true).
				Bold(false)
			s.Selected = s.Selected.
				Foreground(lipgloss.Color("229")).
				Background(lipgloss.Color("57")).
				Bold(false)
			t.SetStyles(s)

			m := listTableModel{t}
			if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
				log.Fatal("failed to display list")
			}
			return nil
		})
	},
}

func init() {
	moduleCmd.AddCommand(listCmd)
}

type listTableModel struct {
	table table.Model
}

func (m listTableModel) Init() tea.Cmd {
	return nil
}

func (m listTableModel) View() string {
	return m.table.View()
}

func (m listTableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		nameCol := table.Column{
			Title: "name",
			Width: msg.Width / 10,
		}

		locationCol := table.Column{
			Title: "location",
			Width: msg.Width - (2 * msg.Width / 10),
		}

		typeCol := table.Column{
			Title: "type",
			Width: msg.Width / 10,
		}

		m.table.SetColumns([]table.Column{nameCol, locationCol, typeCol})
		return m, cmd
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}
