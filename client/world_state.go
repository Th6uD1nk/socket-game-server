package main

import (
  "sync"
  "time"
)

type Vec3 struct {
  X, Y, Z float64
}

type UserType string

const (
  UserTypePlayer UserType = "player"
  UserTypeBot  UserType = "bot"
  UserTypeAdmin  UserType = "admin"
)

type User struct {
  ID               string
  UserType         UserType
  Location         Vec3
  PreviousLocation Vec3
  Orientation      float32
  IsActive         bool
  LastUpdate       time.Time
  Color            [3]uint8
}

type WorldState struct {
  mu  sync.RWMutex
  Users map[string]*User
}

func NewWorldState() *WorldState {
  return &WorldState{
    Users: make(map[string]*User),
  }
}

func (w *WorldState) UpdateUser(user *User) {
  w.mu.Lock()
  defer w.mu.Unlock()
  w.Users[user.ID] = user
}

func (w *WorldState) RemoveUser(id string) {
  w.mu.Lock()
  defer w.mu.Unlock()
  delete(w.Users, id)
}

func (w *WorldState) GetUsers() []*User {
  w.mu.RLock()
  defer w.mu.RUnlock()
  
  users := make([]*User, 0, len(w.Users))
  for _, u := range w.Users {
    users = append(users, u)
  }
  return users
}

func (w *WorldState) GetUsersByType(userType UserType) []*User {
  w.mu.RLock()
  defer w.mu.RUnlock()
  
  users := make([]*User, 0)
  for _, u := range w.Users {
    if u.UserType == userType {
      users = append(users, u)
    }
  }
  return users
}

func GetColorForUserType(userType UserType) [3]uint8 {
  switch userType {
  case UserTypePlayer:
    return [3]uint8{0, 255, 100}
  case UserTypeBot:
    return [3]uint8{255, 150, 0}
  case UserTypeAdmin:
    return [3]uint8{255, 50, 50}
  default:
    return [3]uint8{150, 150, 150}
  }
}
