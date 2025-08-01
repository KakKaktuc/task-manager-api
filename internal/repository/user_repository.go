package repository

import (
	"errors"
	"sync"
	"github.com/KakKaktuc/task-manager-api/pkg/models"
)

type UserRepository struct {
    users  []models.User
    nextID int
    mu     sync.Mutex
}

func NewUserRepository() *UserRepository {
    return &UserRepository{
        users:  make([]models.User, 0),
        nextID: 1,
    }
}

func (r *UserRepository) GetAll() []models.User {
    r.mu.Lock()
    defer r.mu.Unlock()

    return append([]models.User(nil), r.users...) // возвращаем копию
}

func (r *UserRepository) GetByID(id int) (*models.User, error) {
    r.mu.Lock()
    defer r.mu.Unlock()

    for _, user := range r.users {
        if user.ID == id {
            u := user
            return &u, nil
        }
    }
    return nil, errors.New("user not found")
}

func (r *UserRepository) Create(user models.User) models.User {
    r.mu.Lock()
    defer r.mu.Unlock()

    user.ID = r.nextID
    r.nextID++
    r.users = append(r.users, user)
    return user
}

func (r *UserRepository) Update(id int, updated models.User) (*models.User, error) {
    r.mu.Lock()
    defer r.mu.Unlock()

    for i, user := range r.users {
        if user.ID == id {
            updated.ID = id
            r.users[i] = updated
            return &updated, nil
        }
    }
    return nil, errors.New("user not found")
}

func (r *UserRepository) Delete(id int) error {
    r.mu.Lock()
    defer r.mu.Unlock()

    for i, user := range r.users {
        if user.ID == id {
            r.users = append(r.users[:i], r.users[i+1:]...)
            return nil
        }
    }
    return errors.New("user not found")
}
