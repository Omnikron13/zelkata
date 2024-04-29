package config

import (
   "io/fs"
   "os"
   "slices"

   "github.com/adrg/xdg"
)


// TODO: refactor this monstrosity so it isn't reusing quite as disgusting amount of code...
func findYAMLFiles() (yamlFiles []string) {
   for _, dir := range append(xdg.ConfigDirs, xdg.ConfigHome) {
      configDir := os.DirFS(dir + "/zelkata")
      println("looking in", dir + "/zelkata", configDir, "for YAML files")
      if f, err := fs.Glob(configDir, "*.yml"); f != nil && err == nil {
         for i := range f {
            yamlFiles = append(yamlFiles, dir + "/zelkata/" + f[i])
         }
      }
      if f, err := fs.Glob(configDir, "*.yaml"); f != nil && err == nil {
         for i := range f {
            yamlFiles = append(yamlFiles, dir + "/zelkata/" + f[i])
         }
      }

      conf_d := []string{}
      if f, err := fs.Glob(configDir, "\\conf.d/*.yml"); f != nil && err == nil {
         for i := range f {
            conf_d = append(conf_d, dir + "/zelkata/" + f[i])
         }
      }
      if f, err := fs.Glob(configDir, "\\conf.d/*.yaml"); f != nil && err == nil {
         for i := range f {
            conf_d = append(conf_d, dir + "/zelkata/" + f[i])
         }
      }

      slices.Sort(conf_d)
      yamlFiles = append(yamlFiles, conf_d...)
   }
   return
}

