package workerPool

import (
    "context"
    "fmt"
    "os"
    "sync"
)

type worker struct {
    id          int
    messagesCh  chan string
    deletingCh  chan bool
    wg          *sync.WaitGroup
    ctx         context.Context
    afterDelete func()
    setBusy     func()
    setUnbusy   func()
}

func (w *worker) run() {
    defer w.wg.Done()
    defer w.afterDelete()

    fmt.Printf("Worker %d [%d]: запущен\n", w.id, os.Getpid())

    for {
        select {
        // Если контекст завершен, удаляем воркер
        case <-w.ctx.Done():
            fmt.Printf("Worker %d [%d]: удален\n", w.id, os.Getpid())
            return
        case status, ok := <-w.deletingCh:
            // Если канал закрылся или пришел сигнал об удалении воркера, удаляем воркер
            if ok == false || status {
                fmt.Printf("Worker %d [%d]: удален\n", w.id, os.Getpid())
                return
            }
        case msg, ok := <-w.messagesCh:
            // Если канал закрылся, удаляем воркер
            if ok == false {
                fmt.Printf("Worker %d [%d]: удален\n", w.id, os.Getpid())
                return
            }

            w.setBusy()
            fmt.Printf("Worker %d [%d]: %s\n", w.id, os.Getpid(), msg)
            w.setUnbusy()
        }
    }
}
