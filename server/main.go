package main

import (
  "fmt"
  "net"
  "sync"
  "time"
  "math/rand"
  "math"
  "encoding/json"
)

// tmp
type WorldUpdate struct {
  Type  string     `json:"type"`
  Users []UserData `json:"users"`
}

// tmp
type UserData struct {
  ID          string     `json:"id"`
  UserType    string     `json:"user_type"`
  Location    [3]float32 `json:"location"`
  Orientation float32    `json:"orientation"`
  IsActive    bool       `json:"is_active"`
}

type Client struct {
  addr     *net.UDPAddr
  lastSeen time.Time
  user     *User
}

type Server struct {
  conn    *net.UDPConn
  clients map[string]*Client
  mu      sync.RWMutex
}

func NewServer(port int) (*Server, error) {
  fmt.Println("Creating UDP address...")
  addr := net.UDPAddr{
    Port: port,
    IP:   net.ParseIP("0.0.0.0"),
  }
  
  fmt.Printf("Binding to %s:%d...\n", addr.IP, addr.Port)
  conn, err := net.ListenUDP("udp", &addr)
  if err != nil {
    return nil, err
  }
  
  fmt.Println("UDP socket bound successfully")
  
  return &Server{
    conn:    conn,
    clients: make(map[string]*Client),
  }, nil
}

func (server *Server) addOrUpdateClient(addr *net.UDPAddr) {
  server.mu.Lock()
  defer server.mu.Unlock()
  
  key := addr.String()
  if client, exists := server.clients[key]; exists {
    client.lastSeen = time.Now()
  } else {
    server.clients[key] = &Client{
      addr:     addr,
      lastSeen: time.Now(),
    }
    fmt.Printf("+ New client: %s\n", addr.String())
  }
}

func (server *Server) listClients() {
  server.mu.RLock()
  defer server.mu.RUnlock()
  
  fmt.Println("\n=== Clients list ===")
  if len(server.clients) == 0 {
      fmt.Println("No connected client")
  } else {
    for key, client := range server.clients {
      fmt.Printf("- %s\n", key)
      fmt.Printf("  Type: %s\n", client.user.userType)
      fmt.Printf("  Location: (%.2f, %.2f, %.2f)\n", 
          client.user.location.x, client.user.location.y, client.user.location.z)
      fmt.Printf("  Orientation: %.2f°\n", client.user.orientation)
      fmt.Printf("  Active: %t\n", client.user.isActive)
      fmt.Printf("  Last activity: %s\n", time.Since(client.lastSeen).Round(time.Second))
    }
  }
  fmt.Println("========================\n")
}

func (server *Server) cleanInactiveClients(timeout time.Duration) {
  server.mu.Lock()
  defer server.mu.Unlock()
    
  now := time.Now()
  for key, client := range server.clients {
    if now.Sub(client.lastSeen) > timeout {
      fmt.Printf("x Client timeout: %s\n", key)
      delete(server.clients, key)
    }
  }
}

func randomSpawn(id string, userType UserType, conn *net.UDPConn,
  minX, maxX, minY, maxY, minZ, maxZ float32) *User {
  
  var location Vector3

  location.x = float32(math.Round(float64(minX + rand.Float32()*(maxX-minX))))
  location.y = float32(math.Round(float64(minY + rand.Float32()*(maxY-minY))))
  location.z = float32(math.Round(float64(minZ + rand.Float32()*(maxZ-minZ))))

  orientation := rand.Float32() * 360.0
  
  user := NewUser(id, userType, conn)
  user.location = location
  user.orientation = orientation
  
  return user
}

func (server *Server) broadcastWorldState() {
  server.mu.RLock()
  
  worldUpdate := WorldUpdate{
    Type:  "world_update",
    Users: make([]UserData, 0, len(server.clients)),
  }
  
  for _, client := range server.clients {
    if client.user != nil {
      worldUpdate.Users = append(worldUpdate.Users, UserData{
        ID:          client.user.id,
        UserType:    string(client.user.userType),
        Location:    [3]float32{client.user.location.x, client.user.location.y, client.user.location.z},
        Orientation: client.user.orientation,
        IsActive:    client.user.isActive,
      })
    }
  }
  
  server.mu.RUnlock()
  
  // convert to json
  data, err := json.Marshal(worldUpdate)
  if err != nil {
    fmt.Printf("JSON marshal error: %v\n", err)
    return
  }
  
  // send to all clients
  server.mu.RLock()
  for _, client := range server.clients {
    _, err := server.conn.WriteToUDP(data, client.addr)
    if err != nil {
      fmt.Printf("Broadcast error to %s: %v\n", client.addr.String(), err)
    }
  }
  server.mu.RUnlock()
}

func (server *Server) Start() {
  defer server.conn.Close()
  
  fmt.Printf("UDP server started on port %d\n", server.conn.LocalAddr().(*net.UDPAddr).Port)
  
  // clean inactive clients
  go func() {
    ticker := time.NewTicker(10 * time.Second)
    defer ticker.Stop()
    for range ticker.C {
      server.cleanInactiveClients(30 * time.Second)
    }
  }()
  
  // display client list periodically
  go func() {
    ticker := time.NewTicker(15 * time.Second)
    defer ticker.Stop()
    for range ticker.C {
      server.listClients()
    }
  }()
  
  // broadcast world state each 100 ms
  go func() {
    ticker := time.NewTicker(100 * time.Millisecond)
    defer ticker.Stop()
    for range ticker.C {
      server.broadcastWorldState()
    }
  }()

  buffer := make([]byte, 1024)
  
  for {
    nByte, addr, err := server.conn.ReadFromUDP(buffer)
    if err != nil {
      fmt.Printf("Read error %v\n", err)
      continue
    }
    
    clientKey := addr.String()
    if _, exists := server.clients[clientKey]; !exists {
      user := randomSpawn(clientKey, UserTypePlayer, server.conn, 0, 10, 0, 0, 0, 10)
      
      client := &Client{
        addr:     addr,
        lastSeen: time.Now(),
        user:     user,
      }
      
      server.mu.Lock()
      server.clients[clientKey] = client
      server.mu.Unlock()
    
      fmt.Printf("+ new client spawned at (%.2f, %.2f, %.2f) orientation: %.2f\n",
        user.location.x, user.location.y, user.location.z, user.orientation)
    }

    server.addOrUpdateClient(addr)
    
    message := string(buffer[:nByte])
    fmt.Printf("+ received from %s: %s\n", addr.String(), message)
    
    // response := fmt.Sprintf("ACK: %s", message)
    // _, err = server.conn.WriteToUDP([]byte(response), addr)
    // if err != nil {
    //   fmt.Printf("Sent error: %v\n", err)
    // }
  }
}

func main() {
  fmt.Println("=== RTGS Server ===")
  server, err := NewServer(8888)
  if err != nil {
    fmt.Printf("Error on create: %v\n", err)
    return
  }
  server.Start()
}
