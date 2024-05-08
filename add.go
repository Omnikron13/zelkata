package main

import (
   "context"
   "fmt"
   "os"
   "os/exec"
   "path/filepath"
   "strings"

   "github.com/omnikron13/zelkata/note"
   "github.com/omnikron13/zelkata/config"

   "github.com/urfave/cli/v3"
   tea "github.com/charmbracelet/bubbletea"
   "github.com/charmbracelet/bubbles/textinput"
   "github.com/adrg/xdg"
)


var conf *config.Config


func addCmd(ctx context.Context, cmd *cli.Command) error {
   // Initialise the config hierarchy
   var err error
   conf, err = config.Init()
   if err != nil {
      panic(err)
   }

   // TODO: default to a bubbletea(bubbles) TextArea, with a hotkey to launch a full editor?
   // This sets up launching an external editor to write the note body, which is temporarily stored in a state file,
   // which potentially also acts as a draft file if the user saves while editing but the add process is interrupted.
   stateDir := filepath.Join(xdg.StateHome, "zelkata")
   if err := os.MkdirAll(stateDir, 0700); err != nil {
      return err
   }
   newNoteFile := filepath.Join(stateDir, "new-note.md")
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

   // TODO: add a paths module, I think, as we really don't want to be repeating this song and dance all over the place
   // Get the path to the notes directory.
   path, err := config.Get[string](conf, "data-directory")
   if err != nil {
      panic(err)
   }
   // filepath sadly has no convenience function for expanding environment variables...
   path = os.ExpandEnv(path)
   // Converting to an absolute path may not really be worth the time, but we don't want any confusing behaviour
   // cropping up downt the line.
   if path, err = filepath.Abs(filepath.Join(path, "notes")); err != nil {
      panic(err)
   }

   // Ensure the notes directory exists, creating it if necessary (including all directories along the way)
   if err := os.MkdirAll(path, 0700); err != nil {
      return err
   }

   // Generate a byte slice of the note file and write it to the notes directory
   b := note.GenFile()
   if err := os.WriteFile(filepath.Join(path, note.GenFileName()), b, 0600); err != nil {
      return err
   }

   // Return nil if everything went well
   return nil
}


// Bubbletea code for adding the tags to a Note
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

