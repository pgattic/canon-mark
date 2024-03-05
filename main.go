package main

import (
  "bufio"
  "encoding/json"
  "fmt"
  "os"
  "path/filepath"
  "strconv"
  "strings"
)

type Flags struct { // command-line flags
  Paragraph bool // -p
  Verbose bool // -v
  VerseNumbers bool // -n
}

// Highlight struct to match the JSON structure
type Highlight struct {
  Ref []int `json:"ref"`
  Color string `json:"color"`
}

// Marks struct to hold an array of Highlight structs
type Marks struct {
  Highlights []Highlight `json:"highlights"`
}

var verse int
var chapter string
var bookPath string
var marks Marks
var execFlags Flags

func LoadMarks() Marks {
  file, err := os.ReadFile(filepath.Join(getUserHomeDir(), ".canon", "marks", "default", bookPath, chapter+".json"))
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

func setVerse(verseNum int) {
  verse = verseNum
}

func setChapter(chapStr string) {
  chapter = chapStr
  setVerse(0)
}

func setBookPath(bookStr string) {
  bookPath = bookStr
  setChapter("")
}

func printVerse(content string) {

  var appliedMarks []Highlight
  for i := 0; i < len(marks.Highlights); i++ {
    if (marks.Highlights[i].Ref[0] <= verse && marks.Highlights[i].Ref[1] >= verse) {
      appliedMarks = append(appliedMarks, marks.Highlights[i])
    }
  }
  if  execFlags.VerseNumbers {
    fmt.Print(verse, " ")
  }
  if len(appliedMarks) == 0 {
    fmt.Println(content)
    return
  }
  switch appliedMarks[len(appliedMarks)-1].Color {
  case "red":
    fmt.Println("\033[31m" + content + "\033[0m")
  case "green":
    fmt.Println("\033[32m" + content + "\033[0m")
  case "yellow":
    fmt.Println("\033[33m" + content + "\033[0m")
  case "blue":
    fmt.Println("\033[34m" + content + "\033[0m")
  case "magenta":
    fmt.Println("\033[35m" + content + "\033[0m")
  case "cyan":
    fmt.Println("\033[36m" + content + "\033[0m")
  case "white":
    fmt.Println("\033[37m" + content + "\033[0m")
  default:
    fmt.Println("\033[0m" + content)
  }
  //verseParts := strings.Split(content, " ")
}

func handleInput(text string) {
  if text[:3] == "@@@" {
    var verseParts = strings.Split(text[3:], " ")
    verse, _ = strconv.Atoi(verseParts[0])
    setVerse(verse)
    printVerse(text[len(verseParts[0])+4:])
  } else if text[:2] == "@@" {
    setChapter(text[2:])
    marks = LoadMarks()
  } else if text[0] == '@' {
    setBookPath(text[1:])
  }
}

func main() {
  args := os.Args

  if len(args) > 1 {

  //  var refIdx int // index of the args that is the verse index (flags could be before or after the verse ref)
    for i := len(args)-1; i >= 1; i-- {
      if args[i][0] == '-' {
        if args[i][1] == '-' { // args starting with "--"
          switch args[i] {
          case "--paragraph":
            execFlags.Paragraph = true
          case "--verbose":
            execFlags.Verbose = true
          case "--numbered":
            execFlags.VerseNumbers = true
          }
          continue
        }
        for ch := 1; ch < len(args[i]); ch++ {
          switch args[i][ch] { // Execution flags
          case 'p':
            execFlags.Paragraph = true
          case 'v':
            execFlags.Verbose = true
          case 'n':
            execFlags.VerseNumbers = true
          }
        }
//      } else {
//        refIdx = i
      }
    }
  }


  // Create a scanner to read from standard input
  scanner := bufio.NewScanner(os.Stdin)
  // Read input line by line
  for scanner.Scan() {
    line := scanner.Text()
    handleInput(line)
  }

  // Check for errors
  if err := scanner.Err(); err != nil {
    fmt.Fprintln(os.Stderr, "Error reading standard input:", err)
    os.Exit(1)
  }
}

func getUserHomeDir() string {
  if homeDir, err := os.UserHomeDir(); err == nil {
    return homeDir
  }
  // Fallback option if getting user's home directory fails
  return os.Getenv("HOME")
}

