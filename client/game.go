package main

import (
  "image/color"
  "github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
  renderer   *Renderer
  worldState *WorldState
}

func NewGame(worldState *WorldState) *Game {
  return &Game{
    renderer:   NewRenderer(),
    worldState: worldState,
  }
}

func (g *Game) Update() error {
  return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
  screen.Fill(color.RGBA{30, 30, 40, 255})
  
  g.renderer.DrawGrid(screen, 10)
  
  users := g.worldState.GetUsers()
  for _, user := range users {
    if !user.IsActive {
      continue
    }
    
    colorRGB := GetColorForUserType(user.UserType)
    cubeColor := color.RGBA{colorRGB[0], colorRGB[1], colorRGB[2], 255}
    
    g.renderer.DrawCube(screen, user.Location, cubeColor)
  }
}

func (g *Game) Layout(w, h int) (int, int) {
  return 800, 600
}
