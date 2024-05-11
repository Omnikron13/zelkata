package tag

import (
   "testing"

   "github.com/stretchr/testify/assert"
)


func Test_TagMap_Add(t *testing.T) {
   tags := TagMap{}
   tag := &Tag{Name: "Test Tag"}
   err := tags.Add("Test Tag", tag)
   assert.Nil(t, err)
   err = tags.Add("Test Tag", tag)
   assert.NotNil(t, err)
}

func Test_TagMap_Get(t *testing.T) {
   tags := TagMap{}
   tag := &Tag{Name: "Test Tag"}
   if err := tags.Add("Test Tag", tag); err != nil {
      t.Skipf("Error adding tag to TagMap: %v", err)
   }
   assert.Equal(t, tag, tags.Get("Test Tag"))
   assert.Nil(t, tags.Get("Non-existent Tag"))
}

