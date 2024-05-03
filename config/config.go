package config

import (
   "io/fs"
   "os"
   "path/filepath"
   "slices"

   "github.com/adrg/xdg"
)


// findYAMLFiles finds YAML files in the XDG configuration directories.
func findYAMLFiles() (yamlFiles []string) {
   for _, dir := range append(xdg.ConfigDirs, xdg.ConfigHome) {
      yamlFiles = append(yamlFiles, findYAMLFilesIn(filepath.Join(dir, "zelkata"))...)
      yamlFiles = append(yamlFiles, findYAMLFilesIn(filepath.Join(dir, "zelkata", "conf.d"))...)
   }
   return
}


// findYAMLFilesIn finds YAML files in the specified directory.
func findYAMLFilesIn(dir string) (yamlFiles []string) {
   configDir := os.DirFS(dir)
   for _, glob := range []string{"*.yml", "*.yaml"} {
      if f, err := fs.Glob(configDir, glob); f != nil && err == nil {
         for i := range f {
            yamlFiles = append(yamlFiles, filepath.Join(dir, f[i]))
         }
      }
   }
   slices.Sort(yamlFiles)
   return
}

