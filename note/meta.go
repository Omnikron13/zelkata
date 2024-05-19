package note

import (
   "encoding/base32"
   "encoding/base64"
   "fmt"
   "strings"
   "time"

   "github.com/omnikron13/zelkata/config"

   "github.com/google/uuid"
   "gopkg.in/yaml.v3"
)

// Meta is the struct that holds various metadata about a note, which will be converted to the YAML front matter of the
// note file. It is intended to only be used as a field in the Note struct, but is kept separate for organisational
// purposes.
type Meta struct {
   // ID is a unique identifier for the note. It is generated when the note is created, with the specific type and
   // encoding details specified in the config.
   ID string

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


// NewMeta returns a new Meta struct with a new generated unique ID and the current date & time for Created.
func NewMeta() (m Meta) {
   m.ID = encodeID(generateID())
   m.Created = time.Now().UTC()
   return
}


// encodeID encodes an ID as a string ad specified in the config.
func encodeID(id []byte) string {
   format, err := config.Get[string]("notes.metadata.id.encode.format")
   if err != nil {
      panic("error getting config value notes.metadata.id.encode.format: " + err.Error())
   }
   switch format {
      case "base32":
         var encoding base32.Encoding
         charset, err := config.Get[string]("notes.metadata.id.encode.charset")
         if err != nil {
            panic("error getting config value notes.metadata.id.encode.charset: " + err.Error())
         }
         switch charset {
            case "StdEncoding":
               encoding = *base32.StdEncoding
            case "HexEncoding":
               encoding = *base32.HexEncoding
            default:
               // base32 requires precisely 32 characters
               if len(charset) != 32 {
                  panic(fmt.Errorf("invalid encoding charset length: %d", len(charset)))
               }
               // charset must not contain newline/carriage return characters or duplicates
               for i, c := range charset {
                  if c == '\n' || c == '\r' {
                     panic(fmt.Errorf("invalid encoding charset character: %c", c))
                  }
                  if strings.ContainsRune(charset[:i], c) {
                     panic(fmt.Errorf("duplicate encoding charset character: %c", c))
                  }
               }
               encoding = *base32.NewEncoding(charset)
         }
         padChar := base32.NoPadding
         pad, err := config.Get[bool]("notes.metadata.id.encode.padding")
         if err != nil {
            panic("error getting config value notes.metadata.id.encode.padding: " + err.Error())
         }
         if pad {
            padChar = base32.StdPadding
         }
         var b strings.Builder
         encoder := base32.NewEncoder(encoding.WithPadding(padChar), &b)
            defer encoder.Close()
         if _, err := encoder.Write(id); err != nil {
            panic("error encoding ID: " + err.Error())
         }
         return b.String()
      case "base64":
         var encoding base64.Encoding
         charset, err := config.Get[string]("notes.metadata.id.encode.charset")
         if err != nil {
            panic("error getting config value notes.metadata.id.encode.charset: " + err.Error())
         }
         switch charset {
            case "StdEncoding":
               encoding = *base64.StdEncoding
            case "URLEncoding":
               encoding = *base64.URLEncoding
            default:
               // base64 requires precisely 64 characters
               if len(charset) != 64 {
                  panic(fmt.Errorf("invalid encoding charset length: %d", len(charset)))
               }
               // charset must not contain newline/carriage return characters or duplicates
               for i, c := range charset {
                  if c == '\n' || c == '\r' {
                     panic(fmt.Errorf("invalid encoding charset character: %c", c))
                  }
                  if strings.ContainsRune(charset[:i], c) {
                     panic(fmt.Errorf("duplicate encoding charset character: %c", c))
                  }
               }
               encoding = *base64.NewEncoding(charset)
         }
         padChar := base64.NoPadding
         pad, err := config.Get[bool]("notes.metadata.id.encode.padding")
         if err != nil {
            panic("error getting config value notes.metadata.id.encode.padding: " + err.Error())
         }
         if pad {
            padChar = base64.StdPadding
         }
         var b strings.Builder
         encoder := base64.NewEncoder(encoding.WithPadding(padChar), &b)
            defer encoder.Close()
         if _, err := encoder.Write(id); err != nil {
            panic("error encoding ID: " + err.Error())
         }
         return b.String()
      default:
         panic(fmt.Errorf("unsupported encoding: %s", format))
   }
}


// generateID generates an ID as specified in the config.
func generateID() (id []byte) {
   idType, err := config.Get[string]("notes.metadata.id.type")
   if err != nil {
      panic("error getting config value notes.metadata.id.type: " + err.Error())
   }
   switch idType {
      case "UUIDv4":
         uuid := uuid.New()
         id = uuid[:]
      case "UUIDv7":
         uuid, err := uuid.NewV7()
         if err != nil {
            panic("error generating UUIDv7: " + err.Error())
         }
         id = uuid[:]
      default:
         panic(fmt.Errorf("unsupported ID type: %s", idType))
   }
   return
}


// GenFileName generates a filename for a note based on the Meta data.
func (m *Meta) GenFileName() string {
   // TODO: perhaps 'fail' silently on errors, treating a lack config value as a 'no'? with the default config, it
   //       should be impossible to get an error unless something has gone very wrong. better may be to add a function
   //       to the config package like `MustGet()` or `GetOrPanic()` to make access cleaner.
   sb := strings.Builder{}

   if prefixDate, err := config.Get[bool]("notes.filenames.prefix.date"); err != nil {
      panic("error getting config value notes.filenames.prefix.date: " + err.Error())
   } else if prefixDate {
      sb.WriteString(m.Created.Format(time.DateOnly))
      sb.WriteString(".")
   }

   if prefixTime, err := config.Get[bool]("notes.filenames.prefix.time"); err != nil {
      panic("error getting config value notes.filenames.prefix.time: " + err.Error())
   } else if prefixTime{
      sb.WriteString(m.Created.Format("15-04"))
      sb.WriteString(".")
   }

   sb.WriteString(m.ID)

   if suffixExtension, err := config.Get[string]("notes.filenames.suffix.extension"); err != nil {
      panic("error getting config value notes.filenames.suffix.extension: " + err.Error())
   } else {
      sb.WriteString(".")
      sb.WriteString(suffixExtension)
   }

   return sb.String()
}


// MarshalYAML implements the yaml.Marshaler interface for the Meta struct.
func (m *Meta) MarshalYAML() (any, error) {
   data := map[string]any{
      "id": m.ID,
      "tags": m.Tags,
   }
   if created, err := marshalTime(m.Created); err != nil {
      return nil, err
   } else {
      data["created"] = created
   }
   if m.Refs != nil {
      data["refs"] = m.Refs
   }
   if m.Format != nil {
      data["format"] = m.Format
   }
   if m.Title != nil {
      data["title"] = m.Title
   }
   return data, nil
}


// marshalTime is a helper function to marshal a time.Time into a string according to the config.
func marshalTime(t time.Time) (string, error) {
   format, err := config.Get[string]("notes.metadata.date.format")
   if err != nil {
      return "", err
   }

   switch format {
      case "Layout":
         return t.Format(time.Layout), nil
      case "ANSIC":
         return t.Format(time.ANSIC), nil
      case "UnixDate":
         return t.Format(time.UnixDate), nil
      case "RubyDate":
         return t.Format(time.RubyDate), nil
      case "RFC822":
         return t.Format(time.RFC822), nil
      case "RFC822Z":
         return t.Format(time.RFC822Z), nil
      case "RFC850":
         return t.Format(time.RFC850), nil
      case "RFC1123":
         return t.Format(time.RFC1123), nil
      case "RFC1123Z":
         return t.Format(time.RFC1123Z), nil
      case "RFC3339":
         return t.Format(time.RFC3339), nil
      case "RFC3339Nano":
         return t.Format(time.RFC3339Nano), nil
      case "Kitchen":
         return t.Format(time.Kitchen), nil
      case "Stamp":
         return t.Format(time.Stamp), nil
      case "StampMilli":
         return t.Format(time.StampMilli), nil
      case "StampMicro":
         return t.Format(time.StampMicro), nil
      case "StampNano":
         return t.Format(time.StampNano), nil
      case "DateTime":
         return t.Format(time.DateTime), nil
      case "DateOnly":
         return t.Format(time.DateOnly), nil
      case "TimeOnly":
         return t.Format(time.TimeOnly), nil
      default:
         return t.Format(format), nil
   }
}


// UnmarshalYAML implements the yaml.Unmarshaler interface for the Meta struct.
func (m *Meta) UnmarshalYAML(value *yaml.Node) (err error) {
   data := map[string]any{}

   err = value.Decode(&data)
   if err != nil { return }

   m.ID = data["id"].(string)
   if m.ID == "" {
      m = nil
      err = fmt.Errorf("Missing note ID.")
      return
   }

   m.Created = data["created"].(time.Time)

   tags := data["tags"].([]any)
   m.Tags = make([]string, 0, len(tags))
   for _, t := range tags {
      m.Tags = append(m.Tags, t.(string))
   }

   if r, ok := data["refs"]; ok {
      refs := make(map[string]string, len(r.(map[string]any)))
      m.Refs = &refs
      for k, v := range r.(map[string]any) {
         (*m.Refs)[k] = v.(string)
      }
   }

   if format, ok := data["format"]; ok {
      f := format.(string)
      m.Format = &f
   }

   if title, ok := data["title"]; ok {
      t := title.(string)
      m.Title = &t
   }
   return nil
}

