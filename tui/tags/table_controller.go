package tui

import (
   bt "github.com/charmbracelet/bubbletea"
)


func (m *TagsTableModel) Update(msg bt.Msg) (bt.Model, bt.Cmd) {
   if len(m.Tags) < 1 { m.Init() }
   switch msg := msg.(type) {
      case bt.WindowSizeMsg:
         flexboxStyle.Width(max(100, msg.Width))
         flexboxStyle.Height(max(60, msg.Height))
         m.flex.SetWidth(msg.Width)
         m.flex.SetHeight(msg.Height)
         m.table.SetWidth(msg.Width)
         m.table.SetHeight(msg.Height)

      case bt.KeyMsg:
         switch msg.String() {
            case "up", "i", "I":
               m.table.CursorUp()

            case "down", "k", "K":
               m.table.CursorDown()

            case "left", "j", "J":
               if x, _ := m.table.GetCursorLocation(); x > 1 {
                  m.table.CursorLeft()
               }

            case "right", "l", "L":
               m.table.CursorRight()

            case "q", "Q":
               return m, bt.Quit
            default:
               return m, nil
         }
   }
   return m, nil
}

