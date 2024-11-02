package server

import (
	"net"
	"sync"
	"time"
)

const InitialUsersExpected = 8

type User struct {
	conn net.Conn
	Addr string
}

func NewUser(conn net.Conn) *User {
	conn.SetReadDeadline(time.Time{})
	return &User{
		conn: conn,
		Addr: conn.RemoteAddr().String(),
	}
}

func (u *User) Leave() error {
	return u.conn.Close()
}

func (u *User) Send(b []byte) error {
	_, err := u.conn.Write(b)
	return err
}

type UsersRepository struct {
	mu           sync.Mutex
	users        []*User
	vacantPlaces int
}

func NewUsersRepo() *UsersRepository {
	return &UsersRepository{
		mu:    sync.Mutex{},
		users: make([]*User, 0, InitialUsersExpected),
	}
}

func (r *UsersRepository) FindUser(addr string) *User {
	var user *User
	for i := 0; i < len(r.users); i++ {
		curUser := r.users[i]
		if curUser == nil {
			continue
		}
		if r.users[i].Addr == addr {
			user = r.users[i]
		}
	}
	return user
}

func (r *UsersRepository) AddUser(user *User) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.vacantPlaces > 0 {
		for index, u := range r.users {
			if u != nil {
				continue
			}
			r.users[index] = user
		}
		r.vacantPlaces--
	} else {
		r.users = append(r.users, user)
	}
}

func (r *UsersRepository) RemoveUser(addr string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i := 0; i < len(r.users); i++ {
		user := r.users[i]
		if user == nil {
			continue
		}
		if user.Addr == addr {
			r.users[i] = nil
			r.vacantPlaces++
		}
	}
}

func (r *UsersRepository) IterateUsers(f func(*User)) {
	for _, u := range r.users {
		if u != nil {
			f(u)
		}
	}
}
