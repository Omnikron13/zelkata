package tui

import (
   . "cmp"
   "fmt"
   "sort"

   "github.com/omnikron13/zelkata/tags"

   bt "github.com/charmbracelet/bubbletea"
   lg "github.com/charmbracelet/lipgloss"
   bt_list "github.com/charmbracelet/bubbles/list"
)


// tagListItem is is a type to represent a single tag in the list view that bubbletea will render.
type tagListItem struct {
   Tag *tags.Tag
}


// TODO: check best was to actually use this thing...
// FilterValue implements the Item interface
func (t tagListItem) FilterValue() string {
   return t.Tag.Name
}


// Title is part of implementing the DefaultItem interface (allowing use of DefaultDelegate) and returns the tag name.
func (t tagListItem) Title() string {
   return fmt.Sprintf("%s  %s", Or(t.Tag.Icon, "󰓹"), t.Tag.Name)
}


// Description is part of implementing the DefaultItem interface (allowing use of DefaultDelegate) and returns the tag
// description.
func (t tagListItem) Description() string {
   return t.Tag.Description
}


// TagsListModel is the model for a bubbletea component to display a list of Tag objects.
type TagsListModel struct {
   list bt_list.Model;
   Selected int
   style lg.Style
}


// Init initializes the TagsListMode; part of the bubbletea Model interface.
func (m *TagsListModel) Init() (cmd bt.Cmd) {
   tagmap, err := tags.LoadAll()
   if err != nil {
      return bt.Quit
   }

   items := []bt_list.Item{}

   for n, t := range tagmap {
      if n != t.NormalisedName() {
         continue
      }
      items = append(items, tagListItem{Tag: t})
   }

   sort.Slice(items, func(i, j int) bool {
      return items[i].(tagListItem).Tag.Name < items[j].(tagListItem).Tag.Name
   })

   m.list = bt_list.New(items, bt_list.NewDefaultDelegate(), 1, 1)

   m.style = lg.NewStyle()

   return
}


// Update responds to messages (such as user input) and updates the model as necessary in response, returning the
// updated model and any (further) commands to be executed.
func (m *TagsListModel) Update(msg bt.Msg) (model bt.Model, cmd bt.Cmd) {
   switch msg := msg.(type) {
      case bt.KeyMsg:
         switch msg.String() {
            case "ctrl+c":
               cmd = bt.Quit
               return
         }

      case bt.WindowSizeMsg:
         //h, v := m.style.GetFrameSize()
         //m.list.SetSize(msg.Width-h, msg.Height-v)
         _, v := m.style.GetFrameSize()
         m.list.SetSize(40, msg.Height-v)
   }

   m.list, cmd = m.list.Update(msg)
   model = m
   return
}


// View renders the TagsListModel as a string; finishing off the bubbletea Model interface implementation. 
func (m *TagsListModel) View() string {
   return m.style.Render(m.list.View())
}

