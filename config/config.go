package config

import (
   _ "embed"
   "io/fs"
   "maps"
   "os"
   "path/filepath"
   "slices"

   "github.com/adrg/xdg"
   "gopkg.in/yaml.v3"
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


// TODO: probably flatten this out and not bother with an actual struct? Otherwise probably need to make it a singleton
//       anyway to keep things DRY. (alternatively, a Config struct could be a full typed model of the config?)
type Config struct {
   // TODO: replace with a proper stack data structure
   // filesData holds the configuration data from the YAML files, read into raw byte slices as essentially a 'snapshot'
   // of the config file hierarchy on initialisation.
   filesData [][]byte

   // yamlData is a (currently) generic container for unmarshalled YAML data, which should be merged in order into it
   // as the actual config values are needed.
   yamlData map[string]any
}


// Init initialises the config by reading the config file hierarchy from all the applicable YAML files which can
// be found. At the bottom of the hierarchy the embedded default config is inserted.
func Init() (*Config, error) {
   // TODO: consider inverting the order of the filesData slice so the default config is at the end of the hierarchy,
   //       and the hierarchy is read from most important to least important.
   filesData := [][]byte{defaultConfig}

   for _, yamlFile := range findYAMLFiles() {
      data, err := os.ReadFile(yamlFile)
      if err != nil {
         return nil, err
      }
      filesData = append(filesData, data)
   }

   return &Config{filesData: filesData}, nil
}


// unmarshalNext unmarshals the next YAML document from the byte slice and integrates it with the existing config data.
// It operates top-down, effectively lazy-loading the config data.
func (c *Config) unmarshalNext() error {
   if len(c.filesData) == 0 {
      // TODO: check if the slice is wasting memory by not being truncated, and free it up if so?
      // Returning nil because it's irrelevant if there's no more data to unmarshal.
      return nil
   }

   // Generic container for unmarshalled YAML data. Might be better to define a struct for this, but I think it's ok
   // for the time being...
   yml := map[string]any{}

   last := len(c.filesData) - 1
   if err := yaml.Unmarshal(c.filesData[last], &yml); err != nil {
      return err
   }
   c.filesData = c.filesData[:last]

   // TODO: check that this easy approach to merging doesn't become a problem with arrays (which I presume to be the
   //       most likely source of probably undesirable behaviour)
   maps.Copy(yml, c.yamlData)
   c.yamlData = yml

   return nil
}


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

