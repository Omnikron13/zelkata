package tag

import (
   "path/filepath"
   "testing"

   "github.com/stretchr/testify/assert"
)


func Test_LoadPath(t *testing.T) {
   path := filepath.Join("testdata", "test.tag.yaml")
   tag, err := LoadPath(path)
   assert.Nil(t, err)
   assert.Equal(t, Tag{
      Name: "Test Tag",
      Description: "An example tag for testing purposes.",
      Icon: "󰓹",
      Notes: []string{
         "QWERTYUIOP",
         "ASDFGHJKLZ",
      },
   }, *tag)
}


func Test_normaliseName(t *testing.T) {
   assert.Equal(t, "test-tag-name", normaliseName("Test TAG name"))
   assert.Equal(t, "already-normalised", normaliseName("already-normalised"))
 }


// TODO: mock file access and test; Add(), LoadName(), Save()

