package tags

import (
   "testing"

   "github.com/stretchr/testify/assert"
)


func Test_HashName(t *testing.T) {
   name := "Test Tag"
   hash := HashName(name)
   assert.Equal(t, "ZR73CSRY", hash)
}


func Test_hashNameRaw(t *testing.T) {
   name := "Test Tag"
   hash := hashNameRaw(name)
   assert.Equal(t, uint64(0xd2fca0384ab17fcc), hash)
}

