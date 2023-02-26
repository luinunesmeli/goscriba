package view

import (
	"fmt"

	"github.com/mritd/bubbles/common"
	"github.com/mritd/bubbles/selector"
)

//const (
//	listHeight   = 14
//	defaultWidth = 20
//)
//
//var (
//	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
//	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
//	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
//	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
//	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
//	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
//)
//
//type item struct {
//	title   string
//	version string
//}
//
//func (i item) FilterValue() string { return "" }
//
//type itemDelegate struct{}
//
//func (d itemDelegate) Height() int                               { return 1 }
//func (d itemDelegate) Spacing() int                              { return 0 }
//func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
//func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
//	i, ok := listItem.(item)
//	if !ok {
//		return
//	}
//
//	str := fmt.Sprintf("%s: %s", i.title, i.version)
//
//	fn := itemStyle.Render
//	if index == m.Index() {
//		fn = func(s string) string {
//			return selectedItemStyle.Render("> " + s)
//		}
//	}
//
//	fmt.Fprint(w, fn(str))
//}
//
//func newVersionList() list.Model {
//	items := []list.Item{
//		item{
//			title:   "Patch",
//			version: "v0.0.2",
//		},
//		item{
//			title:   "Minor",
//			version: "v0.1.0",
//		},
//		item{
//			title:   "Custom",
//			version: "x.x.x",
//		},
//	}
//
//	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
//	l.Title = "Choose your version type."
//	l.SetShowStatusBar(false)
//	l.SetFilteringEnabled(false)
//	l.Styles.Title = titleStyle
//	l.Styles.PaginationStyle = paginationStyle
//	l.Styles.HelpStyle = helpStyle
//
//	return l
//}

type TypeMessage struct {
	Type    string
	Version string
}

func ss() *selector.Model {
	return &selector.Model{
		Data: []interface{}{
			TypeMessage{Type: "Patch", Version: "v0.0.2"},
			TypeMessage{Type: "Minor", Version: "v0.1.0"},
			TypeMessage{Type: "Custom", Version: ""},
		},
		HeaderFunc: selector.DefaultHeaderFuncWithAppend("Select the type of release:"),
		SelectedFunc: func(m selector.Model, obj interface{}, gdIndex int) string {
			t := obj.(TypeMessage)
			return common.FontColor(fmt.Sprintf("[%d] %s (%s)", gdIndex+1, t.Type, t.Version), selector.ColorSelected)
		},
		UnSelectedFunc: func(m selector.Model, obj interface{}, gdIndex int) string {
			t := obj.(TypeMessage)
			return common.FontColor(fmt.Sprintf(" %d. %s (%s)", gdIndex+1, t.Type, t.Version), selector.ColorUnSelected)
		},
		FooterFunc: func(m selector.Model, obj interface{}, gdIndex int) string {
			return ""
		},
		FinishedFunc: func(s interface{}) string {
			//return ""
			t := s.(TypeMessage)
			msg := fmt.Sprintf("🧑‍🔬Selected version: %s(%s)\n", t.Version, t.Type)
			return common.FontColor(msg, selector.ColorFinished)
		},
	}
}
