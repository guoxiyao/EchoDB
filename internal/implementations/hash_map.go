package implementations

import (
	"echodb/internal/interfaces"
	"echodb/internal/models"
)

// HashMap 哈希表实现
type HashMap struct {
	data map[int]*models.Person
}

func NewHashMap() interfaces.MyMap {
	return &HashMap{
		data: make(map[int]*models.Person),
	}
}

func (hm *HashMap) Put(id int, person *models.Person) {
	hm.data[id] = person
}

func (hm *HashMap) Get(id int) *models.Person {
	return hm.data[id]
}

func (hm *HashMap) Delete(id int) {
	delete(hm.data, id)
}

