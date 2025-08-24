package models

import "fmt"

// Person 结构体定义
type Person struct {
	ID        int
	FirstName string
	LastName  string
	Age       int
	Gender    string
	Email     string
}

// NewPerson 构造函数
func NewPerson(id int, firstName, lastName string, age int, gender, email string) *Person {
	return &Person{
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
		Age:       age,
		Gender:    gender,
		Email:     email,
	}
}

// String 方法实现字符串表示
func (p *Person) String() string {
	return fmt.Sprintf("Person{id=%d, firstName='%s', lastName='%s', age=%d, gender='%s', email='%s'}",
		p.ID, p.FirstName, p.LastName, p.Age, p.Gender, p.Email)
}
