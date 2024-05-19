// Package note covers the core functionality of creating new notes and saving them into the user's Zelkata catalogue.
package note

import (
   "bytes"
   "path/filepath"
   "os"

   "github.com/omnikron13/zelkata/paths"

   "gopkg.in/yaml.v3"
)

// TOTO: split out some of this waffle into the package docstring, and probably workshop some if it into the manual...

// Note is the primary struct that models the core unit of Zelkata,, the note.
// A note should consist of an atomic piece of information; a single idea, concept, or fact.
// By thinking in terms of Tags - _what_ the information relates to (and how this can be range from very broad general
// categories to very specific, niche, or even personal ones, often overlapping in subtle and complex ways that we may
// not even consciously consider much) it becomes clear that it is _disadvantageous_ to allow the scope of an single
// note to become too broad.
// See the manual for discussion on how this is similar to how the brain stores and retrieves information, how this
// might help you to become comfortable with the system, and how you can use the more advanced features of Zelkata in
// productive and flexible ways.
type Note struct {
   // Meta is the struct which, primarily, models the YAML front matter of the note file.
   Meta

   // Body is the actual content of the note.
   // It is perfectly acceptable for it to be composed of simple plaintext, but it is probably more likely to
   // usually utilise one of the many, many minimal markup languages that have proliferated over the years.
   // Most likely it will be MarkDown, but it could be AsciiDoc, reStructuredText, etc. Not to mention the hardly 
   // inconsiderable flavours of MarkDown itself.
   Body string
   // TODO: actually scrap this in favour pf a string balder?
}


// New creates a new Note with a new Meta struct and the given initial Body.
// It is of course probably most likely that a lot of notes will be created with blank bodies as users start manually
// creating notes, but passing an empty string is far from arduous, and it is likely to be convenient for more
// automated processes to be able to create notes with a body already in place.
func New(body string) Note {
   return Note{NewMeta(), body}
}


// genFile generates a byte slice representing the on-disk representation of the note.
func (n *Note) genFile() []byte {
   var b bytes.Buffer
   b.WriteString("---\n")
   yml, err := yaml.Marshal(n.Meta)
   if err != nil {
      // TODO: rework this to more idiomatic Go error handling
      panic(err)
   }
   b.Write(yml)
   // NOTE: the double newline after the front matter is intentional, as it prevents `glow` from rendering the front
   // matter as part of a line 1 heading in the rendered markdown.
   b.WriteString("...\n\n")
   b.WriteString(n.Body)
   return b.Bytes()
}


// ReadFile reads a note file from disk and returns a Note struct.
func ReadFile(path string) (n Note, err error) {
   b, err := os.ReadFile(path)
   if err != nil {
      return
   }
   return readBytes(b)
}


// readBytes reads a byte slice representing the on-disk representation of the note into the Note struct.
func readBytes(b []byte) (n Note, err error) {
   metaEnd := bytes.Index(b, []byte("\n...\n\n"))
   if err = yaml.Unmarshal(b[:metaEnd], &n.Meta); err != nil {
      return
   }
   n.Body = string(b[metaEnd+6:])
   return
}


// Save saves the note to the configured notes directory and filename.
func (n *Note) Save() error {
   return n.saveAs(filepath.Join(paths.Notes(), n.GenFileName()))
}


// saveAs saves the note to the given path.
func (n *Note) saveAs(path string) error {
   return os.WriteFile(path, n.genFile(), 0600)
}

