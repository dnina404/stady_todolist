package todo

import "sync"

type List struct {
	tasks map[string]Task
	mtx   sync.RWMutex
}

func NewList() *List {
	return &List{
		tasks: make(map[string]Task),
	}
}
func (l *List) GetTask(title string) (Task, error) {
	l.mtx.RLock()
	defer l.mtx.RUnlock()

	task, ok := l.tasks[title]
	if !ok {
		return Task{}, ErrTaskNotFound
	}
	return task, nil
}
func (l *List) AddTask(task Task) error {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	if _, ok := l.tasks[task.Title]; ok {
		return ErrTaskAlreadyExists
	}

	l.tasks[task.Title] = task

	return nil
}
func (l *List) CompletesTask(t Task, title string, b bool) Task {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	task := l.tasks[title]
	if b == true {
		task.Complete()
	}
	if b == false {
		task.UnComplete()
	}
	l.tasks[title] = task
	return task
}

func (l *List) ListTasks(completed *bool) map[string]Task {
	l.mtx.RLock()
	defer l.mtx.RUnlock()

	tasks := make(map[string]Task)

	for k, v := range l.tasks {
		if completed == nil {
			tasks[k] = v
			continue
		}

		if v.Completed == *completed {
			tasks[k] = v
		}
	}

	return tasks
}

func (l *List) CompleteTask(title string) error {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	task, ok := l.tasks[title]
	if !ok {
		return ErrTaskNotFound
	}

	task.Complete()

	l.tasks[title] = task

	return nil
}

func (l *List) UncompleteTask(title string) error {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	task, ok := l.tasks[title]
	if !ok {
		return ErrTaskNotFound
	}

	task.UnComplete()

	l.tasks[title] = task

	return nil
}

func (l *List) DeleteTask(title string) error {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	_, ok := l.tasks[title]
	if !ok {
		return ErrTaskNotFound
	}
	delete(l.tasks, title)

	return nil
}
