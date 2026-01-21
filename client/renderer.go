package main

import (
  "fmt"
  "log"
  "github.com/go-gl/gl/v2.1/gl"
  "github.com/go-gl/mathgl/mgl32"
)

type Renderer struct {
  program   uint32
  mvpLoc  int32
  colorLoc  int32
  positionLoc uint32
  cameraDist float32
  angleX   float32
  angleY   float32
}

func NewRenderer() *Renderer {
  r := &Renderer{
  cameraDist: 30,
  angleX:   45,
  angleY:   12,
  }
  r.initShaders()
  return r
}

func (r *Renderer) GetCubeEdgesFromVertices(pos Vec3) []float32 {
  x, y, z := float32(pos.X), float32(pos.Y), float32(pos.Z)
  s := float32(0.5)

  vertices := [8][3]float32{
    {x - s, y - s, z - s}, // v0
    {x + s, y - s, z - s}, // v1
    {x + s, y + s, z - s}, // v2
    {x - s, y + s, z - s}, // v3
    {x - s, y - s, z + s}, // v4
    {x + s, y - s, z + s}, // v5
    {x + s, y + s, z + s}, // v6
    {x - s, y + s, z + s}, // v7
  }

  edgesIdx := [][2]int{
    {0, 1}, {1, 2}, {2, 3}, {3, 0}, // Back
    {4, 5}, {5, 6}, {6, 7}, {7, 4}, // Front
    {0, 4}, {1, 5}, {2, 6}, {3, 7}, // Sides
  }

  edges := make([]float32, 0, len(edgesIdx)*6)
  for _, e := range edgesIdx {
    edges = append(edges,
      vertices[e[0]][0], vertices[e[0]][1], vertices[e[0]][2],
      vertices[e[1]][0], vertices[e[1]][1], vertices[e[1]][2],
    )
  }
  return edges
}

func (r *Renderer) GetCubeVertices(pos Vec3) []float32 {
  x, y, z := float32(pos.X), float32(pos.Y), float32(pos.Z)
  s := float32(0.5)
  return []float32{
    // Front
    x-s, y-s, z+s,  x+s, y-s, z+s,  x+s, y+s, z+s,
    x-s, y-s, z+s,  x+s, y+s, z+s,  x-s, y+s, z+s,
    // Back
    x+s, y-s, z-s,  x-s, y-s, z-s,  x-s, y+s, z-s,
    x+s, y-s, z-s,  x-s, y+s, z-s,  x+s, y+s, z-s,
    // Left
    x-s, y-s, z-s,  x-s, y-s, z+s,  x-s, y+s, z+s,
    x-s, y-s, z-s,  x-s, y+s, z+s,  x-s, y+s, z-s,
    // Right
    x+s, y-s, z+s,  x+s, y-s, z-s,  x+s, y+s, z-s,
    x+s, y-s, z+s,  x+s, y+s, z-s,  x+s, y+s, z+s,
    // Top
    x-s, y+s, z+s,  x+s, y+s, z+s,  x+s, y+s, z-s,
    x-s, y+s, z+s,  x+s, y+s, z-s,  x-s, y+s, z-s,
    // Bottom
    x-s, y-s, z-s,  x+s, y-s, z-s,  x+s, y-s, z+s,
    x-s, y-s, z-s,  x+s, y-s, z+s,  x-s, y-s, z+s,
  }
}

func (r *Renderer) GetGridVertices(gridSize int) []float32 {
  var lines []float32
  for i := -gridSize; i <= gridSize; i++ {
  // X lines
  lines = append(lines,
    float32(-gridSize), 0, float32(i),
    float32(gridSize), 0, float32(i),
  )
  // Z lines
  lines = append(lines,
    float32(i), 0, float32(-gridSize),
    float32(i), 0, float32(gridSize),
  )
  }
  return lines
}

func (r *Renderer) initShaders() {
  
  vertexSrc := 
  `
    attribute vec3 aPosition;
    uniform mat4 uMVP;
    void main() {
      gl_Position = uMVP * vec4(aPosition, 1.0);
    }
  `
  
  fragmentSrc :=
  `
    uniform vec4 uColor;
    void main() {
      gl_FragColor = uColor;
    }
  `
  
  program, err := compileProgram(vertexSrc, fragmentSrc)
  if err != nil {
    log.Fatalln("Failed to compile shaders:", err)
  }
  
  r.program = program
  r.mvpLoc = gl.GetUniformLocation(program, gl.Str("uMVP\x00"))
  r.colorLoc = gl.GetUniformLocation(program, gl.Str("uColor\x00"))
  r.positionLoc = uint32(gl.GetAttribLocation(program, gl.Str("aPosition\x00")))
}

func compileProgram(vertexSrc, fragmentSrc string) (uint32, error) {
  vertexShader, err := compileShader(vertexSrc+"\x00", gl.VERTEX_SHADER)
  if err != nil {
    return 0, err
  }
  fragmentShader, err := compileShader(fragmentSrc+"\x00", gl.FRAGMENT_SHADER)
  if err != nil {
    return 0, err
  }

  program := gl.CreateProgram()
  gl.AttachShader(program, vertexShader)
  gl.AttachShader(program, fragmentShader)
  gl.LinkProgram(program)

  var status int32
  gl.GetProgramiv(program, gl.LINK_STATUS, &status)
  if status == gl.FALSE {
    var logLength int32
    gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)
    log := make([]byte, logLength+1)
    gl.GetProgramInfoLog(program, logLength, nil, &log[0])
    return 0, fmt.Errorf("failed to link program: %s", log)
  }
  return program, nil
}

func compileShader(source string, shaderType uint32) (uint32, error) {
  shader := gl.CreateShader(shaderType)
  csources, free := gl.Strs(source)
  defer free()
  
  gl.ShaderSource(shader, 1, csources, nil)
  gl.CompileShader(shader)

  var status int32
  gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
  if status == gl.FALSE {
    var logLength int32
    gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
    log := make([]byte, logLength+1)
    gl.GetShaderInfoLog(shader, logLength, nil, &log[0])
    return 0, fmt.Errorf("failed to compile shader: %s", log)
  }
  return shader, nil
}

func (r *Renderer) GetMVP(aspect float32) mgl32.Mat4 {

  proj := mgl32.Perspective(mgl32.DegToRad(45), aspect, 0.1, 100.0)

  view := mgl32.Translate3D(0, 0, -r.cameraDist).
  Mul4(mgl32.HomogRotate3DX(mgl32.DegToRad(r.angleX))).
  Mul4(mgl32.HomogRotate3DY(mgl32.DegToRad(r.angleY)))

  model := mgl32.Ident4()

  return proj.Mul4(view).Mul4(model)
}

func (r *Renderer) DrawVertices(vertices []float32, color [4]float32, mode uint32, mvp mgl32.Mat4) {

  gl.UseProgram(r.program)

  data := mvp[:]
  gl.UniformMatrix4fv(r.mvpLoc, 1, false, &data[0])
  gl.Uniform4f(r.colorLoc, color[0], color[1], color[2], color[3])

  var vbo uint32
  gl.GenBuffers(1, &vbo)
  gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
  gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

  gl.EnableVertexAttribArray(r.positionLoc)
  gl.VertexAttribPointer(r.positionLoc, 3, gl.FLOAT, false, 0, nil)

  gl.DrawArrays(mode, 0, int32(len(vertices)/3))

  gl.DisableVertexAttribArray(r.positionLoc)
  gl.BindBuffer(gl.ARRAY_BUFFER, 0)
  gl.DeleteBuffers(1, &vbo)
}
