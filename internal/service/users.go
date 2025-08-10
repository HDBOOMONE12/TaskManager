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

func UpdateUserByID(id int, name, email string) (User, bool) {
	mu.Lock()
	defer mu.Unlock()
	updated, ok := users[id]
	if !ok {
		return User{}, false
	}
	updated.Name = name
	updated.Email = email
	users[id] = updated
	return updated, true
}
func PatchUserByID(id int, name, email *string) (User, bool) {
	mu.Lock()
	defer mu.Unlock()
	updated, ok := users[id]
	if !ok {
		return User{}, false
	}
	if name != nil {
		updated.Name = *name
	}
	if email != nil {
		updated.Email = *email
	}
	users[id] = updated
	return updated, true
}

func DeleteUserByID(id int) bool {
	mu.Lock()
	defer mu.Unlock()
	_, ok := users[id]
	if !ok {
		return false
	}
	delete(users, id)
	return true
}
