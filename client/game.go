package main

import (
  "rtgs-client/rgl"
)

type Game struct {
  renderer   *Renderer
  worldState *WorldState
}

func NewGame(worldState *WorldState, shaders *Shaders) *Game {
  return &Game{
    renderer:   NewRenderer(shaders),
    worldState: worldState,
  }
}

func (g *Game) Draw(width, height int) {
  
  rgl.Viewport(0, 0, int32(width), int32(height))
  rgl.ClearColor(0.118, 0.118, 0.157, 1.0)
  rgl.Clear(rgl.COLOR_BUFFER_BIT | rgl.DEPTH_BUFFER_BIT)
  rgl.Enable(rgl.DEPTH_TEST)
  
  aspect := float32(width) / float32(height)
  mvp := g.renderer.GetMVP(aspect)

  gridVerts := g.renderer.GetGridVertices(10)
  g.renderer.DrawVertices(gridVerts, [4]float32{0.235, 0.235, 0.314, 1.0}, rgl.LINES, mvp)

  for _, user := range g.worldState.GetUsers() {
    if !user.IsActive {
      continue
    }

    rgl.Enable(rgl.BLEND)
    rgl.BlendFunc(rgl.SRC_ALPHA, rgl.ONE_MINUS_SRC_ALPHA)

    pos := Vec3{
      X: user.Location.X + 0.5,
      Y: user.Location.Y + 0.5,
      Z: user.Location.Z + 0.5,
    }
    cubeVerts := g.renderer.GetCubeVertices(pos)
    color := GetColorForUserType(user.UserType)
    g.renderer.DrawVertices(cubeVerts, [4]float32{
      float32(color[0]) / 255.0,
      float32(color[1]) / 255.0,
      float32(color[2]) / 255.0,
      0.5,
    }, rgl.TRIANGLES, mvp)

    edges := g.renderer.GetCubeEdgesFromVertices(pos)
    g.renderer.DrawVertices(edges, [4]float32{
      float32(color[0]) / 255.0,
      float32(color[1]) / 255.0,
      float32(color[2]) / 255.0,
      1.0,
    }, rgl.LINES, mvp)

    rgl.Disable(rgl.BLEND)
  }
}
