package note

import (
   "encoding/base32"
   "fmt"
   "strings"
   "time"

   "github.com/google/uuid"
)


// Meta is the struct that holds various metadata about a note, which will be converted to the YAML front matter of the
// note file. It is intended to only be used as a field in the Note struct, but is kept separate for organisational
// purposes.
type Meta struct {
   // ID is simply a UUID generated to uniquely identify the note.
   ID uuid.UUID

   // Created is the date & time the note was originally created.
   Created time.Time

   // Tags are how notes are (primarily) categorised. They are _technically_ optional, but in practice they should
   // always be used, as they are what creates the hyperlinked web of knowledge that Zelkata is based on.
   // Notes that have no tags should have a 'virtual tag' of 'untagged' or similar applied to them, and the user 
   // warned that they are of little use until linked into the overall system.
   //
   // This simple string slice is very likely to be a temporary stepping stone and placeholder until the much more
   // complex and powerful tag system is implemented.
   Tags []string


   // The following are the optional fields, so may need special handling when converting to/from YAML, etc.
   // They could also perhaps be shifted out to a separate struct and composed here for the sake of keeping the source
   // as clean as possible..


   // Modified is the date & time the note was last modified. It is optional as a lot of notes will likely never be
   // modified, and if Zelkata is set up to use Git (or other possible VCS) then changes will be tracked there.
   // Assuming it does actually get used, it maybe should be replaced with a slice?
   Modified *string

   // Refs are a secondary layer of of hyperlinking that can be necessary in some cases. They are entirely optional for
   // most notes, but can be crucial for some. They can be URLs pointing to related information or resources. They
   // could also be some form of URN identifying a book,research paper, etc. They could even simply be path pointing
   // to a file or directory on the local file system which is related to the note.
   //
   // This basic key:value map is very likely to just serve as a placeholder until a more robust and powerful system
   // is implemented.
   Refs *map[string]string

   // Format can be used to store the format of the note, e.g. MarkDown, AsciiDoc, etc. if it can't be inferred from
   // the note itself, user configuration, or file extension hints.
   Format *string

   // Title could very well be nil, as the note itself is likely to be in a markup format allowing a title there.
   // The fact this is an option raises a point of caution in that duplication (or worse inconsistency) could arise.
   // Plain text notes (etc?) _might_ want an explicit title though?
   Title *string
}


// generateUUID generates a UUID and returns it as simply an array of 16-bytes (128-bits).
// It is sticking with the UUID v4 format for now, though will probably be configurable in the future, and may also
// start defaulting to v7, when it is satisfactorily finalised.
func generateUUID() [16]byte {
   return uuid.New()
}


// NewMeta returns a new Meta struct with a UUID and the current date and time.
func NewMeta() Meta {
   return Meta{
      ID: generateUUID(),
      Created: time.Now(),
   }
}


// base32UUID encodes the UUID of the Meta struct as a base32 string.
// TODO: can probably generalise this if/when I add a config for the encoding, should just need a switch I believe.
func (m *Meta) base32UUID() string {
   var b strings.Builder
   encoder := base32.NewEncoder(base32.StdEncoding.WithPadding(base32.NoPadding), &b)
      defer encoder.Close()
   encoder.Write(m.ID[:])
   return b.String()
}


// GenFileName generates a filename for a note based on the Meta data.
func (m *Meta) GenFileName() string {
   // TODO: move date & time prefixing to config
   // TODO: add config (&default?) for encoding the UUID to a more concise form
   // TODO: add config for file extension? Also probably depend on what the actual Note is told the format is.
   return fmt.Sprintf("%s.%s.%s.md", m.Created.Format(time.DateOnly), m.Created.Format("15-04"), m.ID)
}

