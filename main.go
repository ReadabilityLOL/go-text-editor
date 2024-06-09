/*
TODO:

- Read character
- Display function (Wrapper for tcell)
- Multiple buffers
- Lua API
- be able to open files in dir
- swap between windows
- move to root dir
- make a nicer exit function
*/

/*
DONE:
- Move left
- Move right
- Move up
- Move down
- Load files using buffio
- Write Files
- Close
- Add space after each line
*/

package main

import (
	//ioutil is deprecated, use io or os.
	"bufio"
	"fmt"
	"log"
	_ "math"
	"os"
	_ "path/filepath"
	_ "strings"
	"github.com/gdamore/tcell/v2"
	//"github.com/yuin/gopher-lua"
)

var currentDir string = "."

// var fileContent string = ""
var editorX, editorY, currX, currY, oldX int = 0, 0, 0, 0, 0 // Start the editor below the file list
var currFile = make([]string, 0)

func main() {
	// L := lua.NewState()
	// defer L.Close()
	// if err := L.DoFile(`init.lua`); err != nil {
	// 		panic(err)
	// }
	screen, err := tcell.NewScreen()
	cmd,err := tcell.NewScreen()
	_ = cmd
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := screen.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	defer screen.Fini()
	//displayLine(screen,0,0,tcell.StyleDefault,"hello, world")

	screen.Show()
	screen.Clear()
	var xmax, ymax = screen.Size()
	_, _ = xmax, ymax
	//displayfile(screen,"hello.txt")
	//drawScreen(screen)

	// Initial selection index
	//Handle user input
	buffer := load("hello.txt")
  mode := "normal"
	for {
    visibuffer := bufferize(buffer)
		drawTextEditor(screen, 0, 0, buffer, tcell.StyleDefault,1)
		currX,currY,buffer,mode = switchWithKeybinds(screen,currX,currY,buffer,visibuffer,mode)
	}
}

// func drawScreen(screen tcell.Screen) {
//     screen.Clear()

//     // List and display files in the current directory
//     files, err := listFiles(currentDir)
//     if err != nil {
//         log.Fatalf("Failed to list files: %+v", err)
//     }
//     drawFileList(screen, 1, 1, files, tcell.StyleDefault)

//     // Draw text editor content
//     drawTextEditor(screen, editorX, editorY, fileContent, tcell.StyleDefault)

//     screen.Show()
// }

// func getListLength(path string) int {
//     files, err := listFiles(path)
//     if err != nil {
//         return 0
//     }
//     return len(files) + 1 // +1 for the ".." entry
// }

// func drawFileList(s tcell.Screen, x, y int, files []os.DirEntry, style tcell.Style) {
//     // Draw the ".." entry for going up to the parent directory
//     displayLine(s, x, y, style, "..")
//     for i, file := range files {
//         name := file.Name()
//         displayLine(s, x, y+i+1, style, name)
//     }
//     s.Show()
// }

func splitLines(text string) []string {
	var lines []string
	currentLine := ""
	for _, r := range text {
		if r == '\n' {
			lines = append(lines, currentLine)
			currentLine = ""
		} else {
			currentLine += string(r)
		}
	}
	if currentLine != "" {
		lines = append(lines, currentLine)
	}
	return lines
}


// func highlightSelection(screen tcell.Screen, index int) {
// 	files, err := listFiles(currentDir)
// 	if err != nil {
// 		log.Fatalf("Failed to list files: %+v", err)
// 	}
// 	if index >= 0 && index <= len(files) {
// 		if index == 0 {
// 			// Highlight ".." entry
// 			displayLine(screen, 1, 1, tcell.StyleDefault.Background(tcell.ColorBlue).Foreground(tcell.ColorWhite), "..")
// 		} else {
// 			selected := files[index-1].Name()
// 			displayLine(screen, 1, index+1, tcell.StyleDefault.Background(tcell.ColorBlue).Foreground(tcell.ColorWhite), selected)
// 		}
// 	}
// 	screen.Show()
// }

// func highlight(screen tcell.Screen, text string, x, y int) {

// 	displayLine(screen, 0, y, tcell.StyleDefault, string(text[:x]))
// 	displayLine(screen, x, y, tcell.StyleDefault.Background(tcell.ColorBlue).Foreground(tcell.ColorWhite).Blink(true), string(text[x]))
// 	displayLine(screen, x+1, y, tcell.StyleDefault, string(text[x+1:]))
// }

// Draws text
func drawTextEditor(s tcell.Screen, x, y int, text []string, style tcell.Style,offset int) {
	for i, line := range text {
    s.ShowCursor(currX+offset,currY)	
		displayLine(s, x, y+i, style, line,offset)

	}
	s.Show()
}

// s file and returns text
func load(filename string) []string {
	// content, err := os.ReadFile(filename)
	// if err != nil {
	// 	return err.Error()
	// } else {
	// 	return insertChar(string(content),"e")
	// }
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	sc := bufio.NewScanner(file)
	lines := make([]string, 0)

	// Read through 'tokens' until an EOF is encountered.
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}

	if err := sc.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println(lines)
	return lines
}

// Displays string
func displayLine(s tcell.Screen, x, y int, style tcell.Style, str string,offset int) {
	for i, r := range str {
		s.SetContent(x+i+offset, y, r, nil, style)
	}
}

func insertChar(str string, char string) string {
	return str[:currX] + char + str[currX:]
}

func removeFrontChar(str string) string {
	return str[:currX] + str[currX+1:]
}

func removeBackChar(str string) string {
	return str[:currX-1] + str[currX:]
}

// Returns list of files
func listFiles(path string) ([]os.DirEntry, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func write(filename string, text []string) {
	written := ""
	for _, x := range text {
		written += x + "\n"
	}
	os.WriteFile(filename, []byte(written), 0644)
}

func bufferize(buffer []string)([]string){
  newbuff := make([]string,0)
  for _,x := range buffer{
    newbuff = append(newbuff,x+" ")
  }
  return newbuff
}