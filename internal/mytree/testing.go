package mytree

import (
	"echodb/internal/model"
	"echodb/internal/mymap"
)

type MyTree interface {
	mymap.MyMap
	RangeQuery(minId, maxId int) []*model.Person
}