package main

import (
  "github.com/go-gl/gl/v2.1/gl"
  "github.com/go-gl/glfw/v3.3/glfw"
)

type Game struct {
  renderer   *Renderer
  worldState *WorldState
  window     *glfw.Window
}

func NewGame(worldState *WorldState, window *glfw.Window) *Game {
  return &Game{
    renderer:   NewRenderer(),
    worldState: worldState,
    window:     window,
  }
}

func (g *Game) Draw() {
  w, h := g.window.GetSize()
  gl.Viewport(0, 0, int32(w), int32(h))
  gl.ClearColor(0.118, 0.118, 0.157, 1.0)
  gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
  gl.Enable(gl.DEPTH_TEST)

  aspect := float32(w) / float32(h)
  mvp := g.renderer.GetMVP(aspect)

  gridVerts := g.renderer.GetGridVertices(10)
  g.renderer.DrawVertices(gridVerts, [4]float32{0.235, 0.235, 0.314, 1.0}, gl.LINES, mvp)

  for _, user := range g.worldState.GetUsers() {
    if !user.IsActive {
      continue
    }
    
    gl.Enable(gl.BLEND)
    gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

    pos := Vec3{
      X: user.Location.X + 0.5,
      Y: user.Location.Y + 0.5,
      Z: user.Location.Z + 0.5,
    }
    cubeVerts := g.renderer.GetCubeVertices(pos)
    color := GetColorForUserType(user.UserType)
    g.renderer.DrawVertices(
      cubeVerts,
      [4]float32{
        float32(color[0]) / 255.0,
        float32(color[1]) / 255.0,
        float32(color[2]) / 255.0,
        0.5,
      },
      gl.TRIANGLES,
      mvp,
    )
    
    edges := g.renderer.GetCubeEdgesFromVertices(pos)
    g.renderer.DrawVertices(
      edges,
      [4]float32{
        float32(color[0]) / 255.0,
        float32(color[1]) / 255.0,
        float32(color[2]) / 255.0,
        1.0,
      },
      gl.LINES,
      mvp,
    )
    
    gl.Disable(gl.BLEND)
  }
}
