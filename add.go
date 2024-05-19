package main

import (
   "context"
   "fmt"
   "os"
   "os/exec"
   "path/filepath"
   "strings"

   "github.com/omnikron13/zelkata/note"
   "github.com/omnikron13/zelkata/paths"
   "github.com/omnikron13/zelkata/tag"

   "github.com/charmbracelet/bubbles/textinput"
   tea "github.com/charmbracelet/bubbletea"
   "github.com/urfave/cli/v3"
)


func addCmd(ctx context.Context, cmd *cli.Command) error {
   // TODO: default to a bubbletea(bubbles) TextArea, with a hotkey to launch a full editor?
   // This sets up launching an external editor to write the note body, which is temporarily stored in a state file,
   // which potentially also acts as a draft file if the user saves while editing but the add process is interrupted.
   newNoteFile := filepath.Join(paths.State(), "new-note.md")
   editCmd := exec.Command(os.Getenv("EDITOR"), newNoteFile)
   editCmd.Stdin  = os.Stdin
   editCmd.Stdout = os.Stdout
   editCmd.Stderr = os.Stderr

   // Start the editor and wait for it to finish
   if err := editCmd.Start(); err != nil {
      return err
   }
   editCmd.Wait()

   // Read the draft note file into a string and clear it so the next add has an empty file buffer
   s, err := os.ReadFile(newNoteFile)
   if err != nil {
      return err
   }
   os.Remove(newNoteFile)

   // Create a new Note, including initialising the Meta struct (so this is when the UUID & timestamp are generated)
   note := note.New(string(s))

   // Spin up bubbletea (crudely, for now) to get the tags for the note
   var m tea.Model
   if m, err = tea.NewProgram(initialAddCmdModel()).Run(); err != nil {
      return err
   }
   acm := m.(addCmdModel)

   // Actually add the tags to the Note
   note.Tags = acm.tags

   // Save the note to the configured notes dir with the configured filename
   if err := note.Save(); err != nil {
      return err
   }

   // Update the specified tags with the new note ID
   for _, t := range note.Tags {
      if err := tag.Add(t, note.ID); err != nil {
         return err
      }
   }

   // Return nil if everything went well
   return nil
}


// BubbleTea code for adding the tags to a Note
type addCmdModel struct {
   tags []string
   input textinput.Model
   err error
}

func (m addCmdModel) Init() tea.Cmd {
   // Just return `nil`, which means "no I/O right now, please."
   //return nil
   return textinput.Blink
}

func (m addCmdModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
   var cmd tea.Cmd

   switch msg := msg.(type) {
      case tea.KeyMsg:
         switch msg.Type {
            case tea.KeyCtrlC, tea.KeyEsc:
               return m, tea.Quit

            case tea.KeyEnter:
               s := strings.TrimSpace(m.input.Value())
               if s == "" {
                  return m, tea.Quit
               }
               m.tags = append(m.tags, s)
               m.input.SetValue("")
         }

      // We handle errors just like any other message
      case error:
         m.err = msg
         return m, nil
   }

   m.input, cmd = m.input.Update(msg)
   return m, cmd
}

func (m addCmdModel) View() string {
   var sb strings.Builder
   for _, t := range m.tags {
      sb.WriteString("" + t + " ")
   }
   return fmt.Sprintf("%s\n%s\n(Enter blank to finish adding tags)\n", sb.String(), m.input.View())
}

func initialAddCmdModel() addCmdModel {
   ti := textinput.New()
   ti.Prompt = "Add tag: "
   ti.Focus()
   ti.CharLimit = 156
   ti.Width = 20

   return addCmdModel{
      tags: []string{},
      input: ti,
      err: nil,
   }
}

