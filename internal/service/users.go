package service

import "sync"

var (
	mu     sync.RWMutex
	users  = make(map[int]User)
	nextID = 1
)

type User struct {
	ID    int
	Name  string
	Email string
}

func CreateUser(name, email string) User {
	mu.Lock()
	defer mu.Unlock()

	u := User{
		ID:    nextID,
		Name:  name,
		Email: email,
	}
	users[u.ID] = u
	nextID++
	return u
}

func ListUsers() []User {
	mu.RLock()
	defer mu.RUnlock()
	list := make([]User, 0, len(users))
	for _, v := range users {
		list = append(list, v)
	}
	return list
}

func GetUserByID(id int) (User, bool) {
	mu.RLock()
	defer mu.RUnlock()
	u, ok := users[id]

	return u, ok

}
