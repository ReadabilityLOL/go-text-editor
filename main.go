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
	"github.com/gdamore/tcell/v2"
	"log"
	_ "math"
	"os"
	_ "path/filepath"
	"strings"
)

var currentDir string = "."

// var fileContent string = ""
var editorX, editorY, currX, currY, oldX int = 0, 0, 0, 0, 0 // Start the editor below the file list
var currFile = make([]string, 0)

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
	_, _ = xmax, ymax
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
			case tcell.KeyUp:
				currX,currY = up(currX,currY,buffer)			
			case tcell.KeyDown:
				currX,currY = down(currX,currY,buffer)
			case tcell.KeyRight:
				currX,currY = right(currX,currY,buffer)
			case tcell.KeyLeft:
				currX,currY = left(currX,currY,buffer)
				
			case tcell.KeyCtrlS:
				write("hello.txt", buffer)
			case tcell.KeyEnter:
				newBuff := buffer[currY][:currX]
				afterBuff := buffer[currY][currX:]
				buffer = append(buffer[:currY+1], buffer[currY:]...)
				buffer[currY] = newBuff
				buffer[currY+1] = afterBuff
				currX = 0
				currY++

				screen.Clear()
			default:
				mod, key, ch, name := ev.Modifiers(), ev.Key(), ev.Rune(), ev.Name()
				_, _, _, _ = mod, key, ch, name
				switch name {
				case "Delete":
					if currX < len(buffer[currY])-1 {
						buffer[currY] = removeFrontChar(buffer[currY])
						screen.Clear()
					}
				case "Backspace2":

					if currX > 0 {
						buffer[currY] = removeBackChar(buffer[currY])
						currX--
						screen.Clear()
					} else if currY > 0 {
						for x := range buffer {
							buffer[x] = strings.TrimRight(buffer[x], " ") + " "
						}
						currX = (len(buffer[currY-1]) - 1)
						buffer[currY-1] = buffer[currY-1] + buffer[currY]
						buffer = append(buffer[:currY], buffer[currY+1:]...)
						currY--
						screen.Clear()
					}
				default:
					buffer[currY] = insertChar(buffer[currY], string(ch))
					currX++
				}
			}
			//drawScreen(screen)
			//highlightSelection(screen, selectedIndex)
		case *tcell.EventResize:
			screen.Sync()
			xmax, ymax = screen.Size()
		}
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

	displayLine(screen, 0, y, tcell.StyleDefault, string(text[:x]))
	displayLine(screen, x, y, tcell.StyleDefault.Background(tcell.ColorBlue).Foreground(tcell.ColorWhite).Blink(true), string(text[x]))
	displayLine(screen, x+1, y, tcell.StyleDefault, string(text[x+1:]))
}

// Draws text
func drawTextEditor(s tcell.Screen, x, y int, text []string, style tcell.Style) {
	// text := load(filename)
	//lines := splitLines(text)
	for i, line := range text {
		if i == currY {
			highlight(s, line, currX, y+i)
		} else {
			displayLine(s, x, y+i, style, line)
		}

	}
	s.Show()
}

// Loads file and returns text
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
		lines = append(lines, sc.Text()+" ")
	}

	if err := sc.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println(lines)
	return lines
}

// Displays string
func displayLine(s tcell.Screen, x, y int, style tcell.Style, str string) {
	for i, r := range str {
		s.SetContent(x+i+1, y, r, nil, style)
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

func write(filename string, text []string) {
	written := ""
	for _, x := range text {
		written += x + "\n"
	}
	os.WriteFile(filename, []byte(written), 0644)
}

// func parseInput(event tcell.Event){
// 	switch ev := event.(type) {
// 		cas
// 	}
// }
