package concurrency

type Sentry struct {
	paused  bool
	stopped bool
	resume chan struct{}
}

func (s *Sentry) Sleep() {
	if s.paused {
		<- s.resume
	}
}

func (s *Sentry) Paused() bool {
	return s.paused
}

func (s *Sentry) Stopped() bool {
	return s.stopped
}

type Worker struct {
	sentry *Sentry
	done chan struct{}
}

func NewWorker() *Worker {
	return &Worker{
		sentry: &Sentry{
			stopped: true,
		},
		done: make(chan struct{}),
	}
}

func (w *Worker) Stopped() bool {
	return w.sentry.stopped
}

func (w *Worker) Paused() bool {
	return w.sentry.paused
}

func (w *Worker) Start(task func(*Sentry)) {
	if !w.sentry.stopped {
		panic("worker has started")
	}
	w.sentry.stopped = false
	w.sentry.resume = make(chan struct{})

	go func() {
		task(w.sentry)
		w.done <- struct{}{}
	}()
}

func (w *Worker) Pause() {
	w.sentry.paused = true
}

func (w *Worker) Resume() {
	w.sentry.paused = false
	close(w.sentry.resume)
	w.sentry.resume = make(chan struct{})
}

func (w *Worker) Stop() {
	if w.sentry.stopped {
		return
	}
	w.sentry.stopped = true
	if w.sentry.paused {
		w.sentry.paused = false
		close(w.sentry.resume)
	}
	<- w.done

	// reset
	w.sentry.paused = false
	w.sentry.resume = nil
}
