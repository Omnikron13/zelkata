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

   // These functions control how an item of type T should be rendered, with optional prefix and suffix which can have
   // their own styles, and excluded from any filtering, sorting, etc. of the list.
   renderItem func(T) string
   renderPrefix func(T) string
   renderSuffix func(T) string

   // Seperator that will be rendered between items
   separator string

   // Whether the list can take focus to allow selecting items, and the index of the selected item if so.
   selectable bool
   selected int

   // Customisable styling using lipgloss
   style lg.Style
   itemStyle lg.Style
   prefixStyle lg.Style
   suffixStyle lg.Style
   selectedItemStyle lg.Style
   selectedPrefixStyle lg.Style
   selectedSuffixStyle lg.Style
   seperatorStyle lg.Style
}


// Init initializes the InlineListModel; part of the bubbletea Model interface.
func (m *InlineListModel[T]) Init() (cmd bt.Cmd) {
   m.separator = Or(m.separator, ", ")

   if m.renderItem == nil {
      m.renderItem = func (item T) string { return fmt.Sprintf("%v", item) }
   }

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
      var itemStyle, prefixStyle, suffixStyle lg.Style

      if m.selectable && i == m.selected {
         itemStyle   = m.selectedItemStyle
         prefixStyle = m.selectedPrefixStyle
         suffixStyle = m.selectedSuffixStyle
      } else {
         itemStyle   = m.itemStyle
         prefixStyle = m.prefixStyle
         suffixStyle = m.suffixStyle
      }

      if m.renderPrefix != nil {
         sb.WriteString(prefixStyle.Render(m.renderPrefix(item)))
      }

      sb.WriteString(itemStyle.Render(m.renderItem(item)))

      if m.renderSuffix != nil {
         sb.WriteString(suffixStyle.Render(m.renderSuffix(item)))
      }

      if i < len(m.items)-1 {
         sb.WriteString(m.seperatorStyle.Render(m.separator))
      }
   }

   return m.style.Render(sb.String())
}

