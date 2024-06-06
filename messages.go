package main

import (
	"errors"
	"sync"

	"github.com/google/uuid"
)

type Message struct {
	UUID    uuid.UUID `json:"uuid"`
	Content string    `json:"content"`
}

type MessageBucket struct {
	mu     sync.Mutex
	Bucket map[uuid.UUID]string
}

func (b *MessageBucket) Insert(m *Message) error {
	uuid := uuid.New()
	m.UUID = uuid

	b.mu.Lock()
	if _, ok := b.Bucket[uuid]; ok {
		b.mu.Unlock()
		return errors.New("failed to store content")
	}

	b.Bucket[uuid] = m.Content
	b.mu.Unlock()

	return nil
}

func (b *MessageBucket) RetrieveAndRemove(m *Message) error {
	b.mu.Lock()
	v, ok := b.Bucket[m.UUID]
	if !ok {
		b.mu.Unlock()
		return errors.New("no corresponding content found for given key")
	}

	delete(b.Bucket, m.UUID)
	b.mu.Unlock()

	m.Content = v

	return nil
}
