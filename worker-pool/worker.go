package workerPool

import (
    "context"
    "fmt"
    "os"
    "sync"
)

type worker struct {
    id         int
    messagesCh chan string
    statusCh   chan int
    wg         *sync.WaitGroup
    ctx        context.Context
}

func (w *worker) run() {
    defer w.wg.Done()

    fmt.Printf("Worker %d [%d]: запущен\n", w.id, os.Getpid())

    for {
        select {
        case <-w.ctx.Done():
            fmt.Printf("Worker %d [%d]: удален\n", w.id, os.Getpid())
            return
        case status, ok := <-w.statusCh:
            if ok == false || status < 0 {
                fmt.Printf("Worker %d [%d]: удален\n", w.id, os.Getpid())
                return
            }
        case msg, ok := <-w.messagesCh:
            if ok == false {
                fmt.Printf("Worker %d [%d]: удален\n", w.id, os.Getpid())
                return
            }

            fmt.Printf("Worker %d [%d]: %s\n", w.id, os.Getpid(), msg)
        }
    }
}
