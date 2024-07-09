package inlinelist

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
   Next key.Binding
   Prev key.Binding
}

func defaultKeyMap() KeyMap {
   return KeyMap{
      Next: key.NewBinding(
         key.WithKeys("right", "l"),
         key.WithHelp("󰜶 /l", "focus next item"),
      ),
      Prev: key.NewBinding(
         key.WithKeys("left", "h"),
         key.WithHelp("󰜳 /h", "focus previous item"),
      ),
   }
}

