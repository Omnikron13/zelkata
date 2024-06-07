package tui

import (
   _"fmt"
   "os"
   "context"

   "github.com/urfave/cli/v3"
   bt "github.com/charmbracelet/bubbletea"
)

func MainTui(ctx context.Context, cmd *cli.Command) error {
   tm := TagsTableModel{}
   tm.Init()
   p := bt.NewProgram(&tm, bt.WithAltScreen())
   if m, err := p.Run(); err != nil { bt.ErrProgramKilled.Error() } else {
      _ = m
      os.Exit(0)
   }
   return nil
}

