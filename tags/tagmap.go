package tags

import (
   "errors"
   "os"
   "path/filepath"

   "github.com/omnikron13/zelkata/note"
   "github.com/omnikron13/zelkata/paths"

   "k8s.io/apimachinery/pkg/util/sets"
)

// TagMap is a map of tag names to Tag structs, representing all tags in the system.
// It includes aliases, and normalises all names before lookup.
type TagMap map[string]*Tag


// LoadAll reads all tag files into a TagMap.
func LoadAll() (TagMap, error) {
   tm := TagMap{}
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
      if err := tm.Add(tag.Name, tag); err != nil {
         return nil, err
      }
      for _, t := range tag.Aliases {
         if err := tm.Add(t, tag); err != nil {
            return nil, err
         }
      }
   }
   return tm, nil
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
      tag.Notes = sets.New[string]()
   }

   // Read notes and add their IDs to the appropriate tags, creating new tags as necessary
   files, err := os.ReadDir(paths.Notes())
   if err != nil {
      return err
   }
   for _, file := range files {
      note, err := note.ReadFile(filepath.Join(paths.Notes(), file.Name()))
      if  err != nil {
         return err
      }
      for t := range note.Tags {
         tag := m.Get(t)
         if tag == nil {
            tag = &Tag{Name: t, Notes: sets.New[string]()}
            _ = m.Add(t, tag)
         }
         tag.Notes.Insert(note.ID)
      }
   }
   return nil
}


// Save writes all (non-alias) Tag structs in the TagMap to files in the tags directory.
func (m *TagMap) Save() error {
   for name, tag := range *m {
      if name != normaliseName(tag.Name) {
         continue
      }
      if err := tag.Save(); err != nil {
         return err
      }
   }
   return nil
}

