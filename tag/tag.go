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
   Name string

   // Description is a short human-readable description of the tag. As tags are generally created on the fly, they
   // generally don't start out with a description. Adding well thought out descriptions to tags is one of the small
   // admin tasks a user can do to improve the effectiveness of their Zelkata system. In theory there are a lot of ways
   // the vocabulary used in a description could be cross referenced with that of other tags, and of notes which don't
   // have the tag.
   Description string

   // Virtual indicates whether the tag is able to be directly assigned to notes. This  can be used for more abstract
   // concepts up the tag hierarchy, or for tags automatically applied under certain conditions.
   Virtual bool

   // Icon is a string holding a unicode sequence for an icon to be used to represent the tag.
   Icon string

   // Aliases is a slice of strings that are alternative names for the tag. This allows a user to not be too concerned
   // with the exact name of a tag when adding it to a note, or when searching for notes with a tag, etc.
   Aliases []string


   // TODO: define an actual TagSet type for this kind of functionality?
   // Parents maps 'canonical' (human readable) tag names to their normalised (path friendly) form. It is implemented as
   // a map rather than a simple slice to essentially act as a set, with the normalised forms being just a convenience.
   Parents map[string]string

   // Children ?

   // Relations is a set of Tags that have a less direct or hierarchical relationship to this tag. This could be things
   // culturally or personally associated with the tag, or things which share aspects in common (and may indeed be
   // related like cousins if the Tag graph is filled out enough) but have no direct vertical relationship.
   // They key is the related tags name, and the value in this instance is a short description of the relationship.
   Relations map[string]string

   // Notes is a slice of the UUIDs of notes that have this tag. The canonical connection between note and tag is
   // actually the note file, but it is obviously useful to be able to perform the reverse lookup.
   Notes []string
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


// genFileName generates a filename for a tag file based on the tag name.
func (t *Tag) genFileName() (name string, err error) {
   ext, err := config.Get[string]("tags.extension")
   if err != nil {
      return
   }
   name = t.normalisedName() + ext
   return
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
   if t.Virtual { data["virtual"] = t.Virtual }
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
   if len(t.Relations) > 0 {
      relations := make([]struct{Name, Description string}, 0 , len(t.Relations))
      for k, v := range t.Relations {
         relations = append(relations, struct {Name, Description string}{Name:k, Description:v})
      }
      data["relations"] = relations
   }
   return interface{}(data), nil
}


// normalisedName returns the normalised name of a tag.
func (t *Tag) normalisedName() string {
   return normaliseName(t.Name)
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
   if _, ok := data["virtual"]; ok {
      t.Virtual = true
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
   if relations, ok := data["relations"]; ok {
      t.Relations = map[string]string{}
      for _, r := range relations.([]any) {
         m := r.(map[string]any)
         t.Relations[m["name"].(string)] = m["description"].(string)
      }
   }
   return nil
}


// Save writes a Tag struct to a file in the tags directory.
func (t *Tag) Save() error {
   name, err := t.genFileName()
   if err != nil {
      return err
   }
   path := filepath.Join(paths.Tags(), name)
   return t.saveAs(path)
}


// saveAs writes a Tag struct to an arbitrary file path.
func (t *Tag) saveAs(filePath string) error {
   b, err := yaml.Marshal(t)
   if err != nil {
      return err
   }
   return os.WriteFile(filePath, b, 0600)
}


// normaliseName takes a tag name and returns a normalised (more path friendly, mostly) version of it.
func normaliseName(name string) string {
   // TODO: config allowed characters and strip others? would need to e.g. add a hash to the end of the name to ensure
   // uniqueness in the case of a collision (64bit xxhash base32 encoded = 15 chars, truncated to 32/16 bits = 8/5)
   return strings.ReplaceAll(strings.ToLower(name), " ", "-")
}

