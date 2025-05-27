package main



import (
  "fmt"
  "net/http"
)

// Обработка HTTP-запроса
func hello(res http.ResponseWriter, req *http.Request) { 
  fmt.Fprint(res, "Hello, my name is LN Starmark")  
}

// Основная логика приложения
func main() {
  http.HandleFunc("/", hello)
  http.ListenAndServe("localhost:8080", nil)
}


