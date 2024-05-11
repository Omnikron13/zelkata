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
   assert.Equal(t, Tag{Name: "Test Tag", Description: "An example tag for testing purposes.", Notes: []string{"QWERTYUIOP", "ASDFGHJKLZ"}}, *tag)
}

