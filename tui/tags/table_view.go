package tui


func (m *TagsTableModel) View() string {
   m.flex.ForceRecalculate()

   row := m.flex.GetRow(0)
   if row == nil{
      panic("could not find flexbox table row: 0")
   }
   cell := row.GetCell(0)
   if cell == nil {
      panic("could not find the flexbox table cell: 0")
   }

   m.table.SetWidth(cell.GetWidth())
   m.table.SetHeight(cell.GetHeight())

   cell.SetContent(m.tableStyle.Render(m.table.Render()))

   return m.windowStyle.Render(m.flex.Render())
}
