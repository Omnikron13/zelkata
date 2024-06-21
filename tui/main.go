package tui

import (
   "os"
   "context"

   tags "github.com/omnikron13/zelkata/tui/tags"

   "github.com/urfave/cli/v3"
   bt "github.com/charmbracelet/bubbletea"
)

func MainTui(ctx context.Context, cmd *cli.Command) error {
   tm := tags.TagsTableModel{}
   tm.Init()
   p := bt.NewProgram(&tm, bt.WithAltScreen())
   if m, err := p.Run(); err != nil { bt.ErrProgramKilled.Error() } else {
      _ = m
      os.Exit(0)
   }
   return nil
}

