package timer

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

type Tasks struct {
	Name     string
	Fn       func()
	Interval time.Duration
}

type Timer struct {
	Tasks map[string]Tasks
	lock  sync.Mutex
}

func (t *Timer) Add(name string, interval time.Duration, fn func()) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	if _, ok := t.Tasks[name]; ok {
		return errors.New(fmt.Sprintf("新增任务名称\"%s\"冲突", name))
	}
	t.Tasks[name] = Tasks{Name: name, Fn: fn, Interval: interval}
	go t.startTask(name)
	return nil
}

func (t *Timer) Del(name string) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	if _, ok := t.Tasks[name]; !ok {
		return errors.New(fmt.Sprintf("删除任务名称\"%s\"没有定义", name))
	}
	delete(t.Tasks, name)
	return nil
}

func (t *Timer) Update(name string, interval time.Duration, fn func()) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	if _, ok := t.Tasks[name]; !ok {
		return errors.New(fmt.Sprintf("更新任务名称\"%s\"没有定义", name))
	}
	t.Tasks[name] = Tasks{Name: name, Fn: fn, Interval: interval}
	return nil
}

func (t *Timer) Clean() {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.Tasks = map[string]Tasks{}
}

func (t *Timer) List() []Tasks {
	var tasks []Tasks
	for _, v := range t.Tasks {
		tasks = append(tasks, v)
	}
	return tasks
}

func (t *Timer) GetTask(name string) (*Tasks, error) {
	if task, ok := t.Tasks[name]; !ok {
		return nil, errors.New(fmt.Sprintf("任务名称\"%s\"没有定义", name))
	} else {
		return &task, nil
	}
}

func (t *Timer) execIntervalTask(name string) {
	fn := func() {
		defer func() {
			if r := recover(); r != nil {
				log.Print("error", r)
				// fmt.Println("error", r)
			}
		}()
		t.Tasks[name].Fn()
	}
	select {
	case <-time.After(t.Tasks[name].Interval):
		go fn()
	}
}

func (t *Timer) startTask(name string) {
	for {
		_, ok := t.Tasks[name]
		if !ok {
			fmt.Println("delete")
			break
		}
		t.execIntervalTask(name)
	}
}

func NewTimer() *Timer {
	return &Timer{
		Tasks: map[string]Tasks{},
		lock:  sync.Mutex{},
	}
}
