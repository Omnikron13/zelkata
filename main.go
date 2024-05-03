// This main package, which builds a 'zelkata' binary, will serve as a simple CLI entry point for creatuing, curating,
// querying, and generally _using_ the Zelkata system. The urfave/cli package seems to quite nicely supply all the
// tools I need to not get bogged down in CLI work., and I can focus on the core functionality of the system itself.
package main

import (
   "context"
   "os"

   _"github.com/omnikron13/zelkata/note"
   _"github.com/omnikron13/zelkata/config"

   "github.com/urfave/cli/v3"
)


// main entry point for the 'top level' command 'zelkata'. Handling is then passed over to urfave/cli. This _may_ be
// a temporary arrangement, as I'm unsure how much overlap there is likely to be be with my ultimate plan of a modern
// TUI rather than a 'flags and subcommands' model. I may maintain both though, along with web interfaces, apps, etc.
func main() {
   cmd:= &cli.Command{
      Name:  "Zelkata",
      Usage: "add notes and stuff",

      Commands: []*cli.Command{
         {
            Name: "add",
            Aliases: []string{"a"},
            Usage: "add a note",
            Action: addCmd,
         },
      },
   }

   cmd.Run(context.Background(), os.Args)
}
