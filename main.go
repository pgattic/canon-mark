package main

import (
  "bufio"
  "fmt"
  "os"
  "strconv"
  "strings"
  "github.com/pgattic/marks/marksmanager"
)

type Flags struct { // command-line flags
  Paragraph bool // -p
  Verbose bool // -v
  VerseNumbers bool // -n
  adding bool
}

var verse int
var execFlags Flags
var markToAdd marksmanager.Mark

var marksManager marksmanager.MarksManager

var vRangeStarted bool = false

func setVerse(verseNum int) {
  verse = verseNum
  if execFlags.adding {
    if !vRangeStarted {
      vRangeStarted = true
      markToAdd.Ref = []int{verse, verse}
    } else {
      markToAdd.Ref[1] = verse
    }
  }
}

func setChapter(chapStr string) {
  if execFlags.adding && vRangeStarted {
    marksManager.Add("default", markToAdd)
    marksManager.Store("default")
    markToAdd = marksmanager.Mark{}
    vRangeStarted = false
  }
  marksManager.Init()
  marksManager.Chapter = chapStr
  marksManager.Load("default")
}

func setBookPath(bookStr string) {
  setChapter("")
  marksManager.BookPath = bookStr
}

func resolveBgColor(color string) string { // Get ANSI code for background color
  switch color {
  case "red":
    return "41"
  case "green":
    return "42"
  case "yellow":
    return "43"
  case "blue":
    return "44"
  case "magenta":
    return "45"
  case "cyan":
    return "46"
  case "white":
    return "47"
  default:
    return ""
  }
}

func resolveFgColor(color string) string { // Get ANSI code for foreground color
  switch color {
  case "red":
    return "91"
  case "green":
    return "92"
  case "yellow":
    return "93"
  case "blue":
    return "94"
  case "magenta":
    return "95"
  case "cyan":
    return "96"
  case "white":
    return "97"
  default:
    return ""
  }
}

func printVerse(content string) {
  var fgCol string
  var bgCol string
  var ul bool
  
  if execFlags.adding {
      if markToAdd.Bg != "" {
        bgCol = markToAdd.Bg
      }
      if markToAdd.Fg != "" {
        fgCol = markToAdd.Fg
      }
      ul = markToAdd.Ul
  }

  for i := 0; i < len(marksManager.Profiles["default"].Marks); i++ {
    if (marksManager.Profiles["default"].Marks[i].Ref[0] <= verse && marksManager.Profiles["default"].Marks[i].Ref[1] >= verse) {
      if marksManager.Profiles["default"].Marks[i].Bg != "" {
        bgCol = marksManager.Profiles["default"].Marks[i].Bg
      }
      if marksManager.Profiles["default"].Marks[i].Fg != "" {
        fgCol = marksManager.Profiles["default"].Marks[i].Fg
      }
      ul = marksManager.Profiles["default"].Marks[i].Ul
    }
  }
  if  execFlags.VerseNumbers {
    fmt.Print(" \033[1m", verse, "\033[0m ")
  }
  ANSICode := "\033["
  if fgCol != "" {
    ANSICode += ";" + resolveFgColor(fgCol)
  }
  if bgCol != "" {
    ANSICode += ";" + resolveBgColor(bgCol)
  }
  if ul {
    ANSICode += ";4"
  }
  ANSICode += "m"
  fmt.Println(ANSICode + content + "\033[0m")
  if execFlags.Paragraph {
    fmt.Println()
  }
}

func handleInput(text string) {
  if text[:3] == "@@@" {
    var verseParts = strings.Split(text[3:], " ")
    verse, _ = strconv.Atoi(verseParts[0])
    setVerse(verse)
    printVerse(text[len(verseParts[0])+4:])
  } else if text[:2] == "@@" {
    setChapter(text[2:])
    marksManager.Load("default")
  } else if text[0] == '@' {
    setBookPath(text[1:])
  }
}

func main() {
  args := os.Args

  execFlags = Flags{false, false, false, false}

  if len(args) > 1 {
    switch args[1] {
    case "add":
      execFlags.adding = true
    }

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
//          if len(args[i]) > 10 && args[i][:11] == "--profiles=" {
//            profiles = strings.Split(args[i][11:], ",")
//          }
          if len(args[i]) > 5 && args[i][:5] == "--bg=" {
            markToAdd.Bg = args[i][5:]
          } else if len(args[i]) > 5 && args[i][:5] == "--fg=" {
            markToAdd.Fg = args[i][5:]
          } else if args[i] == "--ul" {
            markToAdd.Ul = true
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
  marksManager.Profiles = make(map[string]marksmanager.Marks)

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

  if execFlags.adding && vRangeStarted {
    marksManager.Add("default", markToAdd)
    marksManager.Store("default")
  }
}


