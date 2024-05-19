package note

import (
   "testing"
   "time"

   "github.com/stretchr/testify/assert"
   "gopkg.in/yaml.v3"
)

func Test_GenFileName(t *testing.T) {
   now, err := time.Parse(time.DateTime, "2024-05-13 01:02:03")
   if err != nil { t.Fatalf("Failed to parse time: %s", err) }

   meta := Meta{ID: "0Q1W2E3R4T5Y6U7I8O9P", Created: now}
   filename := meta.GenFileName()
   assert.Equal(t, "2024-05-13.01-02.0Q1W2E3R4T5Y6U7I8O9P.md", filename)
}

func Test_MarshalYAML(t *testing.T) {
   var err error
   t.Run("simple meta", func(t *testing.T) {
      meta := Meta {
         ID: "123456789",
         Tags: []string{"Foo", "Bar"},
      }
      meta.Created, err = time.Parse(time.RFC3339, "2024-05-13T01:02:03Z")
      if err != nil { t.Fatalf("Failed to parse time: %s", err) }

      data, err := yaml.Marshal(&meta)
      assert.Nil(t, err)
      assert.Equal(t, "created: \"2024-05-13 01:02:03\"\nid: \"123456789\"\ntags:\n    - Foo\n    - Bar\n", string(data))
   })

   t.Run("complex meta", func(t *testing.T) {
      title := "Test Note"
      format := "AsciiDoc"
      meta := Meta {
         ID: "123456789",
         Tags: []string{"Foo", "Bar"},
         Refs: &map[string]string{
            "Website": "https://example.com",
            "Book": "ISBN 1234567890",
         },
         Format: &format,
         Title: &title,
      }
      meta.Created, err = time.Parse(time.RFC3339, "2024-05-13T01:02:03Z")
      if err != nil { t.Fatalf("Failed to parse time: %s", err) }

      data, err := yaml.Marshal(&meta)
      assert.Nil(t, err)
      assert.Equal(t, "created: \"2024-05-13 01:02:03\"\nformat: AsciiDoc\nid: \"123456789\"\nrefs:\n    Book: ISBN 1234567890\n    Website: https://example.com\ntags:\n    - Foo\n    - Bar\ntitle: Test Note\n", string(data))
   })
}

func Test_marshalTime(t *testing.T) {
   now, err := time.Parse(time.DateTime, "2024-05-19 13:20:44")
   if err != nil {
      t.Fatalf("Failed to parse time: %s", err)
   }
   data, err := marshalTime(now)
   assert.Nil(t, err)
   assert.Equal(t, "2024-05-19 13:20:44", string(data))
   // TODO: add extra test cases after implementing config.Set()
}

func Test_UnmarshalYAML(t *testing.T) {
   t.Run("simple meta", func(t *testing.T) {
      data := "created: 2024-05-13 01:02:03\nid: \"123456789\"\ntags:\n    - Foo\n    - Bar\n"
      meta := Meta{}
      err := yaml.Unmarshal([]byte(data), &meta)
      assert.Nil(t, err)

      expected := Meta {
         ID: "123456789",
         Tags: []string{"Foo", "Bar"},
      }
      expected.Created, err = time.Parse(time.RFC3339, "2024-05-13T01:02:03Z")
      if err != nil { t.Fatalf("Failed to parse time: %s", err) }

      assert.Equal(t, expected, meta)
   })

   t.Run("complex meta", func(t *testing.T) {
      data := "created: 2024-05-13T01:02:03Z\nformat: AsciiDoc\nid: \"123456789\"\nrefs:\n    Book: ISBN 1234567890\n    Website: https://example.com\ntags:\n    - Foo\n    - Bar\ntitle: Test Note\n"
      meta := Meta{}
      err := yaml.Unmarshal([]byte(data), &meta)
      assert.Nil(t, err)
      title := "Test Note"
      format := "AsciiDoc"
      expected := Meta{
         ID: "123456789",
         Tags: []string{"Foo", "Bar"},
         Refs: &map[string]string{
            "Website": "https://example.com",
            "Book": "ISBN 1234567890",
         },
         Format: &format,

         Title: &title,
      }
      expected.Created, err = time.Parse(time.RFC3339, "2024-05-13T01:02:03Z")
      if err != nil { t.Fatalf("Failed to parse time: %s", err) }

      assert.Equal(t, expected, meta)
   })
}

