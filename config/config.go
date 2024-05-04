package config

import (
   _ "embed"
   "io/fs"
   "os"
   "path/filepath"
   "slices"

   "github.com/adrg/xdg"
)

// defaultConfig holds the default configuration values for Zelkata.
// Embedding it directly into the binary makes it easy to maintain a canonical DRY version of the entire default
// configuration. This allows devs to view and maintain the default configuration from a central location in a format
// which can be extensively commented/documented inline, which should lead to greater consistency, prevent redundancy,
// reduce maintenance burden, and make it easier to understand the system as a whole.
// Additionally, always-correct default values can be accessed at runtime, for e.g. outputting a complete commented
// config template, dynamically outputting help text for config options, etc.
//go:embed defaults.yaml
var defaultConfig []byte


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

