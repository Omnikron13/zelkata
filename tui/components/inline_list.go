package components

import (
   . "cmp"
   "fmt"
   "strings"

   bt "github.com/charmbracelet/bubbletea"
   lg "github.com/charmbracelet/lipgloss"
   "github.com/charmbracelet/bubbles/key"
)


// InlineListItemStyles is a type grouping lipgloss styles for a list item and its prefix and suffix, to facilitate
// applying a variety of different styles for different states of the list item (e.g. selected, focussed, etc.)
type InlineListItemStyles struct {
   Main lg.Style
   Prefix lg.Style
   Suffix lg.Style
}


// InlineListStyles is a type grouping lipgloss styles for a list and its items, to facilitate applying different styles
// to different states of the list (e.g. focussed, etc.)
type InlineListStyles struct {
   List lg.Style
   Item struct {
      Normal InlineListItemStyles
      Selected InlineListItemStyles
   }
   Seperator lg.Style
}


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

   // Whether the list can take focus to allow selecting items, and a pointer to the selected item (or nil if no item
   // is cyrrently selected).
   selectable bool
   selected *T

   // Whether the list has focus, which will enable or disable keybindings, possibly change styles, etc.
   Focussed bool

   // Customisable keybindings
   KeyMap struct {
      Next key.Binding
      Prev key.Binding
   }

   // Customisable styling using lipgloss
   Styles struct {
      Focussed InlineListStyles
      Unfocussed InlineListStyles
   }
}


// findIndex returns the index in the (displayed) list of the given item, or -1 if not found.
// This is intended to keep track of the selected item if/when the list is changed, e.g. by filtering, sorting, etc.
func (m *InlineListModel[T]) findIndex(item *T) int {
   for i := range m.items {
      if &m.items[i] == item {
         return i
      }
   }
   return -1
}


// GetSelected returns a pointer to the selected item, or nil if nothing is selected (yet?), the selected item is now
// gone (e.g. if the list has been filtered), or the list is not flagged as selectable in the first place.
func (m *InlineListModel[T]) GetSelected() *T {
   if !m.selectable || m.findIndex(m.selected) < 0 {
      return nil
   }
   return m.selected
}


// Init initializes the InlineListModel; part of the bubbletea Model interface.
func (m *InlineListModel[T]) Init() (cmd bt.Cmd) {
   m.separator = Or(m.separator, ", ")

   if m.renderItem == nil {
      m.renderItem = func (item T) string { return fmt.Sprintf("%v", item) }
   }

   m.KeyMap.Next = key.NewBinding(
      key.WithKeys("right", "l"),
      key.WithHelp("󰜶 /l", "focus next item"),
   )

   m.KeyMap.Prev = key.NewBinding(
      key.WithKeys("left", "h"),
      key.WithHelp("󰜳 /h", "focus previous item"),
   )

   return
}


// Update updates the InlineListModel; part of the bubbletea Model interface.
func (m *InlineListModel[T]) Update(msg bt.Msg) (model bt.Model, cmd bt.Cmd) {
   m.KeyMap.Next.SetEnabled(m.Focussed)
   m.KeyMap.Prev.SetEnabled(m.Focussed)

   switch msg := msg.(type) {
      case bt.KeyMsg:
         if m.selectable {
            i := m.findIndex(m.selected)
            switch {
               case key.Matches(msg, m.KeyMap.Next):
                  if i < len(m.items)-1 {
                     m.selected = &m.items[i+1]
                  }
               case key.Matches(msg, m.KeyMap.Prev):
                  if i > 0 {
                     m.selected = &m.items[i-1]
                  } else if i < 0 {
                     m.selected = &m.items[len(m.items)-1]
                  }
            }
         }
   }

   model = m
   return
}


// View renders the InlineListModel; part of the bubbletea Model interface.
func (m *InlineListModel[T]) View() string {
   var sb strings.Builder

   var styles InlineListStyles
   if m.Focussed {
      styles = m.Styles.Focussed
   } else {
      styles = m.Styles.Unfocussed
   }

   for i, item := range m.items {
      var itemStyles InlineListItemStyles

      if m.selectable && &m.items[i] == m.selected {
         itemStyles = styles.Item.Selected
      } else {
         itemStyles = styles.Item.Normal
      }

      if m.renderPrefix != nil {
         sb.WriteString(itemStyles.Prefix.Render(m.renderPrefix(item)))
      }

      sb.WriteString(itemStyles.Main.Render(m.renderItem(item)))

      if m.renderSuffix != nil {
         sb.WriteString(itemStyles.Suffix.Render(m.renderSuffix(item)))
      }

      if i < len(m.items)-1 {
         sb.WriteString(styles.Seperator.Render(m.separator))
      }
   }

   return styles.List.Render(sb.String())
}

