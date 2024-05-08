package config

import (
   "os"
   "path/filepath"
   "testing"

   "github.com/adrg/xdg"
   "github.com/stretchr/testify/assert"
)


// Test_defaultConfig tests that the embedded defaultConfig file correctly matches the contents of defaults.yaml
func Test_defaultConfig(t *testing.T) {
   bytes, _ := os.ReadFile("defaults.yaml")
   assert.Equal(t, defaultConfig, bytes)
}


func Test_unmarshalNext(t *testing.T) {
   // TODO: populate file(s) in the testdata directory and test the merging of multiple files
   c := Config{filesData: [][]byte{defaultConfig}}
   err := c.unmarshalNext()
   assert.Nil(t, err)
   // Ensure nothing strange happens if called with no more data
   err = c.unmarshalNext()
   assert.Nil(t, err)
   // This is the most basic test possible, but it is at least a sanity check...
   assert.NotEmpty(t, c.yamlData)
   // TODO: try and implement a clean & maintainable way to much more robustly test this
   assert.Equal(t, "$XDG_DATA_HOME/zelkata", c.yamlData["data-directory"])
}


func Test_Get(t *testing.T) {
   c := Config{filesData: [][]byte{defaultConfig}}
   if err := c.unmarshalNext(); err != nil {
      t.Skipf("Error initialising Config instance to test: %v", err)
   }

   if c.yamlData == nil {
      t.Skip("Config instance inexplicably has no data to test")
   }

   t.Run("flat", func(t *testing.T) {
      v, err := Get[string](&c, "data-directory")
      assert.Nil(t, err)
      assert.Equal(t, "$XDG_DATA_HOME/zelkata", v)
   })

   t.Run("nested", func(t *testing.T) {
      v, err := Get[string](&c, "notes.metadata.id.type")
      assert.Nil(t, err)
      assert.Equal(t, "UUIDv4", v)
   })

   t.Run("type mismatch", func(t *testing.T) {
      v, err := Get[string](&c, "notes.metadata.id.encode.padding")
      assert.NotNil(t, err)
      assert.Equal(t, "", v)
      t.Logf("Error (expected): %v", err)
   })

   t.Run("non-existant key flat", func(t *testing.T) {
      unobtanium, err := Get[string](&c, "unobtanium")
      assert.NotNil(t, err)
      assert.Equal(t, "", unobtanium)
      t.Logf("Error (expected): %v", err)
   })

   t.Run("non-existant key nested", func(t *testing.T) {
      eludium, err := Get[string](&c, "note.metadata.id.encode.eludium")
      assert.NotNil(t, err)
      assert.Equal(t, "", eludium)
      t.Logf("Error (expected): %v", err)
   })

   t.Run("trigger unmarshal", func(t *testing.T) {
      c := Config{filesData: [][]byte{defaultConfig}}
      v, err := Get[string](&c, "data-directory")
      assert.Nil(t, err)
      assert.Equal(t, "$XDG_DATA_HOME/zelkata", v)
   })
}


func Test_findYAMLFiles(t *testing.T) {
   testdata, _ := filepath.Abs("testdata")
   xdg.ConfigHome = filepath.Join(testdata, "xdg_config_home")
   xdg.ConfigDirs[0] = filepath.Join(testdata, "xdg_config_dirs_0")
   xdg.ConfigDirs[1] = filepath.Join(testdata, "xdg_config_dirs_1")

   var yamlFiles []string
   for _, f := range findYAMLFiles() {
      yamlFiles = append(yamlFiles, filepath.Base(f))
   }
   assert.Equal(t, []string{"config_a.yaml", "config_b.yml", "00_conf.yaml", "01_conf.yml", "config_c.yml", "config_d.yaml", "02_conf.yml", "03_conf.yaml", "config_e.yaml", "config_f.yml", "04_conf.yaml", "05_conf.yml"}, yamlFiles)
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

