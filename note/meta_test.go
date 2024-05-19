package note

import (
   "testing"

   "github.com/stretchr/testify/assert"
   "gopkg.in/yaml.v3"
)

func Test_MarshalYAML(t *testing.T) {
   t.Run("simple meta", func(t *testing.T) {
      meta := Meta {
         ID: "123456789",
         Tags: []string{"Foo", "Bar"},
         Created: "2024-05-13T01:02:03Z00:00",
      }
      data, err := yaml.Marshal(&meta)
      assert.Nil(t, err)
      assert.Equal(t, "created: 2024-05-13T01:02:03Z00:00\nid: \"123456789\"\ntags:\n    - Foo\n    - Bar\n", string(data))
   })

   t.Run("complex meta", func(t *testing.T) {
      modified := "2024-05-13T01:02:03Z00:00"
      title := "Test Note"
      format := "AsciiDoc"
      meta := Meta {
         ID: "123456789",
         Tags: []string{"Foo", "Bar"},
         Created: "2024-05-13T01:02:03Z00:00",
         Modified: &modified,
         Refs: &map[string]string{
            "Website": "https://example.com",
            "Book": "ISBN 1234567890",
         },
         Format: &format,
         Title: &title,
      }
      data, err := yaml.Marshal(&meta)
      assert.Nil(t, err)
      assert.Equal(t, "created: 2024-05-13T01:02:03Z00:00\nformat: AsciiDoc\nid: \"123456789\"\nmodified: 2024-05-13T01:02:03Z00:00\nrefs:\n    Book: ISBN 1234567890\n    Website: https://example.com\ntags:\n    - Foo\n    - Bar\ntitle: Test Note\n", string(data))
   })
}

func Test_UnmarshalYAML(t *testing.T) {
   t.Run("simple meta", func(t *testing.T) {
      data := "created: 2024-05-13T01:02:03Z00:00\nid: \"123456789\"\ntags:\n    - Foo\n    - Bar\n"
      meta := Meta{}
      err := yaml.Unmarshal([]byte(data), &meta)
      assert.Nil(t, err)
      assert.Equal(t, Meta{
         ID: "123456789",
         Tags: []string{"Foo", "Bar"},

         Created: "2024-05-13T01:02:03Z00:00",
      }, meta)
   })

   t.Run("complex meta", func(t *testing.T) {
      data := "created: 2024-05-13T01:02:03Z00:00\nformat: AsciiDoc\nid: \"123456789\"\nmodified: 2024-05-13T01:02:03Z00:00\nrefs:\n    Book: ISBN 1234567890\n    Website: https://example.com\ntags:\n    - Foo\n    - Bar\ntitle: Test Note\n"
      meta := Meta{}
      err := yaml.Unmarshal([]byte(data), &meta)
      assert.Nil(t, err)
      modified := "2024-05-13T01:02:03Z00:00"
      title := "Test Note"
      format := "AsciiDoc"
      expected := Meta{
         ID: "123456789",
         Tags: []string{"Foo", "Bar"},
         Created: "2024-05-13T01:02:03Z00:00",
         Modified: &modified,
         Refs: &map[string]string{
            "Website": "https://example.com",
            "Book": "ISBN 1234567890",
         },
         Format: &format,

         Title: &title,
      }
      assert.Equal(t, expected, meta)
   })
}

