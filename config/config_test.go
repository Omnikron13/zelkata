package config

import (
   "os"
   "path/filepath"
   "testing"

   "github.com/stretchr/testify/assert"
)


// Test_defaultConfig tests that the embedded defaultConfig file correctly matches the contents of defaults.yaml
func Test_defaultConfig(t *testing.T) {
   bytes, _ := os.ReadFile("defaults.yaml")
   assert.Equal(t, defaultConfig, bytes)
}


func Test_findYAMLFilesIn(t *testing.T) {
   testdata, _ := filepath.Abs("testdata")
   configPath := filepath.Join(testdata, "xdg_config_home", "zelkata")

   var yamlFiles []string
   for _, f := range findYAMLFilesIn(configPath) {
      yamlFiles = append(yamlFiles, filepath.Base(f))
   }
   for _, f := range findYAMLFilesIn(filepath.Join(configPath, "conf.d")) {
      yamlFiles = append(yamlFiles, filepath.Base(f))
   }
   assert.Equal(t, []string{"config_e.yaml", "config_f.yml", "04_conf.yaml", "05_conf.yml"}, yamlFiles)
}

