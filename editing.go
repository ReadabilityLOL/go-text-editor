package main

import(
  "strings"
)

//The following function is kinda stupid
func setpos(x,y int)(int,int){
 return x,y 
}

func left(x,y int, buffer []string) (int, int) {
  if x > 0 {
    x--
  }
  return x,y
}

func right(x,y int, buffer []string) (int, int) {
  if x < len(buffer[y])-1 {
    x++
  }
  return x,y
}

func up(x,y int, buffer []string) (int, int){
  if y > 0 {
    y--
    if x >= len(buffer[y])-1 {
      x = len(buffer[y]) - 1
    }
  }
  return x,y
}

func down(x,y int, buffer []string)(int,int){
  //Move selection down
  if y < len(buffer)-1 {
    y++
    if x >= len(buffer[y])-1 {
      x = len(buffer[y]) - 1
    }
  }
  return x,y
}



func insertNewLine(x,y int, buffer []string) (int,int,[]string){
  newBuff := buffer[y][:x]
  afterBuff := buffer[y][x:]
  buffer = append(buffer[:y+1], buffer[y:]...)
  buffer[y] = newBuff
  buffer[y+1] = afterBuff
  x = 0
  y++
  return x, y, buffer
}


func backspace(x,y int, buffer []string) (int,int,[]string){
  if x > 0 {
    buffer[y] = removeBackChar(buffer[y])
    x--
  } else if y > 0 {
    for x := range buffer {
      buffer[x] = strings.TrimRight(buffer[x], " ") + " "
    }
    x = (len(buffer[y-1]) - 1)
    buffer[y-1] = buffer[y-1] + buffer[y]
    buffer = append(buffer[:y], buffer[y+1:]...)
    y--
  }
  return x,y,buffer
}

func delete(x,y int, buffer []string) (int,int,[]string){
  if x < len(buffer[y])-1 {
    buffer[y] = removeFrontChar(buffer[y])
  } 

  return x,y,buffer
}