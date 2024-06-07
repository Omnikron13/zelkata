package tui

import (
   "fmt"

   "github.com/omnikron13/zelkata/tags"

   bt "github.com/charmbracelet/bubbletea"
   "github.com/charmbracelet/lipgloss"

   "github.com/76creates/stickers"
)


var flexboxStyle = lipgloss.NewStyle()

func init() {
   flexboxStyle.Padding(1, 1, 1, 1)
}


type TagsTableModel struct {
   Tags tags.TagMap;
   HashMap map[string]tags.Tag;
   flex *stickers.FlexBox;
   headers []string;
   table *stickers.TableSingleType[string];
   selectedRow uint;
   selectedCol uint;
}


func (m *TagsTableModel) Init() bt.Cmd {
   m.selectedRow = 0
   m.selectedCol = 0
   m.Tags, _ = tags.LoadAll();
   m.HashMap = make(map[string]tags.Tag)
   for _, t := range m.HashMap {
      if fn, err := t.GenFileName(); err != nil {
         panic(err)
      } else { m.HashMap[fn] = t }
   }
   m.headers = []string{
      "󱤇 Name",
      "Description",
      "Aliases",
      "Parents",
      "󰭷 Note Count",
      "󰊕 Hashes",
      "Filename",
   }
   m.table = stickers.NewTableSingleType[string](240, min(60, len(m.Tags)), m.headers)
   m.flex = stickers.NewFlexBox(1, 1)

   rows := make([][]string, 0, 16)
   for _, t := range m.Tags {
      filename, err := t.GenFileName()
      if err != nil { panic(err) }

      r := make([]string, 0, 3)
      r = append(r,
         fmt.Sprintf("%s %s", t.Icon, t.Name),
         t.Description,
         fmt.Sprintf("%v", t.Aliases),
         fmt.Sprintf("%v", t.Parents),
         fmt.Sprintf("%d", len(t.Notes)),
         fmt.Sprintf("%v", t.Notes),
         filename,
      )
      rows = append(rows, r)
   }
   m.table.AddRows(rows)
   return nil
}


func (m *TagsTableModel) Update(msg bt.Msg) (bt.Model, bt.Cmd) {
   if len(m.Tags) < 1 { m.Init() }
   switch msg := msg.(type) {
      case bt.WindowSizeMsg:
         flexboxStyle.Width(max(100, msg.Width))
         flexboxStyle.Height(max(60, msg.Height))

      case bt.KeyMsg:
         switch msg.String() {
            case fmt.Sprintf("%s", bt.KeyUp), "i", "I":
               n := max(m.selectedRow - 1, 0)
               if n != m.selectedRow {
                  m.selectedRow = n
                  m.table.CursorUp()
               }

            case fmt.Sprintf("%s", bt.KeyDown), "k", "K":
               n := min(m.selectedRow + 1, uint(len(m.Tags) - 1))
               if n != m.selectedRow {
                  m.selectedRow = n
                  m.table.CursorDown()
               }

            case fmt.Sprintf("%s", bt.KeyLeft), "j", "J":
               n := max(m.selectedCol- 1, 0)
               if n != m.selectedCol{
                  m.selectedCol = n
                  m.table.CursorLeft()
               }

            case fmt.Sprintf("%s", bt.KeyRight), "l", "L":
               n := min(m.selectedCol + 1, uint(len(m.headers) - 1))
               if n != m.selectedCol {
                  m.selectedCol = n
                  m.table.CursorRight()
               }

            case "q", "Q":
               return m, bt.Quit
            default:
               return m, nil
         }
   }
   return m, nil
}


func (m *TagsTableModel) View() string {
   //m.flex.ForceRecalculate()
   //r := m.flex.NewRow()
   //c := stickers.NewFlexBoxCell(1, 1)
   //c.SetContent(m.table.Render())
   //r.AddCells([]*stickers.FlexBoxCell{c})
   //m.flex.AddRows([]*stickers.FlexBoxRow{r})
   return m.table.Render()
   //return flexboxStyle.Render(m.flex.Render())
}

