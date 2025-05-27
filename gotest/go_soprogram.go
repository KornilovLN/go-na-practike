package main

// go_soprogram.go

import (
  "fmt"
  "time"
)

func count() { // Функция, выполняемая как go-подпрограмма
  for i := 0; i < 15; i++ { 
    fmt.Println(i) 
    go text("text")   // Вызов go-подпрограмм 
    time.Sleep(time.Millisecond * 2)
  } 
}

func text(s string) { // Функция, выполняемая как go-подпрограмма
  for i := 0; i < 5; i++ { 
    fmt.Println(s) 
    time.Sleep(time.Millisecond * 1)
  } 
}

func main() {
  go count()         // Вызов go-подпрограмм
  go text("prompt")
  time.Sleep(time.Millisecond * 2)
  fmt.Println("Hello World")
  time.Sleep(time.Millisecond * 5)
}
