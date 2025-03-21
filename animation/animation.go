package animation

import (
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/ebitenutil"
    "image/color"
    "math"
    "time"
)

// Animation represents a single animation instance.
type Animation struct {
    boardPosition int           // valid if isPixel=false
    pixelX        int           // valid if isPixel=true
    pixelY        int           // valid if isPixel=true
    isPixel       bool          // whether we are drawing at board coords or pixel coords
    startTime     time.Time
    duration      time.Duration
}

// AnimationManager manages active animations.
type AnimationManager struct {
    animations []Animation
}

// NewAnimationManager creates a new animation manager.
func NewAnimationManager() *AnimationManager {
    return &AnimationManager{
        animations: []Animation{},
    }
}

// TriggerMoveAnimation adds a new animation at the center of the specified board cell.
func (am *AnimationManager) TriggerMoveAnimation(position int) {
    am.animations = append(am.animations, Animation{
        boardPosition: position,
        isPixel:       false,
        startTime:     time.Now(),
        duration:      300 * time.Millisecond,
    })
}

// TriggerMoveAnimationAt adds a new animation at specific pixel coordinates (mouse impact).
func (am *AnimationManager) TriggerMoveAnimationAt(x int, y int) {
    am.animations = append(am.animations, Animation{
        pixelX:    x,
        pixelY:    y,
        isPixel:   true,
        startTime: time.Now(),
        duration:  300 * time.Millisecond,
    })
}

// Update removes expired animations.
func (am *AnimationManager) Update() {
    now := time.Now()
    var remaining []Animation
    for _, anim := range am.animations {
        if now.Sub(anim.startTime) < anim.duration {
            remaining = append(remaining, anim)
        }
    }
    am.animations = remaining
}

// Draw renders active animations on the screen.
func (am *AnimationManager) Draw(screen *ebiten.Image) {
    now := time.Now()
    for _, anim := range am.animations {
        elapsed := now.Sub(anim.startTime)
        ratio := float64(elapsed) / float64(anim.duration)
        alpha := uint8(255 - uint8(ratio*255))

        // Gold color that fades out
        clr := color.RGBA{255, 215, 0, alpha}

        if anim.isPixel {
            // Draw at mouse pixel coordinates
            cx := float64(anim.pixelX)
            cy := float64(anim.pixelY)
            r := 20.0
            drawCircle(screen, cx, cy, r, clr)
        } else {
            // Draw in the center of a board cell
            cellWidth := screen.Bounds().Dx() / 3
            cellHeight := screen.Bounds().Dy() / 3
            row := anim.boardPosition / 3
            col := anim.boardPosition % 3
            cx := float64(col*cellWidth + cellWidth/2)
            cy := float64(row*cellHeight + cellHeight/2)
            r := float64(cellWidth) / 4
            drawCircle(screen, cx, cy, r, clr)
        }
    }
}

// drawCircle approximates a circle with short lines.
func drawCircle(screen *ebiten.Image, cx, cy, r float64, clr color.Color) {
    step := 5.0
    for angle := 0.0; angle < 360.0; angle += step {
        x1 := cx + r*math.Cos(angle*math.Pi/180)
        y1 := cy + r*math.Sin(angle*math.Pi/180)
        x2 := cx + r*math.Cos((angle+step)*math.Pi/180)
        y2 := cy + r*math.Sin((angle+step)*math.Pi/180)
        ebitenutil.DrawLine(screen, x1, y1, x2, y2, clr)
    }
}

// Reset clears all active animations.
func (am *AnimationManager) Reset() {
    am.animations = nil
}