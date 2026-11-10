package repository

import (
	"crypto/rand"
	"encoding/hex"
	"sync"

	"pharmacy/medicine"
)

type Memory struct {
	mu        sync.RWMutex
	medicines map[string]medicine.Medicine
}

func NewMemory() *Memory {
	return &Memory{medicines: make(map[string]medicine.Medicine)}
}

func newID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (r *Memory) Create(m medicine.Medicine) (string, error) {
	id, err := newID()
	if err != nil {
		return "", err
	}
	m.ID = id

	r.mu.Lock()
	r.medicines[id] = m
	r.mu.Unlock()
	return id, nil
}

func (r *Memory) GetByID(id string) (medicine.Medicine, error) {
	r.mu.RLock()
	m, ok := r.medicines[id]
	r.mu.RUnlock()
	if !ok {
		return medicine.Medicine{}, medicine.ErrNotFound
	}
	return m, nil
}

func (r *Memory) Sell(id string, count int) error {
	if count <= 0 {
		return medicine.ErrInvalidCount
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	m, ok := r.medicines[id]
	if !ok {
		return medicine.ErrNotFound
	}
	if m.Quantity < count {
		return medicine.ErrOutOfStock
	}
	m.Quantity -= count
	r.medicines[id] = m
	return nil
}
