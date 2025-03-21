package ai

import (
  "math/rand"
  "time"

  "tic-tac-toe/game"
)

// AI defines the interface for AI opponents.
type AI interface {
  Move(board *game.Board, player game.Player) int
}

// EasyAI chooses a random available move.
type EasyAI struct{}

// Move selects a random move from available moves.
func (ai EasyAI) Move(board *game.Board, player game.Player) int {
  var moves []int
  for i, cell := range board.Cells {
    if cell == game.None {
      moves = append(moves, i)
    }
  }
  if len(moves) == 0 {
    return -1 // no moves available
  }
  rand.Seed(time.Now().UnixNano())
  return moves[rand.Intn(len(moves))]
}

// MediumAI checks for winning moves or blocks opponent wins.
type MediumAI struct{}

// Move checks if the AI can win or block before making a random move.
func (ai MediumAI) Move(board *game.Board, player game.Player) int {
  // Check if AI can win in the next move.
  for _, move := range availableMoves(board) {
    board.Cells[move] = player
    if winner, ok := board.CheckWin(); ok && winner == player {
      board.Cells[move] = game.None
      return move
    }
    board.Cells[move] = game.None
  }
  // Block opponent win.
  opponent := game.PlayerX
  if player == game.PlayerX {
    opponent = game.PlayerO
  }
  for _, move := range availableMoves(board) {
    board.Cells[move] = opponent
    if winner, ok := board.CheckWin(); ok && winner == opponent {
      board.Cells[move] = game.None
      return move
    }
    board.Cells[move] = game.None
  }
  // Otherwise, choose a random move.
  return EasyAI{}.Move(board, player)
}

// HardAI uses the minimax algorithm.
type HardAI struct{}

// Move selects the best move using minimax.
func (ai HardAI) Move(board *game.Board, player game.Player) int {
  bestScore := -1000
  bestMove := -1
  for _, move := range availableMoves(board) {
    board.Cells[move] = player
    score := minimax(board, false, player)
    board.Cells[move] = game.None
    if score > bestScore {
      bestScore = score
      bestMove = move
    }
  }
  return bestMove
}

// minimax recursively evaluates board positions.
func minimax(board *game.Board, isMaximizing bool, aiPlayer game.Player) int {
  if winner, ok := board.CheckWin(); ok {
    if winner == aiPlayer {
      return 10
    }
    return -10
  }
  if board.IsDraw() {
    return 0
  }

  if isMaximizing {
    bestScore := -1000
    for _, move := range availableMoves(board) {
      board.Cells[move] = aiPlayer
      score := minimax(board, false, aiPlayer)
      board.Cells[move] = game.None
      if score > bestScore {
        bestScore = score
      }
    }
    return bestScore
  } else {
    // Determine the opponent.
    opponent := game.PlayerX
    if aiPlayer == game.PlayerX {
      opponent = game.PlayerO
    }
    bestScore := 1000
    for _, move := range availableMoves(board) {
      board.Cells[move] = opponent
      score := minimax(board, true, aiPlayer)
      board.Cells[move] = game.None
      if score < bestScore {
        bestScore = score
      }
    }
    return bestScore
  }
}

// availableMoves returns a slice of indices for cells that are not yet occupied.
func availableMoves(board *game.Board) []int {
  var moves []int
  for i, cell := range board.Cells {
    if cell == game.None {
      moves = append(moves, i)
    }
  }
  return moves
}

// NewAI returns an AI instance based on the chosen difficulty.
func NewAI(difficulty string) AI {
  switch difficulty {
  case "easy":
    return EasyAI{}
  case "medium":
    return MediumAI{}
  case "hard":
    return HardAI{}
  default:
    return EasyAI{}
  }
}