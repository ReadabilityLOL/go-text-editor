/*
Roadmap:
- Move left
- Move right
- Move up
- Move down
- Read character
- Display function (Wrapper for tcell)
- Load
- Close
- Write
- Multiple buffers
- Lua API
- be able to open files in dir
- swap between windows
- move to root dir
*/

package main

import (
	//ioutil is deprecated, use io or os.
	_ "bufio"
	"github.com/gdamore/tcell/v2"
	"log"
	"os"
	"fmt"
	_ "path/filepath"
)

var currentDir string = "."
// var fileContent string = ""
var editorX, editorY, currX, currY int = 0, 0, 0, 0 // Start the editor below the file list

func main() {
	screen, err := tcell.NewScreen()
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
	_,_ = xmax,ymax
	//displayfile(screen,"hello.txt")
	//drawScreen(screen)

	// Initial selection index
	//selectedIndex := 0
	buffer := load("hello.txt")
	//highlight(screen,buffer,0,0)
	//Handle user input
	for {
		drawTextEditor(screen, 0, 0, buffer, tcell.StyleDefault)
		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEscape, tcell.KeyCtrlC:
				return // Exit the program
			case tcell.KeyEnter:
				// // Enter the selected directory or open file
				// enterSelection(screen, selectedIndex)
				// selectedIndex = 0
			case tcell.KeyUp:
				//Move selection up
				if currY > 0 {
					currY--
				}
			case tcell.KeyDown:
				currY += 1
				// Move selection down
				// if currY < g(currentDir)-1 {
				//     selectedIndex++
				// }

			case tcell.KeyRight:
				currX++

			case tcell.KeyLeft:
				currX--
			default:
				mod, key, ch := ev.Modifiers(), ev.Key(), ev.Rune()
				fmt.Printf("EventKey Modifiers: %d Key: %d Rune: %s", mod, key, string(ch))
			}
			//drawScreen(screen)
			//highlightSelection(screen, selectedIndex)
		case *tcell.EventResize:
			screen.Sync()
      xmax, ymax = screen.Size()
		}
		//screen.Show()
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

// func enterSelection(screen tcell.Screen, index int) {
//     if index == 0 {
//         // ".." selected, go up to parent directory
//         currentDir = filepath.Dir(currentDir)
//         fileContent = ""
//     } else {
//         files, err := listFiles(currentDir)
//         if err != nil {
//             log.Fatalf("Failed to list files: %+v", err)
//         }
//         if index > 0 && index <= len(files) {
//             selected := files[index-1]
//             if selected.IsDir() {
//                 // Enter directory
//                 currentDir = filepath.Join(currentDir, selected.Name())
//                 fileContent = ""
//                 drawScreen(screen)
//             } else {
//                 // Open file
//                 filePath := filepath.Join(currentDir, selected.Name())
//                 fileContent = load(filePath)
//                 //displayfile(screen,filePath)
//                 // if err != nil {
//                 //     fileContent = "Failed to read file: " + err.Error()
//                 // } else {
//                 //     fileContent = string(content)
//                 // }
//             }
//         }
//     }
// }

func highlightSelection(screen tcell.Screen, index int) {
	files, err := listFiles(currentDir)
	if err != nil {
		log.Fatalf("Failed to list files: %+v", err)
	}
	if index >= 0 && index <= len(files) {
		if index == 0 {
			// Highlight ".." entry
			displayLine(screen, 1, 1, tcell.StyleDefault.Background(tcell.ColorBlue).Foreground(tcell.ColorWhite), "..")
		} else {
			selected := files[index-1].Name()
			displayLine(screen, 1, index+1, tcell.StyleDefault.Background(tcell.ColorBlue).Foreground(tcell.ColorWhite), selected)
		}
	}
	screen.Show()
}

func highlight(screen tcell.Screen, text string, x, y int) {
    if x >= len(text){
        x = len(text) - 1
        currX = len(text) - 1
    }
    if x < 0 {
        x = 0
        currX = 0
    }
	displayLine(screen, 0, y, tcell.StyleDefault, string(text[:x]))
	displayLine(screen, x, y, tcell.StyleDefault.Background(tcell.ColorBlue).Foreground(tcell.ColorWhite), string(text[x]))
	displayLine(screen, x+1, y, tcell.StyleDefault, string(text[x+1:]))
}

// Draws text
func drawTextEditor(s tcell.Screen, x, y int, text string, style tcell.Style) {
	// text := load(filename)
	lines := splitLines(text)
	for i, line := range lines {
		if i == currY {
			highlight(s, line, currX, y+i)
		} else {
			displayLine(s, x, y+i, style, line)
		}

	}
	s.Show()
}

// Loads file and returns text
func load(filename string) string {
	content, err := os.ReadFile(filename)
	if err != nil {
		return err.Error()
	} else {
		return insertChar(string(content),"e")
	}
}

// Displays string
func displayLine(s tcell.Screen, x, y int, style tcell.Style, str string) {
	for i, r := range str {
		s.SetContent(x+i, y, r, nil, style)
	}
}

func insertChar(str string, char string) string{
	return str[:currX]+char+str[currX:]
}

//Loads and displays file
// func displayfile(screen tcell.Screen, filename string){
//     screen.Clear()
//     fileText := load(filename)
//     displayLine(screen,1,1,tcell.StyleDefault,fileText)
//     screen.Sync()
// }

// Returns list of files
func listFiles(path string) ([]os.DirEntry, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	return files, nil
}
