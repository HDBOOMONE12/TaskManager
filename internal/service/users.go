package service

import (
	"errors"
	"sync"
)

var (
	mu     sync.RWMutex
	users  = make(map[int]User)
	nextID = 1

	ErrEmptyName  = errors.New("empty name")
	ErrEmptyEmail = errors.New("empty email")
)

type User struct {
	ID    int
	Name  string
	Email string
}

func CreateUser(name, email string) (User, error) {
	if name == "" {
		return User{}, ErrEmptyName
	}
	if email == "" {
		return User{}, ErrEmptyEmail
	}

	mu.Lock()
	defer mu.Unlock()

	u := User{
		ID:    nextID,
		Name:  name,
		Email: email,
	}
	users[u.ID] = u
	nextID++
	return u, nil
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

func UpdateUserByID(id int, name, email string) (User, error) {
	if name == "" {
		return User{}, ErrEmptyName
	}
	if email == "" {
		return User{}, ErrEmptyEmail
	}

	mu.Lock()
	defer mu.Unlock()

	u, ok := users[id]
	if !ok {
		return User{}, ErrUserNotFound
	}
	u.Name = name
	u.Email = email
	users[id] = u
	return u, nil
}

func PatchUserByID(id int, name, email *string) (User, error) {
	mu.Lock()
	defer mu.Unlock()

	u, ok := users[id]
	if !ok {
		return User{}, ErrUserNotFound
	}

	if name != nil {
		if *name == "" {
			return User{}, ErrEmptyName
		}
		u.Name = *name
	}
	if email != nil {
		if *email == "" {
			return User{}, ErrEmptyEmail
		}
		u.Email = *email
	}

	users[id] = u
	return u, nil
}

func DeleteUserByID(id int) bool {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := users[id]; !ok {
		return false
	}
	delete(users, id)
	return true
}
