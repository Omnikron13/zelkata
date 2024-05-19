package tag

import (
   "path/filepath"
   "testing"

   "github.com/stretchr/testify/assert"
   "gopkg.in/yaml.v3"
   "k8s.io/apimachinery/pkg/util/sets"
)


func Test_LoadPath(t *testing.T) {
   path := filepath.Join("testdata", "test.tag.yaml")
   tag, err := LoadPath(path)
   assert.Nil(t, err)
   assert.Equal(t, Tag{
      Name: "Test Tag",
      Description: "An example tag for testing purposes.",
      Icon: "󰓹",
      Notes: sets.New[string](
         "QWERTYUIOP",
         "ASDFGHJKLZ",
      ),
   }, *tag)
}


func Test_normaliseName(t *testing.T) {
   assert.Equal(t, "test-tag-name", normaliseName("Test TAG name"))
   assert.Equal(t, "already-normalised", normaliseName("already-normalised"))
 }


func Test_MarshalYAML(t *testing.T) {
   t.Run("simple tag", func(t *testing.T) {
      tag := Tag{
         Name: "Test Tag",
      }
      data, err := yaml.Marshal(tag)
      assert.Nil(t, err)
      assert.Equal(t, "name: Test Tag\nnotes: []\n", string(data))
   })

   t.Run("complex tag", func(t *testing.T) {
      tag := Tag{
         Name: "Test Tag",
         Description: "An example tag for testing purposes.",
         Virtual: true,
         Icon: "󰓹",
         Aliases: []string{
            "TestTag",
            "Test",
         },
         Notes: sets.New(
            "QWERTYUIOP",
            "ASDFGHJKLZ",
         ),
         Parents: map[string]string {
            "Parent 1": "parent-1",
            "Parent Number Two": "parent-number-two",
         },
         Relations: map[string]string {
            "Relation 1": "similar subject",
            "Relation Number Two": "first encountered in the same book",
         },
      }
      data, err := yaml.Marshal(tag)
      assert.Nil(t, err)
      assert.Equal(t, "aliases:\n    - TestTag\n    - Test\ndescription: An example tag for testing purposes.\nicon: \"\\U000F04F9\"\nname: Test Tag\nnotes:\n    - ASDFGHJKLZ\n    - QWERTYUIOP\nparents:\n    - Parent 1\n    - Parent Number Two\nrelations:\n    - name: Relation 1\n      description: similar subject\n    - name: Relation Number Two\n      description: first encountered in the same book\nvirtual: true\n", string(data))
   })
}

func Test_UnmarshalYAML(t *testing.T) {
   t.Run("simple tag", func(t *testing.T) {
      data := "name: Test Tag\nnotes: []\n"
      tag := Tag{}
      err := yaml.Unmarshal([]byte(data), &tag)
      assert.Nil(t, err)
      assert.Equal(t, Tag{
         Name: "Test Tag",
         Notes: sets.New[string](),
      }, tag)
   })

   t.Run("complex tag", func(t *testing.T) {
      data := "aliases:\n    - TestTag\n    - Test\ndescription: An example tag for testing purposes.\nicon: \"\\U000F04F9\"\nname: Test Tag\nnotes:\n    - QWERTYUIOP\n    - ASDFGHJKLZ\nparents:\n    - Parent 1\n    - Parent Number Two\nrelations:\n    - name: Relation 1\n      description: similar subject\n    - name: Relation Number Two\n      description: first encountered in the same book\nvirtual: true\n"
      tag := Tag{}
      err := yaml.Unmarshal([]byte(data), &tag)
      assert.Nil(t, err)
      expected := Tag{
         Name: "Test Tag",
         Description: "An example tag for testing purposes.",
         Virtual: true,
         Icon: "󰓹",
         Aliases: []string{
            "TestTag",
            "Test",
         },
         Notes: sets.New(
            "QWERTYUIOP",
            "ASDFGHJKLZ",
         ),
         Parents: map[string]string {
            "Parent 1": "parent-1",
            "Parent Number Two": "parent-number-two",
         },
         Relations: map[string]string {
            "Relation 1": "similar subject",
            "Relation Number Two": "first encountered in the same book",
         },
      }
      assert.Equal(t, expected, tag)
   })
}


// TODO: mock file access and test; Add(), LoadName(), Save()

