package game

import "errors"

// Player represents a player (or empty cell).
type Player int

const (
  None Player = iota
  PlayerX
  PlayerO
)

// Board represents a Tic Tac Toe board.
type Board struct {
  Cells [9]Player
}

// NewBoard creates a new empty board.
func NewBoard() *Board {
  return &Board{}
}

// MakeMove attempts to mark the cell at position pos with the given player's symbol.
func (b *Board) MakeMove(pos int, player Player) error {
  if pos < 0 || pos >= len(b.Cells) {
    return errors.New("invalid move: position out of bounds")
  }
  if b.Cells[pos] != None {
    return errors.New("invalid move: cell already occupied")
  }
  b.Cells[pos] = player
  return nil
}

// CheckWin checks if there is a winning combination on the board.
// Returns the winning player and true if there is a win.
func (b *Board) CheckWin() (Player, bool) {
  wins := [][3]int{
    {0, 1, 2}, {3, 4, 5}, {6, 7, 8}, // rows
    {0, 3, 6}, {1, 4, 7}, {2, 5, 8}, // columns
    {0, 4, 8}, {2, 4, 6},            // diagonals
  }
  for _, line := range wins {
    if b.Cells[line[0]] != None &&
      b.Cells[line[0]] == b.Cells[line[1]] &&
      b.Cells[line[1]] == b.Cells[line[2]] {
      return b.Cells[line[0]], true
    }
  }
  return None, false
}

// IsDraw returns true if the board is full and no winning move exists.
func (b *Board) IsDraw() bool {
  for _, cell := range b.Cells {
    if cell == None {
      return false
    }
  }
  return true
}
