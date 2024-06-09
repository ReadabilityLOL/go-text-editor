package main

import (
	"github.com/gdamore/tcell/v2"
)

func eventSwitch(event tcell.Event) (string,string)  {
  switch ev := event.(type){
  case *tcell.EventKey:
    mod, key, ch, name := ev.Modifiers(), ev.Key(), ev.Rune(), ev.Name()
    _, _, _, _ = mod, key, ch, name
    return name,string(ch)
  case *tcell.EventResize:
    return "resize","c"
  }
  return "resize","c"
}


func switchWithKeybinds(screen tcell.Screen,x,y int, truebuffer, buffer []string, mode string)(int,int,[]string,string){
  switch mode{
  case "insert":
    return insertMode(screen,x,y,truebuffer,buffer)
  case "normal":
    return normalMode(screen,x,y,truebuffer,buffer)
  }
  // ev := screen.PollEvent()
  // name,ch := eventSwitch(ev)
  // switch name {
  //   case "Esc", "Ctrl+C":
  //      // Exit the program
  //   case "Up":
  //     x,y = up(x,y,buffer)			
  //   case "Down":
  //     x,y = down(x,y,buffer)
  //   case "Right":
  //     x,y = right(x,y,buffer)
  //   case "Left":
  //     x,y = left(x,y,buffer)
  //   case "Ctrl+S":
  //     write("hello.txt", truebuffer)
  //   case "Enter":
  //     x,y,buffer = insertNewLine(x,y,truebuffer)
  //     screen.Clear()
  // case "Delete":
  //   x,y,buffer = delete(x,y,truebuffer)
  //   screen.Clear()
  // case "Backspace2":
  //   x,y,buffer = backspace(x,y,truebuffer)
  //   screen.Clear()
  // case "resize":
  //   screen.Sync()
  // default:
  //   truebuffer[y] = insertChar(truebuffer[y], string(ch))
  //   x++
  // }  
  // return x,y,truebuffer
  return x,y,truebuffer,mode
}

func insertMode(screen tcell.Screen,x,y int, truebuffer, buffer []string)(int,int,[]string,string){
  
  screen.SetCursorStyle(5)
  ev := screen.PollEvent()
  name,ch := eventSwitch(ev)
  mode := "insert"
  switch name {
    case "Esc","Ctrl+O":
      mode = "normal"
      screen.SetCursorStyle(0)
    case "Up":
      x,y = up(x,y,buffer)			
    case "Down":
      x,y = down(x,y,buffer)
    case "Right":
      x,y = right(x,y,buffer)
    case "Left":
      x,y = left(x,y,buffer)
    case "Ctrl+S":
      write("hello.txt", truebuffer)
    case "Enter":
      x,y,buffer = insertNewLine(x,y,truebuffer)
      screen.Clear()
  case "Delete":
    x,y,buffer = delete(x,y,truebuffer)
    screen.Clear()
  case "Backspace2":
    x,y,buffer = backspace(x,y,truebuffer)
    screen.Clear()
  case "resize":
    screen.Sync()
  default:
    truebuffer[y] = insertChar(truebuffer[y], string(ch))
    x++
  }  
  return x,y,truebuffer,mode
}

func normalMode(screen tcell.Screen,x,y int, truebuffer, buffer []string)(int,int,[]string,string){

  screen.SetCursorStyle(1)
  ev := screen.PollEvent()
  name,ch := eventSwitch(ev)
  mode := "normal"
  switch name {
    case "Up":
      x,y = up(x,y,buffer)			
    case "Down":
      x,y = down(x,y,buffer)
    case "Right":
      x,y = right(x,y,buffer)
    case "Left":
      x,y = left(x,y,buffer)
    case "Ctrl+S":
      write("hello.txt", truebuffer)
    case "Enter":
      x,y,buffer = insertNewLine(x,y,truebuffer)
      screen.Clear()
  case "Delete":
    x,y,buffer = delete(x,y,truebuffer)
    screen.Clear()
  case "Backspace2":
    x,y,buffer = backspace(x,y,truebuffer)
    screen.Clear()
  case "resize":
    screen.Sync()
  default:
    switch ch{
      case "i":
        mode = "insert"
        screen.SetCursorStyle(5)
      case "j":
        x,y = up(x,y,buffer)			
      case "k":
        x,y = down(x,y,buffer)
      case "l":
        x,y = right(x,y,buffer)
      case "h":
        x,y = left(x,y,buffer)
      case "r":
        
      }  
    }

  return x,y,truebuffer,mode
}