package tag

import (
   "errors"
   "os"
   "path/filepath"

   "github.com/omnikron13/zelkata/paths"

   "gopkg.in/yaml.v3"
)

// TagMap is a map of tag names to Tag structs, representing all tags in the system.
// It includes aliases, and normalises all names before lookup.
type TagMap map[string]*Tag


// LoadAll reads all tag files into a TagMap.
func LoadAll() (TagMap, error) {
   tags := TagMap{}
   files, err := os.ReadDir(paths.Tags())
   if err != nil {
      return nil, err
   }
   for _, file := range files {
      if !file.Type().IsRegular() {
         continue
      }
      tag, err := LoadPath(filepath.Join(paths.Tags(), file.Name()))
      if err != nil {
         return nil, err
      }
      if err := tags.Add(tag.Name, tag); err != nil {
         return nil, err
      }
      for _, t := range tag.Aliases {
         if err := tags.Add(t, tag); err != nil {
            return nil, err
         }
      }
   }
   return tags, nil
}


// Add adds a new reference to a given tag to the TagMap with the given name, returning an error if a reference with
// that name already exists.
func (m *TagMap) Add(name string, tag *Tag) error {
   name = normaliseName(name)
   if _, exists := (*m)[name]; exists {
      return errors.New("Tag already exists")
   }
   (*m)[name] = tag
   return nil
}


// Get returns a reference from the TagMap by name, or nil if it doesn't exist. The name is normalised before lookup.
func (m *TagMap) Get(name string) *Tag {
   return (*m)[normaliseName(name)]
}


// Reindex clears the Notes field of all tags in the TagMap, then repopulates them by scanning the notes directory.
func (m *TagMap) Reindex() error {
   // Clear existing note references
   for name, tag := range *m {
      if name != normaliseName(tag.Name) {
         continue
      }
      tag.Notes = []string{}
   }

   // Read notes and add their IDs to the appropriate tags, creating new tags as necessary
   files, err := os.ReadDir(paths.Notes())
   if err != nil {
      return err
   }
   for _, file := range files {
      // TODO: rework Note/Meta so they can be unmarshalled into properly and replace this hacky generic map
      note := map[string]any{}
      b, err := os.ReadFile(filepath.Join(paths.Notes(), file.Name()))
      if err != nil {
         return err
      }
      if err := yaml.Unmarshal(b, note); err != nil {
         return err
      }
      tags, ok := note["tags"].([]any)
      if !ok {
         return errors.New("Tags field is not an array")
      }
      for _, t := range tags {
         t, ok := t.(string)
         if !ok {
            return errors.New("Tag is not a string")
         }
         tag := m.Get(t)
         if tag == nil {
            tag = &Tag{Name: t, Notes: []string{}}
            _ = m.Add(t, tag)
         }
         tag.Notes = append(tag.Notes, note["id"].(string))
      }
   }
   return nil
}

