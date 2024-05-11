package tag

import (
   "errors"
   "os"
   "path/filepath"

   "github.com/omnikron13/zelkata/paths"
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

