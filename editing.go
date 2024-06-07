package main


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
    if x >= len(buffer[currY])-1 {
      x = len(buffer[currY]) - 1
    }
  }
  return x,y
}

func down(x,y int,buffer []string)(int,int){
  //Move selection down
  if y < len(buffer)-1 {
    y++
    if x >= len(buffer[y])-1 {
      x = len(buffer[y]) - 1
    }
  }
  return x,y
}