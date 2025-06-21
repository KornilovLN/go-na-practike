package main

// Пример веб-приложения с использованием переменных окружения: env_config.go
// Этот пример демонстрирует, как использовать переменные окружения для настройки веб-сервера на Go.
// Установите переменную окружения PORT и запустите сервер
// export PORT=8181  # Для Linux
// после этого запустите приложение командой:
// go run cmd/app-12factor/env_config.go
// При запросе в браузере: http://localhost:8181/ → вернёт "The homepage."
// Можно и так сделать запрос: curl http://localhost:8181/

import (
	"fmt"
	"net/http"
	"os"
)

// Проверка наличия переменной окружения PORT
func main() {
	// Роутинг: связываем корневой URL с функцией homePage
	http.HandleFunc("/", homePage)

	// Извлечение  значения переменной PORT из окружения
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}

// homePage обрабатывает запросы к корневому URL и возвращает текст "The homepage."
func homePage(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.NotFound(res, req)
		return
	}
	fmt.Fprint(res, "The homepage.")
}
