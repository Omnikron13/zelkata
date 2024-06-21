package tui

import (
   "fmt"
   "strings"
   "unicode/utf8"

   "github.com/omnikron13/zelkata/tags"
   "k8s.io/apimachinery/pkg/util/sets"

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
   HashMap map[string]*tags.Tag;
   flex *stickers.FlexBox;
   headers []string;
   widthRatio []int;
   widthMin []int;
   table *stickers.TableSingleType[string];
   selectedRow uint;
   selectedCol uint;
   windowStyle lipgloss.Style;
}


func (m *TagsTableModel) Init() bt.Cmd {
   m.headers = []string{
      "󱤇 ",
      "Name",
      "Description",
      "Aliases",
      "Parents",
      "󰭷 Note Count",
      "Note IDs",
      "Filename",
   }

   m.selectedRow = 0
   m.selectedCol = 1

   m.widthRatio = []int{1, 10, 20, 20, 20, 5, 100, 10}
   m.widthMin = make([]int, len(m.widthRatio))

   for i, s := range m.headers {
      m.widthMin[i] = utf8.RuneCountInString(s)
   }

   m.windowStyle = lipgloss.NewStyle().
      Padding(1).
      Margin(1).
      Border(lipgloss.RoundedBorder())
      //Background(lipgloss.Color("#24273A"))

   var err error
   if m.Tags, err = tags.LoadAll(); err != nil {
      panic(fmt.Errorf("Error loading TagMap for tabs table: %v\n", err))
   }
   m.HashMap = make(map[string]*tags.Tag)
   for _, t := range m.Tags {
      if fn, err := t.GenFileName(); err != nil {
         panic(err)
      } else {
         m.HashMap[fn] = t
         if length := utf8.RuneCountInString(t.Name); length > m.widthMin[1] {
            m.widthMin[1] = length
         }
         if length := utf8.RuneCountInString(fn); length > m.widthMin[7] {
            m.widthMin[7] = length
         }
      }
   }

   m.table = stickers.NewTableSingleType[string](min(250, m.windowStyle.GetWidth()), max(20, m.windowStyle.GetHeight()), m.headers)
   m.flex = stickers.NewFlexBox(1, 1)

   rows := make([][]string, 0, 16)
   for _, t := range m.Tags {
      filename, err := t.GenFileName()
      if err != nil { panic(err) }

      aliasesStr := strings.Join(t.Aliases, ", ")
      parentsStr := strings.Join(sets.List(t.Parents), ", ")
      notesStr := strings.Join(sets.List(t.Notes), ", ")

      r := make([]string, 0, 3)
      r = append(r,
         fmt.Sprintf("%s ", t.Icon),
         t.Name,
         t.Description,
         aliasesStr,
         parentsStr,
         fmt.Sprintf("%d", len(t.Notes)),
         notesStr,
         filename,
      )
      rows = append(rows, r)
   }
   m.table.AddRows(rows)

   m.table.SetRatio(m.widthRatio)
   m.table.SetMinWidth(m.widthMin)

   // Move the table cursor from the icon to name column by default
   m.table.CursorRight()

   return nil
}

