package note

import (
   "testing"

   "github.com/stretchr/testify/assert"
)

func Test_ReadFile(t *testing.T) {
   n, err := ReadFile("testdata/testnote.md")
   assert.Nil(t, err)
   assert.NotNil(t, n)
   assert.Equal(t, "123456789", n.ID)
   assert.Equal(t, "2024-05-13T01:02:03Z00:00", n.Created)
   assert.Equal(t, []string{"Foo", "Bar"}, n.Tags)
   assert.Equal(t, "A Test Note\n===========\n\nThis is a test note.\n\n\nDetails\n-------\n\nPrimarily this note file is to test loading and parsing of note files, ensuring both the YAML front matter and MarkDown\nbody are correctly read back into a Note object that can be manipulated and saved back to disk.\n\n", n.Body)
}

