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
               n := max(m.selectedRow - 1, 0)
               if n != m.selectedRow {
                  m.selectedRow = n
                  m.table.CursorUp()
               }

            case "down", "k", "K":
               n := min(m.selectedRow + 1, uint(len(m.HashMap) - 1))
               if n != m.selectedRow {
                  m.selectedRow = n - 1
                  m.table.CursorDown()
               }

            case "left", "j", "J":
               n := max(m.selectedCol - 1, 1)
               if n != m.selectedCol{
                  m.selectedCol = n
                  m.table.CursorLeft()
               }

            case "right", "l", "L":
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

