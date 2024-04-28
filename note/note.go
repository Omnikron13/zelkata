// Package note covers the core functionality of creating new notes and saving them into the user's Zelkata catalogue.
package note

import (
   _"io"
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


// Write implements the io.Writer interface, allowing a Note to be written to a file (or elsewhere, theoretically).
func (n Note) Write(p []byte) (int, error) {
   // TODO
   panic("Unimplemented!")
   //return 0, nil
}


// TODO: fill out this utter stub o stud/interface/