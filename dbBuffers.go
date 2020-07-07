package main

import (
    "time"
    "sync"
)


type History struct {
    timestamp time.Time
    op        string
    id        int
    data      Person
}

type HistoryStore struct {
    data      []History
    sync.RWMutex
}

func InitHistoryStore() *HistoryStore {
    return &HistoryStore{
        data: []History{},
    }
}

func (h *HistoryStore) Append(op string, id int, data Person) {
    h.Lock()
    defer h.Unlock()
    newEntry := History{
        timestamp: time.Now(),
        op: op,
        id: id,
        data: data,
    }
    h.data = append(h.data, newEntry)
}

func (h *HistoryStore) Get() []History {
    h.RLock()
    defer h.RUnlock()
    data := h.data
    return data
}
