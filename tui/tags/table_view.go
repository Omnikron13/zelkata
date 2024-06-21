package tui


func (m *TagsTableModel) View() string {
   m.flex.ForceRecalculate()

   row, exists := m.flex.GetRow(0)
   if !exists {
      panic("could not find flexbox table row: 0")
   }
   cell := row.Cell(0)
   if cell == nil {
      panic("could not find the flexbox table cell: 0")
   }

   m.table.SetWidth(cell.GetWidth())
   m.table.SetHeight(cell.GetHeight())

   cell.SetContent(m.table.Render())

   return flexboxStyle.Render(m.flex.Render())
}
