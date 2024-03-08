package marksmanager

import (
  "encoding/json"
  "os"
  "path/filepath"
  "sort"
)

// Match the JSON structure
type Marks struct {
  Marks []Mark `json:"marks"`
}

type Mark struct {
  Ref []int `json:"ref"`
  Bg string `json:"bg"`
  Fg string `json:"fg"`
  Ul bool `json:"ul"`
}

type MarksManager struct {
  BookPath string
  Chapter string
  Profiles map[string]Marks
}

func (m *MarksManager) Init() {
//  m.BookPath = ""
//  m.Chapter = ""
  m.Profiles = make(map[string]Marks)
}

func (m *MarksManager) Load(profile string) {
  file, err := os.ReadFile(filepath.Join(getUserHomeDir(), ".marks", profile, m.BookPath, m.Chapter+".json"))
  if err == nil {
    // Loading existing profile for this path
    var marks Marks
    err_1 := json.Unmarshal(file, &marks)
    if err_1 != nil {
      panic(err_1)
    }
    m.Profiles[profile] = marks
  } else {
    // Creating a new profile for this path
    m.Profiles[profile] = Marks{}
  }
}

func (m *MarksManager) LoadAll() {
  dirs, _ := os.ReadDir(filepath.Join(getUserHomeDir(), ".marks"))
  for _, d := range dirs {
    m.Load(d.Name())
  }
}

func (m *MarksManager) Add(profile string, mark Mark) {
  marks := m.Profiles[profile]
  marks.Marks = append(marks.Marks, mark)
  m.Profiles[profile] = marks
}

func (m *MarksManager) GetMergedMarks() Marks {
  result := Marks{}
  for _, marks := range m.Profiles {
    result.Marks = append(result.Marks, marks.Marks...)
  }
  return result
}

func (m *MarksManager) Store(profile string) {
  if _, err := os.Stat(m.BookPath); os.IsNotExist(err) {
    // Directory does not exist, create it
    err := os.MkdirAll(filepath.Join(getUserHomeDir(), ".marks", profile, m.BookPath), 0755) // 0755 is the directory permissions
    if err != nil {
      return
    }
  } else if err != nil {
    return
  }

  marks := m.Profiles[profile]

  // Sort the marks
  sort.Slice(marks.Marks, func(i, j int) bool {
    return marks.Marks[i].Ref[0] < marks.Marks[j].Ref[0]
  })
  // Store the Marks as JSON
  jsonFile, _ := json.MarshalIndent(marks, "", " ")
  _ = os.WriteFile(filepath.Join(getUserHomeDir(), ".marks", "default", m.BookPath, m.Chapter+".json"), jsonFile, 0644)
}

func (m *MarksManager) StoreAll() {
  for profile := range m.Profiles {
    m.Store(profile)
  }
}

func getUserHomeDir() string {
  if homeDir, err := os.UserHomeDir(); err == nil {
    return homeDir
  }
  // Fallback option
  return os.Getenv("HOME")
}

