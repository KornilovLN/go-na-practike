package main

// HTTP-запрос GET: http_get.go

import (
  "fmt"
  "io/ioutil"
  "net/http"
)

func main() {
  resp, _ := http.Get("http://example.com/")  // Создание HTTP-запроса GET
  body, _ := ioutil.ReadAll(resp.Body)        // Чтение тела ответа
  fmt.Println(string(body))                   // Вывод тела в виде строки
  resp.Body.Close()                           // Закрытие соединения
}
