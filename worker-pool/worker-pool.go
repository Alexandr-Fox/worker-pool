package workerPool

import (
    "context"
    "errors"
    "fmt"
    "sync"
)

type WorkerPool struct {
    messagesCh      chan string
    deletingCh      chan bool
    counter         int
    workers         int
    workersMutex    *sync.Mutex
    busyWorkers     int
    busyMutex       *sync.Mutex
    ctx             context.Context
    cancel          context.CancelFunc
    closed          bool
    wg              *sync.WaitGroup
    wgAfterDeleting func()
    wgDeletingMutex *sync.Mutex
}

func NewWorkerPull() *WorkerPool {
    return NewWorkerPullWithContext(context.Background())
}

func NewWorkerPullWithContext(ctx context.Context) *WorkerPool {
    ctx, cancel := context.WithCancel(ctx)
    var wg sync.WaitGroup
    var wm sync.Mutex
    var bm sync.Mutex
    var wdm sync.Mutex

    wp := &WorkerPool{
        messagesCh:      make(chan string),
        deletingCh:      make(chan bool),
        counter:         1,
        workers:         0,
        workersMutex:    &wm,
        busyWorkers:     0,
        busyMutex:       &bm,
        ctx:             ctx,
        cancel:          cancel,
        wg:              &wg,
        wgDeletingMutex: &wdm,
        closed:          false,
    }

    wp.wgAfterDeleting = func() {
        wp.afterDeleteWorker()
    }

    return wp
}

func (wp *WorkerPool) afterDeleteWorker() {
    wp.workersMutex.Lock()
    wp.workers--
    wp.workersMutex.Unlock()
}

func (wp *WorkerPool) setBusyWorker() {
    wp.busyMutex.Lock()
    wp.busyWorkers++
    wp.busyMutex.Unlock()
}

func (wp *WorkerPool) unsetBusyWorker() {
    wp.busyMutex.Lock()
    wp.busyWorkers--
    wp.busyMutex.Unlock()
}

func (wp *WorkerPool) AddWorkers(count int) error {
    if wp.closed {
        return errors.New("WorkerPool уже закрыт")
    }

    for i := 0; i < count; i++ {
        w := &worker{
            id:         wp.counter,
            messagesCh: wp.messagesCh,
            deletingCh: wp.deletingCh,
            wg:         wp.wg,
            ctx:        wp.ctx,
            afterDelete: func() {
                wp.wgAfterDeleting()
            },
            setBusy: func() {
                wp.setBusyWorker()
            },
            setUnbusy: func() {
                wp.unsetBusyWorker()
            },
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

    var wg sync.WaitGroup

    wp.wgDeletingMutex.Lock()

    wp.wgAfterDeleting = func() {
        wp.afterDeleteWorker()
        defer wg.Done()
    }

    for i := 0; i < count; i++ {
        wg.Add(1)
        wp.deletingCh <- true
    }

    wg.Wait()
    wp.wgAfterDeleting = func() {
        wp.afterDeleteWorker()
    }
    wp.wgDeletingMutex.Unlock()

    return nil
}

func (wp *WorkerPool) Close() error {
    if wp.closed {
        return errors.New("WorkerPool уже закрыт")
    }

    wp.cancel()
    wp.closed = true

    wp.wg.Wait()

    close(wp.deletingCh)
    close(wp.messagesCh)

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

func (wp *WorkerPool) GetWorkers() int {
    return wp.workers
}

func (wp *WorkerPool) GetBusyWorkers() int {
    return wp.busyWorkers
}

