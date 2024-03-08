package marksmanager

import (
  "encoding/json"
  "os"
  "path/filepath"
  "sort"
)

// Marks struct to hold an array of Highlight structs
type Marks struct {
  Marks []Mark `json:"marks"`
}

// Mark struct to match the JSON structure
type Mark struct {
  Ref []int `json:"ref"`
  Bg string `json:"bg"`
  Fg string `json:"fg"`
  Ul bool `json:"ul"`
}


func LoadMarks(bookPath string, chapter string) Marks {
  file, err := os.ReadFile(filepath.Join(getUserHomeDir(), ".marks", "default", bookPath, chapter+".json"))
  if err != nil {
    return Marks{}
  }

  // Unmarshal JSON into Config struct
  var marks Marks
  err_1 := json.Unmarshal(file, &marks)
  if err_1 != nil {
    panic(err_1)
  }
  return marks
}

func StoreMarks(marks Marks, bookPath string, chapter string) {
  if _, err := os.Stat(bookPath); os.IsNotExist(err) {
    // Directory does not exist, create it
    err := os.MkdirAll(filepath.Join(getUserHomeDir(), ".marks", "default", bookPath), 0755) // 0755 is the directory permissions
    if err != nil {
      return
    }
  } else if err != nil {
    return
  }



  // Sort the marks
  sort.Slice(marks.Marks, func(i, j int) bool {
    return marks.Marks[i].Ref[0] < marks.Marks[j].Ref[0]
  })
  // Store the Marks as JSON
  jsonFile, _ := json.MarshalIndent(marks, "", " ")
  _ = os.WriteFile(filepath.Join(getUserHomeDir(), ".marks", "default", bookPath, chapter+".json"), jsonFile, 0644)
}

func getUserHomeDir() string {
  if homeDir, err := os.UserHomeDir(); err == nil {
    return homeDir
  }
  // Fallback option
  return os.Getenv("HOME")
}

