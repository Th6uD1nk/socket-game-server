package main

import (
  "log"
  "github.com/hajimehoshi/ebiten/v2"
)

func main() {
  worldState := NewWorldState()
  
  // Start UDP client
  client, err := NewUDPClient("127.0.0.1:8888", worldState)
  if err != nil {
    log.Fatalf("Cannot create UDP client: %v", err)
  }
  defer client.Conn.Close()
  
  client.StartReceiving()
  client.StartSending()
  
  // Start Ebiten GUI
  ebiten.SetWindowSize(800, 600)
  ebiten.SetWindowTitle("RTGS Client")
  
  game := NewGame(worldState)
  
  if err := ebiten.RunGame(game); err != nil {
    log.Fatalf("Ebiten run failed: %v", err)
  }
}
