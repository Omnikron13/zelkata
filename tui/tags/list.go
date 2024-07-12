package tui

import (
   . "cmp"
   "fmt"

   "github.com/omnikron13/zelkata/tags"

   bt "github.com/charmbracelet/bubbletea"
   lg "github.com/charmbracelet/lipgloss"
   bt_list "github.com/charmbracelet/bubbles/list"
)


// TagListItem is is a type to represent a single tag in the list view that bubbletea will render.
type TagListItem struct {
   Item *tags.Tag
}


// TODO: check best was to actually use this thing...
// FilterValue implements the Item interface
func (t TagListItem) FilterValue() string {
   return t.Item.Name
}


// Title is part of implementing the DefaultItem interface (allowing use of DefaultDelegate) and returns the tag name.
func (t TagListItem) Title() string {
   return fmt.Sprintf("%s  %s", Or(t.Item.Icon, "󰓹"), t.Item.Name)
}


// Description is part of implementing the DefaultItem interface (allowing use of DefaultDelegate) and returns the tag
// description.
func (t TagListItem) Description() string {
   return t.Item.Description
}


// TagsListModel is the model for a bubbletea component to display a list of Tag objects.
type TagsListModel struct {
   List bt_list.Model;
   Selected int
   style lg.Style
}


// Init initializes the TagsListMode; part of the bubbletea Model interface.
func (m *TagsListModel) Init() (cmd bt.Cmd) {
   m.style = lg.NewStyle()

   cmd = func() bt.Msg {
      return m.List.SelectedItem().(TagListItem).Item
   }

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
         _, v := m.style.GetFrameSize()
         m.List.SetSize(40, msg.Height-v)
   }

   cmd = func() bt.Msg {
      return m.List.SelectedItem().(TagListItem).Item
   }

   m.List, _ = m.List.Update(msg)
   model = m
   return
}


// View renders the TagsListModel as a string; finishing off the bubbletea Model interface implementation. 
func (m *TagsListModel) View() string {
   return m.style.Render(m.List.View())
}

