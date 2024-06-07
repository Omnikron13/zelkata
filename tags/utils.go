package tags

import (
   "encoding/base32"
   "encoding/base64"
   "encoding/binary"
   "fmt"
   "strings"

   "github.com/omnikron13/zelkata/config"

   "github.com/cespare/xxhash"
)

// HashName generates a hash of the tag name to be used in e.g. filenames.
func HashName(name string) (hash string) {
   raw := hashNameRaw(name)
   hash = encodeHash(raw)
   if n := config.GetOrPanic[int]("tags.metadata.hash.truncate"); n > 0 {
      hash = hash[0:n]
   }
   return
}


// hashNameRaw generates a raw hash of the tag name to be used in e.g. filenames.
func hashNameRaw(name string) (raw uint64) {
   raw = xxhash.Sum64String(name)
   return
}


// encodeHash encodes a 64-bit hash as a string as specified in the config.
func encodeHash(raw uint64) string {
   bytes := make([]byte, 8)
   binary.LittleEndian.PutUint64(bytes, raw)

   format := config.GetOrPanic[string]("tags.metadata.hash.encode.format")
   switch format {
      case "base32":
         var encoding base32.Encoding
         charset := config.GetOrPanic[string]("tags.metadata.hash.encode.charset")
         switch charset {
            case "StdEncoding":
               encoding = *base32.StdEncoding
            case "HexEncoding":
               encoding = *base32.HexEncoding
            default:
               // base32 requires precisely 32 characters
               if len(charset) != 32 {
                  panic(fmt.Errorf("invalid encoding charset length: %d", len(charset)))
               }
               // charset must not contain newline/carriage return characters or duplicates
               for i, c := range charset {
                  if c == '\n' || c == '\r' {
                     panic(fmt.Errorf("invalid encoding charset character: %c", c))
                  }
                  if strings.ContainsRune(charset[:i], c) {
                     panic(fmt.Errorf("duplicate encoding charset character: %c", c))
                  }
               }
               encoding = *base32.NewEncoding(charset)
         }
         padChar := base32.NoPadding
         pad := config.GetOrPanic[bool]("tags.metadata.hash.encode.padding")
         if pad {
            padChar = base32.StdPadding
         }
         var sb strings.Builder
         encoder := base32.NewEncoder(encoding.WithPadding(padChar), &sb)
            defer encoder.Close()
         if _, err := encoder.Write(bytes); err != nil {
            panic("error encoding ID: " + err.Error())
         }
         return sb.String()

      case "base64":
         var encoding base64.Encoding
         charset := config.GetOrPanic[string]("tags.metadata.hash.encode.charset")
         switch charset {
            case "StdEncoding":
               encoding = *base64.StdEncoding
            case "URLEncoding":
               encoding = *base64.URLEncoding
            default:
               // base64 requires precisely 64 characters
               if len(charset) != 64 {
                  panic(fmt.Errorf("invalid encoding charset length: %d", len(charset)))
               }
               // charset must not contain newline/carriage return characters or duplicates
               for i, c := range charset {
                  if c == '\n' || c == '\r' {
                     panic(fmt.Errorf("invalid encoding charset character: %c", c))
                  }
                  if strings.ContainsRune(charset[:i], c) {
                     panic(fmt.Errorf("duplicate encoding charset character: %c", c))
                  }
               }
               encoding = *base64.NewEncoding(charset)
         }
         padChar := base64.NoPadding
         pad := config.GetOrPanic[bool]("tags.metadata.hash.encode.padding")
         if pad {
            padChar = base64.StdPadding
         }
         var sb strings.Builder
         encoder := base64.NewEncoder(encoding.WithPadding(padChar), &sb)
            defer encoder.Close()
         if _, err := encoder.Write(bytes); err != nil {
            panic("error encoding ID: " + err.Error())
         }
         return sb.String()

      default:
         panic(fmt.Errorf("unsupported encoding: %s", format))
   }
}

