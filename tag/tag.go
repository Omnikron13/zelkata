package tag

import (
   "fmt"
   "os"
   "path/filepath"
   "strings"

   "github.com/omnikron13/zelkata/config"
   "github.com/omnikron13/zelkata/paths"

   "gopkg.in/yaml.v3"
)


// Tag is a struct that holds various metadata about a tag, which is primarily stored in plaintext YAML files in the
// user's data directory
type Tag struct {
   // Name is the canonical name of the tag. In practice tags are case-insensitive, but this is how the tag will
   // _output_ regardless of how it was input.
   Name string `yaml:"name"`

   // Description is a short human-readable description of the tag. As tags are generally created on the fly, they
   // generally don't start out with a description. Adding well thought out descriptions to tags is one of the small
   // admin tasks a user can do to improve the effectiveness of their Zelkata system. In theory there are a lot of ways
   // the vocabulary used in a description could be cross referenced with that of other tags, and of notes which don't
   // have the tag.
   Description string `yaml:"description,omitempty"`

   // Icon is a string holding a unicode sequence for an icon to be used to represent the tag.
   Icon string `yaml:"icon,omitempty"`

   // Aliases is a slice of strings that are alternative names for the tag. This allows a user to not be too concerned
   // with the exact name of a tag when adding it to a note, or when searching for notes with a tag, etc.
   Aliases []string `yaml:"aliases,omitempty"`


   // TODO: define an actual TagSet type for this kind of functionality?
   // Parents maps 'canonical' (human readable) tag names to their normalised (path friendly) form. It is implemented as
   // a map rather than a simple slice to essentially act as a set, with the normalised forms being just a convenience.
   Parents map[string]string `yaml:"parents,omitempty"`

   // Children ?

   // Notes is a slice of the UUIDs of notes that have this tag. The canonical connection between note and tag is
   // actually the note file, but it is obviously useful to be able to perform the reverse lookup.
   Notes []string `yaml:"notes"`
}


// Add either adds a new note ID to an existing tag, or creates a new tag with the given name, and its first note ID.
func Add(name, noteID string) error {
   tags, err := LoadAll()
   if err != nil {
      return err
   }
   tag := tags.Get(name)
   // Create the tag if it doesn't exist
   if tag == nil {
      tag = &Tag{Name: name, Notes: []string{}}
   }
   tag.Notes = append(tag.Notes, noteID)
   return tag.Save()
}


// LoadPath reads a tag file and returns a Tag struct
func LoadPath(filePath string) (*Tag, error) {
   t := Tag{}
   b, err := os.ReadFile(filePath)
   if err != nil {
      return nil, err
   }
   err = yaml.Unmarshal(b, &t)
   if err != nil {
      return nil, err
   }
   return &t, nil
}


// LoadName reads a tag file by name and returns a Tag struct.
// This is a convenience function that calls LoadPath with the full path and normalised tag name.
func LoadName(name string) (*Tag, error) {
   ext, err := config.Get[string]("tags.extension")
   if err != nil {
      return nil, err
   }
   return LoadPath(filepath.Join(paths.Tags(), normaliseName(name) + ext))
}


// MarshalYAML implements the yaml.Marshaler interface for the Tag struct.
func (t Tag) MarshalYAML() (interface{}, error) {
   data := map[string]any{
      "name": t.Name,
      "notes": t.Notes,
   }
   if len(t.Aliases) > 0 { data["aliases"] = t.Aliases }
   if t.Description != "" { data["description"] = t.Description }
   if t.Icon!= "" { data["icon"] = t.Icon }
   if len(t.Parents) > 0 {
      parents := make([]string, 0, len(t.Parents))
      for k := range t.Parents {
         parents = append(parents, k)
      }
      data["parents"] = parents
   }
   return interface{}(data), nil
}


// UnmarshalYAML implements the yaml.Unmarshaler interface for the Tag struct.
func (t *Tag) UnmarshalYAML(value *yaml.Node) error {
   data := map[string]any{}
   if err := value.Decode(&data); err != nil {
      return err
   }
   t.Name = data["name"].(string)
   if t.Name== "" {
      t = nil
      return fmt.Errorf("Missing tag name.")
   }
   t.Notes = []string{}
   for _, n := range data["notes"].([]any){
      t.Notes = append(t.Notes, n.(string))
   }
   if aliases, ok := data["aliases"]; ok {
      t.Aliases = []string{}
      for _, a := range aliases.([]any) {
         t.Aliases = append(t.Aliases, a.(string))
      }
   }
   if description, ok := data["description"]; ok {
      t.Description = description.(string)
   }
   if icon, ok := data["icon"]; ok {
      t.Icon = icon.(string)
   }
   if parents, ok := data["parents"]; ok {
      t.Parents = map[string]string{}
      for _, p := range parents.([]any) {
         t.Parents[p.(string)] = normaliseName(p.(string))
      }
   }
   return nil
}


// Save writes a Tag struct to a file in the tags directory.
func (t *Tag) Save() error {
   ext, err := config.Get[string]("tags.extension")
   if err != nil {
      return err
   }
   b, err := yaml.Marshal(t)
   if err != nil {
      return err
   }
   return os.WriteFile(filepath.Join(paths.Tags(), normaliseName(t.Name) + ext), b, 0600)
}


// normaliseName takes a tag name and returns a normalised (more path friendly, mostly) version of it.
func normaliseName(name string) string {
   return strings.ReplaceAll(strings.ToLower(name), " ", "-")
}

