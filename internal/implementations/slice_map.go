package implementations

import (
	"echodb/internal/interfaces"
	"echodb/internal/models"
)

// SliceMap 切片实现
type SliceMap struct {
	data []*models.Person
}

func NewSliceMap() interfaces.MyMap {
	return &SliceMap{
		data: make([]*models.Person, 0),
	}
}

func (sm *SliceMap) Put(id int, person *models.Person) {
	for i, p := range sm.data {
		if p.ID == id {
			sm.data[i] = person
			return
		}
	}
	sm.data = append(sm.data, person)
}

func (sm *SliceMap) Get(id int) *models.Person {
	for _, p := range sm.data {
		if p.ID == id {
			return p
		}
	}
	return nil
}

func (sm *SliceMap) Delete(id int) {
	for i, p := range sm.data {
		if p.ID == id {
			sm.data[i] = sm.data[len(sm.data)-1]
			sm.data = sm.data[:len(sm.data)-1]
			return
		}
	}
}

