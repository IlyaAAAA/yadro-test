package queue

import (
	"container/list"
	"errors"
)

// Queue is a queue
type Queue interface {
	Poll() (string, error)
	Len() int
	Add(string)
	Remove()
}

type queueImpl struct {
	*list.List
}

func (q *queueImpl) Add(v string) {
	q.PushBack(v)
}

func (q *queueImpl) Remove() {
	e := q.Front()
	q.List.Remove(e)
}

func (q *queueImpl) Len() int {
	return q.List.Len()
}

func (q *queueImpl) Poll() (string, error) {

	if q.List.Len() == 0 {
		return "", errors.New("очередь пустая")
	}

	front := q.Front()
	q.List.Remove(front)

	return front.Value.(string), nil
}

func New() Queue {
	return &queueImpl{list.New()}
}
