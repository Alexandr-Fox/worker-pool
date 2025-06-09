package workerPool

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestWorkerPool_AddWorkers(t *testing.T) {
	wp := NewWorkerPull()
	defer wp.Close()

	err := wp.AddWorkers(3)
	if err != nil {
		t.Fatalf("Ошибка при добавлении воркеров: %v", err)
	}

	if wp.workers != 3 {
		t.Errorf("Ожидалось 3 воркера, получено %d", wp.workers)
	}
}

func TestWorkerPool_DeleteWorkers(t *testing.T) {
	wp := NewWorkerPull()
	defer wp.Close()

	wp.AddWorkers(5)

	err := wp.DeleteWorkers(2)
	if err != nil {
		t.Fatalf("Ошибка при удалении воркеров: %v", err)
	}

	if wp.workers != 3 {
		t.Errorf("Ожидалось 3 воркера после удаления, получено %d", wp.workers)
	}
}

func TestWorkerPool_AddMessage(t *testing.T) {
	wp := NewWorkerPull()
	defer wp.Close()

	wp.AddWorkers(1)

	err := wp.AddMessage("test message")
	if err != nil {
		t.Fatalf("Ошибка при отправке сообщения: %v", err)
	}
}

func TestWorkerPool_Close(t *testing.T) {
	wp := NewWorkerPull()

	wp.AddWorkers(2)
	err := wp.Close()
	if err != nil {
		t.Fatalf("Ошибка при закрытии пула: %v", err)
	}

	if !wp.closed {
		t.Error("Пул не помечен как закрытый после Close()")
	}

	// Проверяем, что после закрытия операции возвращают ошибки
	err = wp.AddWorkers(1)
	if err == nil {
		t.Error("Ожидалась ошибка при добавлении воркеров в закрытый пул")
	}

	err = wp.DeleteWorkers(1)
	if err == nil {
		t.Error("Ожидалась ошибка при удалении воркеров из закрытого пула")
	}

	err = wp.AddMessage("test")
	if err == nil {
		t.Error("Ожидалась ошибка при отправке сообщения в закрытый пул")
	}
}

func TestWorkerPool_ConcurrentOperations(t *testing.T) {
	wp := NewWorkerPull()
	defer wp.Close()

	var wg sync.WaitGroup
	wg.Add(4)

	// Конкурентные операции
	go func() {
		defer wg.Done()
		wp.AddWorkers(10)
	}()

	go func() {
		defer wg.Done()
		time.Sleep(10 * time.Millisecond)
		wp.DeleteWorkers(5)
	}()

	go func() {
		defer wg.Done()
		time.Sleep(20 * time.Millisecond)
		for i := 0; i < 5; i++ {
			wp.AddMessage("concurrent message")
		}
	}()

	go func() {
		defer wg.Done()
		time.Sleep(30 * time.Millisecond)
		wp.DeleteWorkers(3)
	}()

	wg.Wait()

	if wp.workers != 2 {
		t.Errorf("Ожидалось 2 воркера после конкурентных операций, получено %d", wp.workers)
	}
}

func TestWorkerPool_ContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	wp := NewWorkerPullWithContext(ctx)

	wp.AddWorkers(3)
	cancel() // Отменяем контекст

	// Даем время на завершение воркеров
	time.Sleep(100 * time.Millisecond)

	if wp.workers != 0 {
		t.Errorf("Ожидалось 0 воркеров после отмены контекста, получено %d", wp.workers)
	}
}