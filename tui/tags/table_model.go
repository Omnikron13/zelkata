package tui

import (
   "fmt"
   //"strings"
   "unicode/utf8"

   "github.com/omnikron13/zelkata/tags"
   //"k8s.io/apimachinery/pkg/util/sets"

   bt "github.com/charmbracelet/bubbletea"
   "github.com/charmbracelet/lipgloss"

   "github.com/omnikron13/stickers"
)


type TagsTableModel struct {
   Tags tags.TagMap;
   HashMap map[string]*tags.Tag;
   flex *stickers.FlexBox;
   headers []string;
   widthRatio []int;
   widthMin []int;
   table *stickers.TableSingleType[string];
   windowStyle lipgloss.Style;
   tableStyle lipgloss.Style;
}


func (m *TagsTableModel) Init() bt.Cmd {
   m.headers = []string{
      "󱤇 ",
      "Name",
      "Description",
      //"Aliases",
      //"Parents",
      "󰭷 Notes",
      //"Note IDs",
      //"Filename",
   }

   m.widthRatio = []int{
      1,
      20,
      200,
      //20,
      //20,
      1,
      //100,
      //10,
   }
   m.widthMin = make([]int, len(m.widthRatio))

   for i, s := range m.headers {
      m.widthMin[i] = utf8.RuneCountInString(s) + 1
   }

   m.windowStyle = lipgloss.NewStyle().
      Background(lipgloss.Color("#24273A"))

   m.tableStyle = lipgloss.NewStyle().
      Background(lipgloss.Color("#24273A"))

   var err error
   if m.Tags, err = tags.LoadAll(); err != nil {
      panic(fmt.Errorf("Error loading TagMap for tabs table: %v\n", err))
   }
   m.HashMap = make(map[string]*tags.Tag)
   for _, t := range m.Tags {
      m.HashMap[t.Name] = t
      if length := utf8.RuneCountInString(t.Name) + 1; length > m.widthMin[1] {
         m.widthMin[1] = length
      }
      //fn, err := t.GenFileName()
      //if err != nil { panic(err) }
      //if length := utf8.RuneCountInString(fn); length > m.widthMin[7] {
      //   m.widthMin[7] = length
      //}
   }

   m.table = stickers.NewTableSingleType[string](120, 20, m.headers)

   rows := make([][]string, 0, 16)
   for _, t := range m.HashMap {
      //filename, err := t.GenFileName()
      //if err != nil { panic(err) }

      //aliasesStr := strings.Join(t.Aliases, ", ")
      //parentsStr := strings.Join(sets.List(t.Parents), ", ")
      //notesStr := strings.Join(sets.List(t.Notes), ", ")

      r := make([]string, 0, 3)
      r = append(r,
         fmt.Sprintf("%s ", t.Icon),
         t.Name,
         t.Description,
         //aliasesStr,
         //parentsStr,
         fmt.Sprintf("%d", len(t.Notes)),
         //notesStr,
         //filename,
      )
      rows = append(rows, r)
   }
   m.table.AddRows(rows)

   m.table.SetRatio(m.widthRatio)
   m.table.SetMinWidth(m.widthMin)

   m.table.OrderByColumn(1)

   // Move the table cursor from the icon to name column by default
   m.table.CursorRight()

   m.flex = stickers.NewFlexBox(1, 1)
   row := m.flex.NewRow()
   cell := stickers.NewFlexBoxCell(1, 1)
   row.AddCells([]*stickers.FlexBoxCell{cell})
   m.flex.AddRows([]*stickers.FlexBoxRow{row})

   return nil
}

