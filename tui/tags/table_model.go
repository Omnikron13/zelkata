package tui

import (
   "fmt"
   "unicode/utf8"

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
   m.selectedRow = 0
   m.selectedCol = 0

   m.widthRatio = []int{1, 4, 4, 4, 1, 4, 1}
   m.widthMin = []int{10, 20, 20, 20, 10, 20, 10}

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
         if length := utf8.RuneCountInString(t.Name) + 2; length > m.widthMin[0] {
            m.widthMin[0] = length
         }
         if length := utf8.RuneCountInString(fn); length > m.widthMin[6] {
            m.widthMin[6] = length
         }
      }
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
   m.table = stickers.NewTableSingleType[string](min(250, m.windowStyle.GetWidth()), max(20, m.windowStyle.GetHeight()), m.headers)
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

   m.table.SetRatio(m.widthRatio)
   m.table.SetMinWidth(m.widthMin)

   return nil
}

