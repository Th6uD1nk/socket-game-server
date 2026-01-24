package core

import "sync"

const (
  LPAD_UP    = 0
  LPAD_DOWN  = 1
  LPAD_LEFT  = 2
  LPAD_RIGHT = 3

  RPAD_UP     = 10
  RPAD_DOWN   = 11
  RPAD_LEFT   = 12
  RPAD_RIGHT  = 13
  //
  RPAD_LCLICK = 14
  RPAD_RCLICK = 15
)

type InputManager struct {
  mu sync.RWMutex
  
  padStates     [20]bool
  padPrevStates [20]bool
  
  // pc
  mouseX      float64
  mouseY      float64
  mousePrevX  float64
  mousePrevY  float64
  mouseMoved  bool
  lpadKeys    [4]int
  
  // mobile
  leftPadTouch  *TouchState
  rightPadTouch *TouchState
  leftZone      PadZone
  rightZone     PadZone
}

type PadZone struct {
  Left   float32
  Top    float32
  Right  float32
  Bottom float32
}

type TouchState struct {
  Active bool
  X      float32
  Y      float32
  StartX float32
  StartY float32
}

var inputMgr *InputManager
var once sync.Once

func NewInputManager() *InputManager {
  once.Do(func() {
    inputMgr = &InputManager{
      leftPadTouch:  &TouchState{},
      rightPadTouch: &TouchState{},
    }
  })
  return inputMgr
}

func GetInputManager() *InputManager {
  if inputMgr == nil {
    return NewInputManager()
  }
  return inputMgr
}

/**
 * global functions to get pad states from game
 */
func GetPad(padCode int) bool {
  inputMgr.mu.RLock()
  defer inputMgr.mu.RUnlock()
  
  if padCode < 0 || padCode >= len(inputMgr.padStates) {
    return false
  }
  return inputMgr.padStates[padCode]
}

func GetPadDown(padCode int) bool {
  inputMgr.mu.RLock()
  defer inputMgr.mu.RUnlock()
  
  if padCode < 0 || padCode >= len(inputMgr.padStates) {
    return false
  }
  return inputMgr.padStates[padCode] && !inputMgr.padPrevStates[padCode]
}

func GetPadUp(padCode int) bool {
  inputMgr.mu.RLock()
  defer inputMgr.mu.RUnlock()
  
  if padCode < 0 || padCode >= len(inputMgr.padStates) {
    return false
  }
  return !inputMgr.padStates[padCode] && inputMgr.padPrevStates[padCode]
}

/**
 * pc bindings functions
 */
func BindLPadKeys(up, down, left, right int) {
  inputMgr.mu.Lock()
  defer inputMgr.mu.Unlock()
  
  inputMgr.lpadKeys[0] = up
  inputMgr.lpadKeys[1] = down
  inputMgr.lpadKeys[2] = left
  inputMgr.lpadKeys[3] = right
}

func ProcessKeyboard(key int, pressed bool) {
  inputMgr.mu.Lock()
  defer inputMgr.mu.Unlock()

  for i, boundKey := range inputMgr.lpadKeys {
    if key == boundKey {
      inputMgr.padStates[i] = pressed
      return
    }
  }
}

func ProcessMouseButton(button int, pressed bool) {
  inputMgr.mu.Lock()
  defer inputMgr.mu.Unlock()
  
  if button == 0 {
    inputMgr.padStates[RPAD_LCLICK] = pressed
  } else if button == 1 {
    inputMgr.padStates[RPAD_RCLICK] = pressed
  }
}

func ProcessMouseMove(x, y float64) {
  inputMgr.mu.Lock()
  defer inputMgr.mu.Unlock()
  
  inputMgr.mousePrevX = inputMgr.mouseX
  inputMgr.mousePrevY = inputMgr.mouseY
  inputMgr.mouseX = x
  inputMgr.mouseY = y
  inputMgr.mouseMoved = true
}

/**
 * mobile bindings functions
 */
 func SetLeftPadZone(left, top, right, bottom float32) {
  inputMgr.mu.Lock()
  defer inputMgr.mu.Unlock()
  
  inputMgr.leftZone = PadZone{
    Left:   left,
    Top:    top,
    Right:  right,
    Bottom: bottom,
  }
}

func SetRightPadZone(left, top, right, bottom float32) {
  inputMgr.mu.Lock()
  defer inputMgr.mu.Unlock()
  
  inputMgr.rightZone = PadZone{
    Left:   left,
    Top:    top,
    Right:  right,
    Bottom: bottom,
  }
}

func ProcessTouch(x, y float32, pressed bool, touchID int) {
  inputMgr.mu.Lock()
  defer inputMgr.mu.Unlock()
  
  inLeftZone := x >= inputMgr.leftZone.Left && x <= inputMgr.leftZone.Right &&
    y >= inputMgr.leftZone.Top && y <= inputMgr.leftZone.Bottom
  
  inRightZone := x >= inputMgr.rightZone.Left && x <= inputMgr.rightZone.Right &&
    y >= inputMgr.rightZone.Top && y <= inputMgr.rightZone.Bottom
  
  if inLeftZone {
    processLeftPadTouch(x, y, pressed)
  } else if inRightZone {
    processRightPadTouch(x, y, pressed)
  }
}

func processLeftPadTouch(x, y float32, pressed bool) {
  if !pressed {
    inputMgr.leftPadTouch.Active = false
    inputMgr.padStates[LPAD_UP] = false
    inputMgr.padStates[LPAD_DOWN] = false
    inputMgr.padStates[LPAD_LEFT] = false
    inputMgr.padStates[LPAD_RIGHT] = false
    return
  }
  
  if !inputMgr.leftPadTouch.Active {
    inputMgr.leftPadTouch.Active = true
    inputMgr.leftPadTouch.StartX = x
    inputMgr.leftPadTouch.StartY = y
  }
  
  inputMgr.leftPadTouch.X = x
  inputMgr.leftPadTouch.Y = y
  
  dx := x - inputMgr.leftPadTouch.StartX
  dy := y - inputMgr.leftPadTouch.StartY
  
  threshold := float32(20.0) // tmp
  
  // v
  if dy < -threshold {
    inputMgr.padStates[LPAD_UP] = true
    inputMgr.padStates[LPAD_DOWN] = false
  } else if dy > threshold {
    inputMgr.padStates[LPAD_DOWN] = true
    inputMgr.padStates[LPAD_UP] = false
  } else {
    inputMgr.padStates[LPAD_UP] = false
    inputMgr.padStates[LPAD_DOWN] = false
  }
  
  // h
  if dx < -threshold {
    inputMgr.padStates[LPAD_LEFT] = true
    inputMgr.padStates[LPAD_RIGHT] = false
  } else if dx > threshold {
    inputMgr.padStates[LPAD_RIGHT] = true
    inputMgr.padStates[LPAD_LEFT] = false
  } else {
    inputMgr.padStates[LPAD_LEFT] = false
    inputMgr.padStates[LPAD_RIGHT] = false
  }
}

func processRightPadTouch(x, y float32, pressed bool) {
  if !pressed {
    inputMgr.rightPadTouch.Active = false
    inputMgr.padStates[RPAD_UP] = false
    inputMgr.padStates[RPAD_DOWN] = false
    inputMgr.padStates[RPAD_LEFT] = false
    inputMgr.padStates[RPAD_RIGHT] = false
    inputMgr.padStates[RPAD_LCLICK] = false
    return
  }
  
  if !inputMgr.rightPadTouch.Active {
    inputMgr.rightPadTouch.Active = true
    inputMgr.rightPadTouch.StartX = x
    inputMgr.rightPadTouch.StartY = y
    inputMgr.padStates[RPAD_LCLICK] = true
  }
  
  inputMgr.rightPadTouch.X = x
  inputMgr.rightPadTouch.Y = y
  
  dx := x - inputMgr.rightPadTouch.StartX
  dy := y - inputMgr.rightPadTouch.StartY
  
  threshold := float32(20.0)
  
  if dy < -threshold {
    inputMgr.padStates[RPAD_UP] = true
    inputMgr.padStates[RPAD_DOWN] = false
  } else if dy > threshold {
    inputMgr.padStates[RPAD_DOWN] = true
    inputMgr.padStates[RPAD_UP] = false
  } else {
    inputMgr.padStates[RPAD_UP] = false
    inputMgr.padStates[RPAD_DOWN] = false
  }
  
  if dx < -threshold {
    inputMgr.padStates[RPAD_LEFT] = true
    inputMgr.padStates[RPAD_RIGHT] = false
  } else if dx > threshold {
    inputMgr.padStates[RPAD_RIGHT] = true
    inputMgr.padStates[RPAD_LEFT] = false
  } else {
    inputMgr.padStates[RPAD_LEFT] = false
    inputMgr.padStates[RPAD_RIGHT] = false
  }
}

/**
 * update functions
 */
func UpdateInputs() {
  inputMgr.mu.Lock()
  defer inputMgr.mu.Unlock()
  
  copy(inputMgr.padPrevStates[:], inputMgr.padStates[:])
  
  if inputMgr.mouseMoved {
    dx := inputMgr.mouseX - inputMgr.mousePrevX
    dy := inputMgr.mouseY - inputMgr.mousePrevY
    
    threshold := 5.0 // tmp
    
    if dy < -threshold {
      inputMgr.padStates[RPAD_UP] = true
      inputMgr.padStates[RPAD_DOWN] = false
    } else if dy > threshold {
      inputMgr.padStates[RPAD_DOWN] = true
      inputMgr.padStates[RPAD_UP] = false
    } else {
      inputMgr.padStates[RPAD_UP] = false
      inputMgr.padStates[RPAD_DOWN] = false
    }
    
    if dx < -threshold {
      inputMgr.padStates[RPAD_LEFT] = true
      inputMgr.padStates[RPAD_RIGHT] = false
    } else if dx > threshold {
      inputMgr.padStates[RPAD_RIGHT] = true
      inputMgr.padStates[RPAD_LEFT] = false
    } else {
      inputMgr.padStates[RPAD_LEFT] = false
      inputMgr.padStates[RPAD_RIGHT] = false
    }
    
    inputMgr.mouseMoved = false
  } else {
    inputMgr.padStates[RPAD_UP] = false
    inputMgr.padStates[RPAD_DOWN] = false
    inputMgr.padStates[RPAD_LEFT] = false
    inputMgr.padStates[RPAD_RIGHT] = false
  }
}
