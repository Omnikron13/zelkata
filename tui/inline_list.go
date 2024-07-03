package tui

import (
   . "cmp"
   "fmt"
   "strings"

   bt "github.com/charmbracelet/bubbletea"
   lg "github.com/charmbracelet/lipgloss"
)


// InlineListModel is a widget that displays a list of items that flow horizontally as a paragraph, joined by a
// separator, with optional prefix and/or suffix for each item. Additionally if suooorts selecting items with a
// 'curosr'.
// TODO: add interface requirement for generic type T
type InlineListModel[T any] struct {
   items []T

   prefix string
   suffix string
   separator string

   selectable bool
   selected int

   // Customisable styling using lipgloss
   style lg.Style
   seperatorStyle lg.Style
   itemStyle lg.Style
   prefixStyle lg.Style
   suffixStyle lg.Style
   selectedItemStyle lg.Style
   selectedPrefixStyle lg.Style
   selectedSuffixStyle lg.Style
}


// Init initializes the InlineListModel; part of the bubbletea Model interface.
func (m *InlineListModel[T]) Init() (cmd bt.Cmd) {
   m.separator = Or(m.separator, ", ")

   return
}


// Update updates the InlineListModel; part of the bubbletea Model interface.
func (m *InlineListModel[T]) Update(msg bt.Msg) (model bt.Model, cmd bt.Cmd) {
   switch msg := msg.(type) {
      case bt.KeyMsg:
         switch msg.String() {
            case "ctrl+c":
               cmd = bt.Quit
               return
         }
   }

   model = m
   return
}


// View renders the InlineListModel; part of the bubbletea Model interface.
func (m *InlineListModel[T]) View() string {
   var sb strings.Builder

   for i, item := range m.items {
      if m.selectable && i == m.selected {
         sb.WriteString(m.selectedPrefixStyle.Render(m.prefix))
         sb.WriteString(m.selectedItemStyle.Render(fmt.Sprintf("%v", item)))
         sb.WriteString(m.selectedSuffixStyle.Render(m.suffix))
      } else {
         sb.WriteString(m.prefixStyle.String())
         sb.WriteString(m.itemStyle.Render(fmt.Sprintf("%v", item)))
         sb.WriteString(m.suffixStyle.String())
      }

      if i < len(m.items)-1 {
         sb.WriteString(m.seperatorStyle.Render(m.separator))
      }
   }

   return m.style.Render(sb.String())
}

