# WorkerPool

## Описание проекта

WorkerPool - это реализация пула воркеров на Go с возможностью динамического добавления и удаления воркеров. Пул предназначен для параллельной обработки задач, представленных в виде строковых сообщений.

## Особенности

- Динамическое добавление и удаление воркеров
- Безопасная работа в конкурентной среде
- Гибкое управление через контекст
- Поддержка graceful shutdown
- Интеграция с системными сигналами

## Как это работает

### Основные компоненты

1. **WorkerPool** - центральная структура, управляющая пулом
2. **worker** - отдельный воркер, обрабатывающий сообщения
3. Каналы:
   - `messagesCh` - для передачи сообщений
   - `deletingCh` - для управления воркерами

### Жизненный цикл

1. Создается пул с указанным количеством воркеров
2. Воркеры ожидают сообщения в своих горутинах
3. При получении сообщения воркер обрабатывает его (выводит в stdout)
4. По команде или при завершении контекста воркеры останавливаются

## Интерфейсы

```go
// Конструкторы
func NewWorkerPull() *WorkerPool
func NewWorkerPullWithContext(ctx context.Context) *WorkerPool

// Основные методы
func (wp *WorkerPool) AddWorkers(count int) error
func (wp *WorkerPool) DeleteWorkers(count int) error
func (wp *WorkerPool) AddMessage(msg string) error
func (wp *WorkerPool) Close() error
func (wp *WorkerPool) GetWorkerCount() int
```

## Пример использования

```go
package main

import (
    workerPool "github.com/Alexandr-Fox/worker-pool/worker-pool"
)

func main() {
    wp := workerPool.NewWorkerPull()
    defer wp.Close()

    // Добавляем 3 воркера
    wp.AddWorkers(3)

    // Отправляем сообщения
    wp.AddMessage("Hello")
    wp.AddMessage("World")

    // Удаляем 1 воркера
    wp.DeleteWorkers(1)

    // Отправляем еще сообщения
    wp.AddMessage("Goodbye")

    // Закрываем WorkerPool и завершаем все воркеры
    wp.Close()
}
```

## Тестирование

Проект включает комплексные тесты:

1. **Базовые тесты**:
   - `TestWorkerPool_AddWorkers` - проверка добавления воркеров
   - `TestWorkerPool_DeleteWorkers` - проверка удаления воркеров
   - `TestWorkerPool_AddMessage` - проверка отправки сообщений

2. **Тесты корректности**:
   - `TestWorkerPool_Close` - проверка graceful shutdown
   - `TestWorkerPool_ContextCancel` - проверка отмены через контекст
   - `TestWorkerPool_AfterDeleteCallback` - проверка callback'ов

3. **Конкурентные тесты**:
   - `TestWorkerPool_ConcurrentOperations` - проверка работы при конкурентном доступе

Для запуска тестов:
```bash
go test -v ./worker-pool
```

## Интерактивный пример

В файле `main.go` представлен пример, который позволяет из консоли создавать воркеры, удалять воркеры, отправлять сообщения (для генерации сообщений используется lorem ipsum), а так же контролировать количество воркеров.

### Сборка

```bash
go build .
```

### Запуск

Если проект был предварительно собран:

```bash
./worker-pool
```

Если проект не был предварительно сборан (для разработчиков):

```bash
go run .
```