// Package paths provides very light utility functions for getting path strings to the various data directories etc.
package paths

import (
   "os"
   "path/filepath"

   "github.com/omnikron13/zelkata/config"

   "github.com/adrg/xdg"
)

// Cached path strings
var dataDir string
var noteDir string
var tagDir  string
var stateDir string


// Data returns the path to the root data directory that Zelkata is to use.
func Data() string {
   if dataDir != "" {
      return dataDir
   }

   // Get the path to the notes directory.
   path, err := config.Get[string]("data-directory")
   if err != nil {
      panic(err)
   }

   // filepath sadly has no convenience function for expanding environment variables...
   path = os.ExpandEnv(path)

   // Converting to an absolute path may not really be worth the time, but we don't want any confusing behaviour
   // cropping up down the line.
   if dataDir, err = filepath.Abs(path); err != nil {
      panic(err)
   }

   // Ensure the directory exists, creating it if necessary (including all directories along the way)
   if err := os.MkdirAll(dataDir, 0700); err != nil {
      panic(err)
   }

   return dataDir
}


// Notes returns the path to the notes directory.
func Notes() string {
   if noteDir != "" {
      return noteDir
   }
   noteDir = filepath.Join(Data(), "notes")
   if err := os.MkdirAll(noteDir, 0700); err != nil {
      panic(err)
   }
   return noteDir
}


// Tags returns the path to the tags directory.
func Tags() string {
   if tagDir != "" {
      return tagDir
   }
   tagDir = filepath.Join(Data(), "tags")
   if err := os.MkdirAll(tagDir, 0700); err != nil {
      panic(err)
   }
   return tagDir
}


// State returns the path to the state directory.
func State() string {
   if stateDir != "" {
      return stateDir
   }
   stateDir := filepath.Join(xdg.StateHome, "zelkata")
   if err := os.MkdirAll(stateDir, 0700); err != nil {
      panic(err)
   }
   return stateDir
}

