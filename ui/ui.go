package ui

import (
    "fmt"
    "image/color"
    "log"
    "math"

    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/ebitenutil"
    "tic-tac-toe/ai"
    "tic-tac-toe/animation"
    "tic-tac-toe/game"
)

const (
    screenWidth  = 600
    screenHeight = 600
)

type GameState int

const (
    StateMenu GameState = iota
    StatePlaying
    StateResult
)

type GameUI struct {
    state         GameState
    board         *game.Board
    currentPlayer game.Player
    aiOpponent    bool
    aiPlayer      game.Player
    aiEngine      ai.AI
    difficulty    string
    winner        game.Player
    animation     *animation.AnimationManager

    playerXWins   int
    playerOWins   int
    draws         int
    inputDebounce bool
}

func NewGameUI() *GameUI {
    return &GameUI{
        state:     StateMenu,
        animation: animation.NewAnimationManager(),
    }
}

func (g *GameUI) Update() error {
    // Simple input debounce
    if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
        if g.inputDebounce {
            return nil
        }
    } else {
        g.inputDebounce = false
    }

    switch g.state {
    case StateMenu:
        if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
            _, y := ebiten.CursorPosition()
            if y > 200 && y < 250 {
                g.startGame("easy")
                g.inputDebounce = true
            } else if y > 260 && y < 310 {
                g.startGame("medium")
                g.inputDebounce = true
            } else if y > 320 && y < 370 {
                g.startGame("hard")
                g.inputDebounce = true
            }
        }

    case StatePlaying:
        // If it's AI's turn, place a move in the board cell (and animate at the cell center).
        if g.aiOpponent && g.currentPlayer == g.aiPlayer {
            move := g.aiEngine.Move(g.board, g.aiPlayer)
            if err := g.board.MakeMove(move, g.aiPlayer); err != nil {
                log.Println(err)
            }
            g.animation.TriggerMoveAnimation(move) // center of the cell
            g.checkGameOver()
            g.switchPlayer()
            return nil
        }

        // Human move: place in the board, plus show an "impact" circle at the mouse
        if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
            x, y := ebiten.CursorPosition()
            pos := g.screenPosToBoardIndex(x, y)
            if pos != -1 && g.board.Cells[pos] == game.None {
                if err := g.board.MakeMove(pos, g.currentPlayer); err != nil {
                    log.Println(err)
                } else {
                    // Animate an "impact circle" at the mouse location
                    g.animation.TriggerMoveAnimationAt(x, y)

                    g.checkGameOver()
                    g.switchPlayer()
                }
                g.inputDebounce = true
            }
        }

    case StateResult:
        if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
            _, y := ebiten.CursorPosition()
            if y > 500 && y < 550 {
                g.startGame(g.difficulty)
                g.inputDebounce = true
            } else if y > 560 && y < 610 {
                g.state = StateMenu
                g.inputDebounce = true
            }
        }
    }

    // Update any active animations
    g.animation.Update()
    return nil
}

func (g *GameUI) Draw(screen *ebiten.Image) {
    switch g.state {
    case StateMenu:
        screen.Fill(color.RGBA{100, 149, 237, 255})
        ebitenutil.DebugPrintAt(screen, "Tic Tac Toe - Choose Difficulty", 200, 150)
        drawButton(screen, "Easy", 200)
        drawButton(screen, "Medium", 260)
        drawButton(screen, "Hard", 320)

    case StatePlaying:
        // Fill entire window with a background color
        screen.Fill(color.RGBA{230, 230, 250, 255}) // a light lavender

        g.drawBoard(screen)
        g.animation.Draw(screen)

    case StateResult:
        screen.Fill(color.RGBA{240, 240, 240, 255})
        msg := "Draw!"
        if g.winner == game.PlayerX {
            msg = "Player X Wins!"
        } else if g.winner == game.PlayerO {
            msg = "Player O Wins!"
        }
        ebitenutil.DebugPrintAt(screen, msg, 240, 100)
        ebitenutil.DebugPrintAt(
            screen,
            "Score - X: "+itoa(g.playerXWins)+" O: "+itoa(g.playerOWins)+" Draws: "+itoa(g.draws),
            200, 150,
        )
        drawButton(screen, "Play Again", 500)
        drawButton(screen, "Back to Menu", 560)
    }
}

func (g *GameUI) Layout(outsideWidth, outsideHeight int) (int, int) {
    // Return the same size as screenWidth and screenHeight constants
    // to ensure the window matches the game size exactly
    return screenWidth, screenHeight
}

func (g *GameUI) startGame(difficulty string) {
    g.state = StatePlaying
    g.board = game.NewBoard()
    g.currentPlayer = game.PlayerX
    g.aiOpponent = true
    g.aiPlayer = game.PlayerO
    g.aiEngine = ai.NewAI(difficulty)
    g.difficulty = difficulty
    g.animation.Reset()
    g.winner = game.None
}

// drawBoard uses thick lines for the grid and draws X/O with thicker strokes.
func (g *GameUI) drawBoard(screen *ebiten.Image) {
    cellSize := screenWidth / 3

    // 1) Draw thick grid lines using rectangles
    lineColor := color.RGBA{60, 60, 60, 255} // dark gray
    lineThickness := 6.0

    // Vertical lines
    for i := 1; i <= 2; i++ {
        x := float64(i * cellSize)
        ebitenutil.DrawRect(screen, x-(lineThickness/2), 0, lineThickness, float64(cellSize*3), lineColor)
    }

    // Horizontal lines
    for i := 1; i <= 2; i++ {
        y := float64(i * cellSize)
        ebitenutil.DrawRect(screen, 0, y-(lineThickness/2), float64(cellSize*3), lineThickness, lineColor)
    }

    // 2) Draw the X and O
    for i, cell := range g.board.Cells {
        x, y := g.boardIndexToScreenPos(i)
        switch cell {
        case game.PlayerX:
            // Draw thick X
            drawThickLine(screen,
                float64(x+20), float64(y+20),
                float64(x+cellSize-20), float64(y+cellSize-20),
                6, color.RGBA{220, 20, 60, 255}) // crimson

            drawThickLine(screen,
                float64(x+cellSize-20), float64(y+20),
                float64(x+20), float64(y+cellSize-20),
                6, color.RGBA{220, 20, 60, 255})

        case game.PlayerO:
            centerX := float64(x + cellSize/2)
            centerY := float64(y + cellSize/2)
            radius := float64(cellSize/2 - 20)
            drawThickCircle(screen, centerX, centerY, radius, 6, color.RGBA{30, 144, 255, 255}) // dodger blue
        }
    }
}

func (g *GameUI) screenPosToBoardIndex(x, y int) int {
    cellWidth := screenWidth / 3
    // Only consider the top 600x600 for the board
    if x < 0 || x >= cellWidth*3 || y < 0 || y >= cellWidth*3 {
        return -1
    }
    col := x / cellWidth
    row := y / cellWidth
    return row*3 + col
}

func (g *GameUI) boardIndexToScreenPos(index int) (int, int) {
    cellWidth := screenWidth / 3
    row := index / 3
    col := index % 3
    return col * cellWidth, row * cellWidth
}

func (g *GameUI) switchPlayer() {
    if g.currentPlayer == game.PlayerX {
        g.currentPlayer = game.PlayerO
    } else {
        g.currentPlayer = game.PlayerX
    }
}

func (g *GameUI) checkGameOver() {
    if winner, ok := g.board.CheckWin(); ok {
        g.winner = winner
        g.state = StateResult
        if winner == game.PlayerX {
            g.playerXWins++
        } else {
            g.playerOWins++
        }
    } else if g.board.IsDraw() {
        g.winner = game.None
        g.state = StateResult
        g.draws++
    }
}

// -------------------- Thick Drawing Helpers --------------------

// drawThickLine draws a line as a rotated rectangle.
func drawThickLine(screen *ebiten.Image, x1, y1, x2, y2, thickness float64, clr color.Color) {
    dx := x2 - x1
    dy := y2 - y1
    length := math.Sqrt(dx*dx + dy*dy)

    // Create a small rectangle of width=length and height=thickness
    lineImage := ebiten.NewImage(int(length), int(thickness))
    lineImage.Fill(clr)

    op := &ebiten.DrawImageOptions{}
    // Shift origin so rotation is around rectangle center
    op.GeoM.Translate(-length/2, -thickness/2)
    // Rotate
    angle := math.Atan2(dy, dx)
    op.GeoM.Rotate(angle)
    // Move it to midpoint
    midX := (x1 + x2) / 2
    midY := (y1 + y2) / 2
    op.GeoM.Translate(midX, midY)

    screen.DrawImage(lineImage, op)
}

// drawThickCircle approximates a circle with multiple thick line segments.
func drawThickCircle(screen *ebiten.Image, cx, cy, r, thickness float64, clr color.Color) {
    step := 5.0 // in degrees
    for angle := 0.0; angle < 360.0; angle += step {
        x1 := cx + r*math.Cos(angle*math.Pi/180)
        y1 := cy + r*math.Sin(angle*math.Pi/180)
        x2 := cx + r*math.Cos((angle+step)*math.Pi/180)
        y2 := cy + r*math.Sin((angle+step)*math.Pi/180)
        drawThickLine(screen, x1, y1, x2, y2, thickness, clr)
    }
}

// drawButton draws a simple rectangular button with text.
func drawButton(screen *ebiten.Image, text string, y int) {
    buttonColor := color.RGBA{173, 216, 230, 255}
    ebitenutil.DrawRect(screen, 200, float64(y), 200, 40, buttonColor)
    ebitenutil.DebugPrintAt(screen, text, 270, y+10)
}

// itoa is a helper to convert int to string.
func itoa(n int) string {
    return fmt.Sprintf("%d", n)
}