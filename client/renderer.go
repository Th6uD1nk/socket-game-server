package main

import (
  "image/color"
  "math"
  "github.com/hajimehoshi/ebiten/v2"
  "github.com/hajimehoshi/ebiten/v2/vector"
)

type Renderer struct {
  angleX     float64
  angleY     float64
  fov        float64
  cameraDist float64
}

func NewRenderer() *Renderer {
  return &Renderer{
    angleX:     math.Pi / 180.0 * 45.0,
    angleY:     math.Pi / 180.0 * 10.0,
    fov:        400.0,
    cameraDist: 20.0,
  }
}

func (r *Renderer) project(v Vec3, screenW, screenH float64) (float32, float32) {
  xh := v.X*math.Cos(r.angleY) - v.Z*math.Sin(r.angleY)
  zh := v.X*math.Sin(r.angleY) + v.Z*math.Cos(r.angleY)

  yh := v.Y
  ax := -r.angleX

  yv := yh*math.Cos(ax) - zh*math.Sin(ax)
  zv := yh*math.Sin(ax) + zh*math.Cos(ax)

  depth := zv + r.cameraDist
  if depth <= 0.1 {
    depth = 0.1
  }

  factor := r.fov / depth

  px := xh*factor + screenW/2
  py := -yv*factor + screenH/2

  return float32(px), float32(py)
}

func (r *Renderer) DrawGrid(screen *ebiten.Image, gridSize int) {
  w, h := screen.Size()
  screenW, screenH := float64(w), float64(h)
  
  gridColor := color.RGBA{60, 60, 80, 255}
  
  for i := -gridSize; i <= gridSize; i++ {
    x1, y1 := r.project(Vec3{X: float64(i), Y: 0, Z: float64(-gridSize)}, screenW, screenH)
    x2, y2 := r.project(Vec3{X: float64(i), Y: 0, Z: float64(gridSize)}, screenW, screenH)
    vector.StrokeLine(screen, x1, y1, x2, y2, 1, gridColor, false)
    
    z1, w1 := r.project(Vec3{X: float64(-gridSize), Y: 0, Z: float64(i)}, screenW, screenH)
    z2, w2 := r.project(Vec3{X: float64(gridSize), Y: 0, Z: float64(i)}, screenW, screenH)
    vector.StrokeLine(screen, z1, w1, z2, w2, 1, gridColor, false)
  }
}

func (r *Renderer) DrawCube(screen *ebiten.Image, pos Vec3, cubeColor color.RGBA) {
  w, h := screen.Size()
  screenW, screenH := float64(w), float64(h)
  
  cubeVertices := []Vec3{
    {pos.X - 0.5, pos.Y - 0.5, pos.Z - 0.5}, // bottom front left
    {pos.X + 0.5, pos.Y - 0.5, pos.Z - 0.5}, // bottom front right
    {pos.X + 0.5, pos.Y - 0.5, pos.Z + 0.5}, // bottom back right
    {pos.X - 0.5, pos.Y - 0.5, pos.Z + 0.5}, // bottom back left
    {pos.X - 0.5, pos.Y + 0.5, pos.Z - 0.5}, // top front left
    {pos.X + 0.5, pos.Y + 0.5, pos.Z - 0.5}, // top front right
    {pos.X + 0.5, pos.Y + 0.5, pos.Z + 0.5}, // top back right
    {pos.X - 0.5, pos.Y + 0.5, pos.Z + 0.5}, // top back left
  }
  
  var projected [][2]float32
  for _, v := range cubeVertices {
    px, py := r.project(v, screenW, screenH)
    projected = append(projected, [2]float32{px, py})
  }
  
  faces := [][4]int{
    {0, 1, 5, 4}, // front
    {1, 2, 6, 5}, // right
    {2, 3, 7, 6}, // back
    {3, 0, 4, 7}, // left
    {4, 5, 6, 7}, // top
    {0, 1, 2, 3}, // bottom
  }
  
  faceColor := color.RGBA{cubeColor.R / 2, cubeColor.G / 2, cubeColor.B / 2, 200}
  
  for _, face := range faces {
    var path vector.Path
    path.MoveTo(projected[face[0]][0], projected[face[0]][1])
    for i := 1; i < 4; i++ {
      path.LineTo(projected[face[i]][0], projected[face[i]][1])
    }
    path.Close()
    
    vertices, indices := path.AppendVerticesAndIndicesForFilling(nil, nil)
    for i := range vertices {
      vertices[i].ColorR = float32(faceColor.R) / 255
      vertices[i].ColorG = float32(faceColor.G) / 255
      vertices[i].ColorB = float32(faceColor.B) / 255
      vertices[i].ColorA = float32(faceColor.A) / 255
    }
    screen.DrawTriangles(vertices, indices, emptyImage, nil)
  }
  
  edges := [][2]int{
    {0, 1}, {1, 2}, {2, 3}, {3, 0},
    {4, 5}, {5, 6}, {6, 7}, {7, 4},
    {0, 4}, {1, 5}, {2, 6}, {3, 7},
  }
  
  for _, edge := range edges {
    p1 := projected[edge[0]]
    p2 := projected[edge[1]]
    vector.StrokeLine(screen, p1[0], p1[1], p2[0], p2[1], 1, cubeColor, false)
  }
}

var emptyImage = ebiten.NewImage(1, 1)

func init() {
  emptyImage.Fill(color.White)
}
