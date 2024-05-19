package note

import (
   "testing"
   "time"

   "github.com/stretchr/testify/assert"
)

func Test_ReadFile(t *testing.T) {
   now, err := time.Parse(time.DateTime, "2024-05-13 01:02:03")
   if err != nil { t.Fatalf("Failed to parse time: %s", err) }

   n, err := ReadFile("testdata/testnote.md")
   assert.Nil(t, err)
   assert.NotNil(t, n)
   assert.Equal(t, "0Q1W2E3R4T5Y6U7I8O9P", n.ID)
   assert.Equal(t, now, n.Created)
   assert.Equal(t, []string{"Foo", "Bar"}, n.Tags)
   assert.Equal(t, "A Test Note\n===========\n\nThis is a test note.\n\n\nDetails\n-------\n\nPrimarily this note file is to test loading and parsing of note files, ensuring both the YAML front matter and MarkDown\nbody are correctly read back into a Note object that can be manipulated and saved back to disk.\n\n", n.Body)
}

