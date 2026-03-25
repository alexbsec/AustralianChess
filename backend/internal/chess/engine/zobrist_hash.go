package engine

import (
	"math/rand"

	"github.com/alexbsec/AustralianChess/backend/internal/chess"
)


const (
    numPieceKinds  = 8  
    numColors      = 2
    numSquares     = 144 
)

var (
    zobristPiece [numPieceKinds][numColors][numSquares]uint64
    zobristTurn  uint64
)

func init() {
    rng := rand.New(rand.NewSource(0x54415950414e))
    for k := range numPieceKinds {
        for c := range numColors {
            for sq := range numSquares {
                zobristPiece[k][c][sq] = rng.Uint64()
            }
        }
    }
    zobristTurn = rng.Uint64()
}

func ZobristHash(board chess.Board, color chess.PieceColor) uint64 {
    var h uint64
    rows, cols := board.Dimensions()
    for row := range rows {
        for col := range cols {
            pos := chess.Position{Row: row, Col: col}
            piece := board.GetPiece(pos)
            if piece != nil {
                h ^= zobristPiece[piece.Kind][piece.Color][squareIndex(pos)]
            }
        }
    }
    if color == chess.PieceBlack {
        h ^= zobristTurn
    }
    return h
}

func ZobristUpdate(hash uint64, board chess.Board, from, to chess.Position) uint64 {
    piece := board.GetPiece(from)
    if piece == nil {
        return hash
    }

    // XOR out piece at source
    hash ^= zobristPiece[piece.Kind][piece.Color][squareIndex(from)]

    // XOR out captured piece at destination if any
    captured := board.GetPiece(to)
    if captured != nil {
        hash ^= zobristPiece[captured.Kind][captured.Color][squareIndex(to)]
    }

    // XOR in piece at destination
    hash ^= zobristPiece[piece.Kind][piece.Color][squareIndex(to)]

    // Flip turn
    hash ^= zobristTurn

    return hash
}


func squareIndex(pos chess.Position) int {
    return pos.Row*12 + pos.Col
}
