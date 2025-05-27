package main

// Чтение состояния по протоколу TCP: read_status.go

import (
  "bufio"
  "fmt"
  "net"
)

func main() {
  conn, _ := net.Dial("tcp", "localhost:8089") // Соединение по TCP
  fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")   // Отправка строки через соединение
  status, _ := bufio.NewReader(conn).ReadString('\n') // Вывод первой строки ответа
  fmt.Println(status)
}
