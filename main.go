package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	workerPool "github.com/Alexandr-Fox/worker-pool/worker-pool"
	"github.com/go-loremipsum/loremipsum"
)

func main() {
	// Инициализация генератора случайного текста
	loremIpsumGenerator := loremipsum.NewWithSeed(1234)

	// Создаем пул воркеров
	wp := workerPool.NewWorkerPull()
	defer wp.Close()

	// Канал для сигналов ОС
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Горутина для обработки команд пользователя
	go func() {
		for {
			fmt.Println("\nВыберите действие:")
			fmt.Println("1 - Добавить воркеров")
			fmt.Println("2 - Удалить воркеров")
			fmt.Println("3 - Отправить сообщения")
			fmt.Println("4 - Статус пула")
			fmt.Println("5 - Выход")

			var choice int
			fmt.Scan(&choice)

			switch choice {
			case 1:
				fmt.Print("Сколько воркеров добавить? ")
				var count int
				fmt.Scan(&count)
				if err := wp.AddWorkers(count); err != nil {
					fmt.Printf("Ошибка добавления воркеров: %v\n", err)
				} else {
					fmt.Printf("Добавлено %d воркеров\n", count)
				}

			case 2:
				fmt.Print("Сколько воркеров удалить? ")
				var count int
				fmt.Scan(&count)
				if err := wp.DeleteWorkers(count); err != nil {
					fmt.Printf("Ошибка удаления воркеров: %v\n", err)
				} else {
					fmt.Printf("Удалено %d воркеров\n", count)
				}

			case 3:
				fmt.Print("Сколько сообщений отправить? ")
				var count int
				fmt.Scan(&count)

				var wg sync.WaitGroup
				wg.Add(count)

				for i := 0; i < count; i++ {
					go func(id int) {
						defer wg.Done()
						msg := fmt.Sprintf("Сообщение %d: %s", id, loremIpsumGenerator.Sentence())
						if err := wp.AddMessage(msg); err != nil {
							fmt.Printf("Ошибка отправки сообщения: %v\n", err)
						}
					}(i)
				}

				wg.Wait()
				fmt.Printf("Отправлено %d сообщений\n", count)

			case 4:
				fmt.Printf("Текущее количество воркеров: %d\n", wp.GetWorkers())

			case 5:
				sigChan <- syscall.SIGTERM
				return

			default:
				fmt.Println("Неверный выбор, попробуйте снова")
			}
		}
	}()

	// Ожидаем сигнал завершения
	<-sigChan
	fmt.Println("\nПолучен сигнал завершения, закрываем пул...")
}
