package main

// Использование каналов: go_chan.go

import (
  "fmt"
  "time"
)

func printCount(c chan int) { // аргумент - канал для передачи целого значения
  num := 0
  for num >= 0 {
    num = <-c            // Ожидание целого значения
    fmt.Print(num, " ")
  }
}

// В канал c будут отправляться числа массива a
func getValue(c chan int, arr [] int) {
  for _, v := range arr {  // Запись из массива целого значения в канал
    c <- v
  }
}

func main() {
  c := make(chan int)    // Создание канала
  a := []int{8, 6, 7, 5, 3, 0, 9, -1} // Тестовый массив целых

  go printCount(c)       // Вызов сопрограммы которая печатает приход из канала
  go getValue(c, a)      // Вызов сопрограммы которая отправляет числа в канал

  time.Sleep(time.Millisecond * 1)  // Функция main приостанавливается
  fmt.Println("End of main")        // перед завершением
}
