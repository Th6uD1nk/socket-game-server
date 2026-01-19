package main

import (
  "encoding/json"
  "fmt"
  "net"
  "time"
)

type UDPClient struct {
  Conn       *net.UDPConn
  WorldState *WorldState
}

func NewUDPClient(addr string, worldState *WorldState) (*UDPClient, error) {
  serverAddr, err := net.ResolveUDPAddr("udp", addr)
  if err != nil {
    return nil, err
  }
  conn, err := net.DialUDP("udp", nil, serverAddr)
  if err != nil {
    return nil, err
  }
  return &UDPClient{
    Conn:       conn,
    WorldState: worldState,
  }, nil
}

type ServerMessage struct {
  Type  string       `json:"type"`
  Users []UserUpdate `json:"users,omitempty"`
}

type UserUpdate struct {
  ID          string      `json:"id"`
  UserType    string      `json:"user_type"`
  Location    [3]float32  `json:"location"`
  Orientation float32     `json:"orientation"`
  IsActive    bool        `json:"is_active"`
}

func (c *UDPClient) StartReceiving() {
  go func() {
    buffer := make([]byte, 4096)
    for {
      n, err := c.Conn.Read(buffer)
      if err != nil {
        fmt.Printf("Receive error: %v\n", err)
        continue
      }
      
      var msg ServerMessage
      if err := json.Unmarshal(buffer[:n], &msg); err != nil {
        fmt.Printf("Parse error: %v\n", err)
        continue
      }
      
      if msg.Type == "world_update" {
        for _, userUpdate := range msg.Users {
          user := &User{
            ID:          userUpdate.ID,
            UserType:    UserType(userUpdate.UserType),
            Location:    Vec3{X: float64(userUpdate.Location[0]), Y: float64(userUpdate.Location[1]), Z: float64(userUpdate.Location[2])},
            Orientation: userUpdate.Orientation,
            IsActive:    userUpdate.IsActive,
            LastUpdate:  time.Now(),
            Color:       GetColorForUserType(UserType(userUpdate.UserType)),
          }
          c.WorldState.UpdateUser(user)
        }
        fmt.Printf("+ World updated: %d users\n", len(msg.Users)) // tmp
      }
    }
  }()
}

func (c *UDPClient) StartSending() {
  go func() {
    counter := 1
    ticker := time.NewTicker(2 * time.Second)
    defer ticker.Stop()
    for range ticker.C {
      message := fmt.Sprintf("Message %d", counter)
      _, err := c.Conn.Write([]byte(message))
      if err != nil {
        fmt.Printf("Send error: %v\n", err)
        continue
      }
      fmt.Printf("- Sent: %s\n", message)
      counter++
    }
  }()
}

