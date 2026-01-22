//go:build !mobile

package main

import (
  "log"
  "runtime"
  "rtgs-client/rgl"
  "github.com/go-gl/glfw/v3.3/glfw"
)

const (
  width  = 800
  height = 600
)

func init() {
  runtime.LockOSThread()
}

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

  if err := glfw.Init(); err != nil {
    log.Fatalln("failed to init glfw:", err)
  }
  defer glfw.Terminate()

  glfw.WindowHint(glfw.ContextVersionMajor, 2)
  glfw.WindowHint(glfw.ContextVersionMinor, 1)
  glfw.WindowHint(glfw.Resizable, glfw.True)

  window, err := glfw.CreateWindow(width, height, "RTGS Client", nil, nil)
  if err != nil {
    panic(err)
  }

  window.MakeContextCurrent()
  glfw.SwapInterval(1)

  if err := rgl.Init(); err != nil {
    panic(err)
  }

  shaders := &Shaders{
    Vertex: `
      attribute vec3 aPosition;
      uniform mat4 uMVP;
      void main() {
        gl_Position = uMVP * vec4(aPosition, 1.0);
      }
    `,
    Fragment: `
      uniform vec4 uColor;
        void main() {
          gl_FragColor = uColor;
        }
    `,
  }
  
  game := NewGame(worldState, shaders)

  // Main loop
  for !window.ShouldClose() {
    w, h := window.GetSize()
    game.Draw(w, h)

    window.SwapBuffers()
    glfw.PollEvents()
  }
}
