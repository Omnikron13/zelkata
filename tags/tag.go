package tags

import (
   . "cmp"
   "fmt"
   "os"
   "path/filepath"

   "github.com/omnikron13/zelkata/config"
   "github.com/omnikron13/zelkata/note"
   "github.com/omnikron13/zelkata/paths"

   "gopkg.in/yaml.v3"
   "k8s.io/apimachinery/pkg/util/sets"
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


   // Parents is a set of Tags that can be considered to be directly 'above' this tag in a perceptual hierarchy. It is
   // probably best to lean on the broad side when considering what conceptually constitutes a 'parent'. For example,
   // 'mathematics' could be considered a parent of 'algebra', 'geometry', etc. but also 'physics', 'engineering', even
   // 'music'.
   Parents sets.Set[string]

   // Children ?

   // Relations is a set of Tags that have a less direct or hierarchical relationship to this tag. This could be things
   // culturally or personally associated with the tag, or things which share aspects in common (and may indeed be
   // related like cousins if the Tag graph is filled out enough) but have no direct vertical relationship.
   // They key is the related tags name, and the value in this instance is a short description of the relationship.
   Relations map[string]string

   // Notes is a set of the UUIDs of notes that have this tag. The canonical connection between note and tag is
   // actually the note file, but it is obviously useful to be able to perform the reverse lookup.
   Notes sets.Set[string]
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
      tag = &Tag{Name: name, Notes: sets.New[string]()}
   }
   tag.Notes.Insert(noteID)
   return tag.Save()
}


// AddNote adds a note ID to the tag.
// n is a Note struct to take the ID from.
func (t *Tag) AddNote(n *note.Note) {
   t.Notes.Insert(n.ID)
}


// GenFileName generates a filename for a tag file based on the tag name.
func (t *Tag) GenFileName() (name string, err error) {
   ext := config.GetOrPanic[string]("tags.metadata.extension")
   name = fmt.Sprintf("%s%s", t.normalisedName(), ext)
   return
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


// LoadOrCreate reads a tag file by name and returns a Tag struct, or creates a new Tag struct with the given name if
// the tag doesn't already exist. This mirrors how a user would generally treat tags when adding them to their new note
// before saving it; near-zero friction when making notes is paramount.
func LoadOrCreate(name string) (t *Tag, err error) {
   var tags TagMap
   if tags, err = LoadAll(); err != nil { return } else
      { t = Or(tags.Get(name), &Tag{Name: name}) }
   return
}


// LoadPath reads a tag file and returns a Tag struct
func LoadPath(filePath string) (t *Tag, err error) {
   var b []byte
   if b, err = os.ReadFile(filePath); err != nil { return } else
      { err = yaml.Unmarshal(b, &t) }
   return
}


// MarshalYAML implements the yaml.Marshaler interface for the Tag struct.
func (t Tag) MarshalYAML() (interface{}, error) {
   data := map[string]any{
      "name": t.Name,
      "notes": sets.List(t.Notes),
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


// Save writes a Tag struct to a file in the tags directory.
func (t *Tag) Save() error {
   name, err := t.GenFileName()
   if err != nil {
      return err
   }
   path := filepath.Join(paths.Tags(), name)
   return t.SaveAs(path)
}


// SaveAs writes a Tag struct to an arbitrary file path.
func (t *Tag) SaveAs(filePath string) error {
   b, err := yaml.Marshal(t)
   if err != nil {
      return err
   }
   if err = os.WriteFile(filePath, b, 0600); err == nil { return nil }
   return fmt.Errorf("Failure writing tag file at %s during Save()", filePath)
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
   t.Notes = sets.New[string]()
   for _, n := range data["notes"].([]any){
      t.Notes.Insert(n.(string))
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
      t.Parents = sets.New[string]()
      for _, p := range parents.([]any) {
         t.Parents.Insert(p.(string))
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

