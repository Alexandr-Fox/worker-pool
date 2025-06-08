package workerPool

import (
    "context"
    "errors"
    "fmt"
    "sync"
)

type WorkerPool struct {
    messagesCh chan string
    statusCh   chan int
    counter    int
    workers    int
    ctx        context.Context
    cancel     context.CancelFunc
    closed     bool
    wg         *sync.WaitGroup
}

func NewWorkerPull() *WorkerPool {
    ctx, cancel := context.WithCancel(context.Background())
    var wg sync.WaitGroup

    return &WorkerPool{
        messagesCh: make(chan string),
        statusCh:   make(chan int),
        counter:    1,
        workers:    0,
        ctx:        ctx,
        cancel:     cancel,
        wg:         &wg,
        closed:     false,
    }
}

func (wp *WorkerPool) AddWorkers(count int) error {
    if wp.closed {
        return errors.New("WorkerPool уже закрыт")
    }

    for i := 0; i < count; i++ {
        w := &worker{
            id:         wp.counter,
            messagesCh: wp.messagesCh,
            statusCh:   wp.statusCh,
            wg:         wp.wg,
            ctx:        wp.ctx,
        }

        wp.counter++
        wp.workers++
        go w.run()
        wp.wg.Add(1)
    }

    return nil
}

func (wp *WorkerPool) DeleteWorkers(count int) error {
    if wp.closed {
        return errors.New("WorkerPool уже закрыт")
    }

    if count > wp.workers {
        return errors.New(fmt.Sprintf("число запущенных воркеров: %d, вы пытаетесь удалить %d", wp.workers, count))
    }

    for i := 0; i < count; i++ {
        wp.statusCh <- -1
        wp.workers--
    }

    return nil
}

func (wp *WorkerPool) Close() error {
    if wp.closed {
        return errors.New("WorkerPool уже закрыт")
    }

    wp.cancel()
    wp.workers = 0
    close(wp.statusCh)
    close(wp.messagesCh)
    wp.closed = true

    wp.wg.Wait()

    return nil
}

func (wp *WorkerPool) AddMessage(msg string) error {
    if wp.closed {
        return errors.New("WorkerPool уже закрыт")
    }

    if wp.workers == 0 {
        return errors.New("нет воркеров")
    }

    wp.messagesCh <- msg

    return nil
}
