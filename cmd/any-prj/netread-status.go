package main

// Чтение состояния по протоколу TCP: netread_status.go

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

type ipaddr struct { // Структура для хранения адресов
	IP   string // IP-адрес
	name string // описание адреса
}

const ( // Константы для использования в коде
	HTTP_GET  = "GET / HTTP/1.0\r\n\r\n" // Строка запроса HTTP
	NUMB_ADDR = 6                        // Количество адресов в срезе
	CAP_ADDR  = NUMB_ADDR * 2            // Вместимость среза адресов
)

// Метод для форматирования вывода структуры ipaddr
func (a ipaddr) String() string {
	return fmt.Sprintf("IP: %s, Name: %s", a.IP, a.name)
}

// Функция для сканирования массива адресов и вывода их состояния
func ScanAddrArray(addr []ipaddr) {
	fmt.Println("=== ScanAddrArray: Сканирование адресов ===")

	for i := 0; i < len(addr); i++ {
		fmt.Printf("addr[%d] = %s\n", i, addr[i]) // Используем метод String()

		con, err := net.Dial("tcp", addr[i].IP) // Соединение по TCP
		if err != nil {
			fmt.Printf("Ошибка подключения к %s: %v\n", addr[i].IP, err)
			continue // Пропуск итерации в случае ошибки
		}
		defer con.Close() // Закрытие соединения

		fmt.Fprintf(con, HTTP_GET)                           // Отправка строки запроса
		status, err := bufio.NewReader(con).ReadString('\n') // Чтение ответа
		if err != nil {
			fmt.Printf("Ошибка чтения ответа от %s: %v\n", addr[i].IP, err)
			continue
		}
		fmt.Printf("Ответ: %s\n", status)
		fmt.Println("---")
	}
}

func main() {
	// срез адресов для контроля соединения
	addr := make([]ipaddr, NUMB_ADDR, CAP_ADDR)
	// Инициализация адресов
	addr[0] = ipaddr{IP: "localhost:8089", name: "Страница ссылок на сайты"}
	addr[1] = ipaddr{IP: "localhost:8082", name: "Документ: СПД на gitlab"}
	addr[2] = ipaddr{IP: "localhost:8083", name: "Документ: sunpp_comment"}
	addr[3] = ipaddr{IP: "192.168.88.102:8081", name: "Тесты на 102 машине по GO"}
	addr[4] = ipaddr{IP: "example.com:80", name: "example.com:80"}
	addr[5] = ipaddr{IP: "brama.sunpp.cns.atom:64080", name: "brama.sunpp.cns.atom:64080"}

	// Вызов функции сканирования
	go ScanAddrArray(addr)

	// Ждем завершения (простой способ)
	time.Sleep(10 * time.Second)
	// https://www.facebook.com/reel/606518561930101
}
