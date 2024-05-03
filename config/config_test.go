package config

import (
   "path/filepath"
   "testing"

   "github.com/stretchr/testify/assert"
)


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
   assert.Equal(t, []string{"config_a.yaml", "config_b.yml", "00_conf.yaml", "01_conf.yml"}, yamlFiles)
}

