package tui


func (m *TagsTableModel) View() string {
   //m.flex.ForceRecalculate()
   //r := m.flex.NewRow()
   //c := stickers.NewFlexBoxCell(1, 1)
   //c.SetContent(m.table.Render())
   //r.AddCells([]*stickers.FlexBoxCell{c})
   //m.flex.AddRows([]*stickers.FlexBoxRow{r})
   //return m.windowStyle.Render(m.table.Render())
   //return flexboxStyle.Render(m.flex.Render())
   return m.table.Render()
}
