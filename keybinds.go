package main

import(
  "github.com/gdamore/tcell/v2"
)


func switchWithKeybinds(screen tcell.Screen,x,y int, buffer []string)(int,int,[]string){
  ev := screen.PollEvent()
  switch ev := ev.(type) {
  case *tcell.EventKey:
    switch ev.Key() {
    default:
      mod, key, ch, name := ev.Modifiers(), ev.Key(), ev.Rune(), ev.Name()
      _, _, _, _ = mod, key, ch, name
      switch name {
        case "Esc", "Ctrl+C":
           // Exit the program
        case "Up":
          x,y = up(x,y,buffer)			
        case "Down":
          x,y = down(x,y,buffer)
        case "Right":
          x,y = right(x,y,buffer)
        case "Left":
          x,y = left(x,y,buffer)
        case "Ctrl+S":
          write("hello.txt", buffer)
        case "Enter":
          x,y,buffer = insertNewLine(x,y,buffer)
          screen.Clear()
      case "Delete":
        x,y,buffer = delete(x,y,buffer)
        screen.Clear()
      case "Backspace2":
        x,y,buffer = backspace(x,y,buffer)
        screen.Clear()
      default:
        buffer[y] = insertChar(buffer[y], string(ch))
        x++
      }
    }
    //drawScreen(screen)
    //highlightSelection(screen, selectedIndex)
  case *tcell.EventResize:
    screen.Sync()
  }
  return x,y,buffer
}