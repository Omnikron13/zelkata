package components

import (
   . "cmp"
   "fmt"
   "strings"
   "unicode/utf8"

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


// inlineListCachedItem is a struct that holds _unstyled_ rendered strings for the prefix, item, and suffix of an item,
// and the _styled_ rendered string.
type inlineListCachedItem[T any] struct {
   item *T
   suffix, main, prefix, focussedNormal, unfocussedNormal, focussedSelected, unfocussedSelected string
}


// InlineListModel is a widget that displays a list of items that flow horizontally as a paragraph, joined by a
// separator, with optional prefix and/or suffix for each item. Additionally if supports selecting items with a
// 'cursor'.
// TODO: add interface requirement for generic type T
type InlineListModel[T any] struct {
   Items []T

   // itemRenderCache is a cache of the strings the render functions return, operating on the assumption that they are
   // supposed to always return the same string for a given item (i.e. pure functions).
   // TODO: add a method to explicitly rebuild the cache for situations where the render functions are not pure.
   itemRenderCache []inlineListCachedItem[T]
   itemRenderCacheChannel chan *inlineListCachedItem[T]

   // These functions control how an item of type T should be rendered, with optional prefix and suffix which can have
   // their own styles, and excluded from any filtering, sorting, etc. of the list.
   RenderItem func(T) string
   RenderPrefix func(T) string
   RenderSuffix func(T) string

   // Separator that will be rendered between items
   separator string

   // Whether the list can take focus to allow selecting items, and a pointer to the selected item (or nil if no item
   // is currently selected).
   Selectable bool
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
   for i := range m.Items {
      if &m.Items[i] == item {
         return i
      }
   }
   return -1
}


// GetSelected returns a pointer to the selected item, or nil if nothing is selected (yet?), the selected item is now
// gone (e.g. if the list has been filtered), or the list is not flagged as selectable in the first place.
func (m *InlineListModel[T]) GetSelected() *T {
   if !m.Selectable || m.findIndex(m.selected) < 0 {
      return nil
   }
   return m.selected
}


// Init initializes the InlineListModel; part of the bubbletea Model interface.
func (m *InlineListModel[T]) Init() (cmd bt.Cmd) {
   m.separator = Or(m.separator, ", ")

   if m.RenderItem == nil {
      m.RenderItem = func (item T) string { return fmt.Sprintf("%v", item) }
   }

   m.itemRenderCache = make([]inlineListCachedItem[T], 0, len(m.Items))
   m.itemRenderCacheChannel = make(chan *inlineListCachedItem[T], len(m.Items))
   go func() {
      for i := range m.Items {
         item := &m.Items[i]
         c := inlineListCachedItem[T] {
            item: item,
            main: m.RenderItem(*item),
         }
         if m.RenderPrefix != nil {
            c.prefix = m.RenderPrefix(*item)
         }
         if m.RenderSuffix != nil {
            c.suffix = m.RenderSuffix(*item)
         }
         c.unfocussedNormal = fmt.Sprintf(
            "%s%s%s",
            m.Styles.Unfocussed.Item.Normal.Prefix.Render(c.prefix),
            m.Styles.Unfocussed.Item.Normal.Main.Render(c.main),
            m.Styles.Unfocussed.Item.Normal.Suffix.Render(c.suffix),
         )
         c.focussedNormal = fmt.Sprintf(
            "%s%s%s",
            m.Styles.Focussed.Item.Normal.Prefix.Render(c.prefix),
            m.Styles.Focussed.Item.Normal.Main.Render(c.main),
            m.Styles.Focussed.Item.Normal.Suffix.Render(c.suffix),
         )
         c.unfocussedSelected = fmt.Sprintf(
            "%s%s%s",
            m.Styles.Unfocussed.Item.Selected.Prefix.Render(c.prefix),
            m.Styles.Unfocussed.Item.Selected.Main.Render(c.main),
            m.Styles.Unfocussed.Item.Selected.Suffix.Render(c.suffix),
         )
         c.focussedSelected = fmt.Sprintf(
            "%s%s%s",
            m.Styles.Focussed.Item.Selected.Prefix.Render(c.prefix),
            m.Styles.Focussed.Item.Selected.Main.Render(c.main),
            m.Styles.Focussed.Item.Selected.Suffix.Render(c.suffix),
         )
         m.itemRenderCacheChannel <- &c
      }
      close(m.itemRenderCacheChannel)
   }()

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


// itemToString converts an item of type T to a string, using RenderPrefix(), RenderItem(), and RenderSuffix()
// functions (if set), returning a styled string and the unstyled rune length.
func (m *InlineListModel[T]) itemToString(item *T, style InlineListItemStyles) (string, int) {
   var sb strings.Builder
   var n int

   if m.RenderPrefix != nil {
      s := m.RenderPrefix(*item)
      n += utf8.RuneCountInString(s)
      sb.WriteString(style.Prefix.Render(s))
   }

   s := m.RenderItem(*item)
   n += utf8.RuneCountInString(s)
   sb.WriteString(style.Main.Render(s))

   if m.RenderSuffix != nil {
      s := m.RenderSuffix(*item)
      n += utf8.RuneCountInString(s)
      sb.WriteString(style.Suffix.Render(s))
   }

   return sb.String(), n
}


// Update updates the InlineListModel; part of the bubbletea Model interface.
func (m *InlineListModel[T]) Update(msg bt.Msg) (model bt.Model, cmd bt.Cmd) {
   m.KeyMap.Next.SetEnabled(m.Focussed)
   m.KeyMap.Prev.SetEnabled(m.Focussed)

   switch msg := msg.(type) {
      case bt.KeyMsg:
         if m.Selectable {
            i := m.findIndex(m.selected)
            switch {
               case key.Matches(msg, m.KeyMap.Next):
                  if i < len(m.Items)-1 {
                     m.selected = &m.Items[i+1]
                  }
               case key.Matches(msg, m.KeyMap.Prev):
                  if i > 0 {
                     m.selected = &m.Items[i-1]
                  } else if i < 0 {
                     m.selected = &m.Items[len(m.Items)-1]
                  }
            }
         }
   }

   model = m
   return
}


// View renders the InlineListModel; part of the bubbletea Model interface.
func (m *InlineListModel[T]) View() string {
   for c := range m.itemRenderCacheChannel {
      m.itemRenderCache = append(m.itemRenderCache, *c)
   }

   var sb strings.Builder

   var styles InlineListStyles
   if m.Focussed {
      styles = m.Styles.Focussed
   } else {
      styles = m.Styles.Unfocussed
   }

   for i, c := range m.itemRenderCache {
      if m.Selectable && c.item == m.selected {
         if m.Focussed {
            sb.WriteString(c.focussedSelected)
         } else {
            sb.WriteString(c.unfocussedSelected)
         }
      } else {
         if m.Focussed {
            sb.WriteString(c.focussedNormal)
         } else {
            sb.WriteString(c.unfocussedNormal)
         }
      }

      if i < len(m.Items)-1 {
         sb.WriteString(styles.Seperator.Render(m.separator))
      }
   }

   return styles.List.Render(sb.String())
}

