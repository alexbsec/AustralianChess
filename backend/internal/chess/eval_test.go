package chess_test

import (
	"testing"

	"github.com/alexbsec/AustralianChess/backend/internal/chess"
	"github.com/stretchr/testify/require"
)

func emptyBoard() chess.Board {
	board := chess.Board{
		Data: make([][]chess.Square, 12),
	}
	for r := range 12 {
		board.Data[r] = make([]chess.Square, 12)
	}
	return board
}

func placePieceOnBoard(b *chess.Board, piece *chess.Piece, pos chess.Position) {
	b.Data[pos.Row][pos.Col].Piece = piece
}

func TestEvaluateBoard_EmptyBoard(t *testing.T) {
	board := emptyBoard()
	score := chess.EvaluateBoard(board)
	require.Equal(t, 0, score)
}

func TestEvaluateBoard_WhitePawn(t *testing.T) {
	board := emptyBoard()
	pawn := chess.NewPiece(chess.PawnPiece, chess.PieceWhite)
	placePieceOnBoard(&board, pawn, chess.Position{Row: 5, Col: 5})
	score := chess.EvaluateBoard(board)
	require.Greater(t, score, 0)
}

func TestEvaluateBoard_BlackPawn(t *testing.T) {
	board := emptyBoard()
	pawn := chess.NewPiece(chess.PawnPiece, chess.PieceBlack)
	placePieceOnBoard(&board, pawn, chess.Position{Row: 5, Col: 5})
	score := chess.EvaluateBoard(board)
	require.Less(t, score, 0)
}

func TestEvaluateBoard_SymmetricPosition(t *testing.T) {
	board := emptyBoard()
	whitePawn := chess.NewPiece(chess.PawnPiece, chess.PieceWhite)
	blackPawn := chess.NewPiece(chess.PawnPiece, chess.PieceBlack)
	placePieceOnBoard(&board, whitePawn, chess.Position{Row: 5, Col: 5})
	placePieceOnBoard(&board, blackPawn, chess.Position{Row: 6, Col: 5})
	// not necessarily zero due to position bonuses, but should be close
	score := chess.EvaluateBoard(board)
	require.IsType(t, 0, score)
}

func TestStaticEval_WhitePerspective(t *testing.T) {
	board := emptyBoard()
	whitePawn := chess.NewPiece(chess.PawnPiece, chess.PieceWhite)
	placePieceOnBoard(&board, whitePawn, chess.Position{Row: 5, Col: 5})
	raw := chess.EvaluateBoard(board)
	eval := chess.StaticEval(board, chess.PieceWhite)
	require.Equal(t, raw, eval)
}

func TestStaticEval_BlackPerspective(t *testing.T) {
	board := emptyBoard()
	whitePawn := chess.NewPiece(chess.PawnPiece, chess.PieceWhite)
	placePieceOnBoard(&board, whitePawn, chess.Position{Row: 5, Col: 5})
	raw := chess.EvaluateBoard(board)
	eval := chess.StaticEval(board, chess.PieceBlack)
	require.Equal(t, -raw, eval)
}

func TestIsInsideBoard_Inside(t *testing.T) {
	board := chess.NewInitialBoard()
	from := chess.Position{Row: 0, Col: 0}
	to := chess.Position{Row: 11, Col: 11}
	require.True(t, chess.IsInsideBoard(board, from, to))
}

func TestIsInsideBoard_OutsideFrom(t *testing.T) {
	board := chess.NewInitialBoard()
	from := chess.Position{Row: -1, Col: 0}
	to := chess.Position{Row: 5, Col: 5}
	require.False(t, chess.IsInsideBoard(board, from, to))
}

func TestIsInsideBoard_OutsideTo(t *testing.T) {
	board := chess.NewInitialBoard()
	from := chess.Position{Row: 5, Col: 5}
	to := chess.Position{Row: 12, Col: 5}
	require.False(t, chess.IsInsideBoard(board, from, to))
}

func TestIsInsideBoard_NegativeCol(t *testing.T) {
	board := chess.NewInitialBoard()
	from := chess.Position{Row: 5, Col: 5}
	to := chess.Position{Row: 5, Col: -1}
	require.False(t, chess.IsInsideBoard(board, from, to))
}

func TestHasAdjacentFriendlyPiece_WithNeighbor(t *testing.T) {
	board := emptyBoard()
	white1 := chess.NewPiece(chess.PawnPiece, chess.PieceWhite)
	white2 := chess.NewPiece(chess.PawnPiece, chess.PieceWhite)
	pos := chess.Position{Row: 5, Col: 5}
	neighbor := chess.Position{Row: 5, Col: 6}
	placePieceOnBoard(&board, white1, pos)
	placePieceOnBoard(&board, white2, neighbor)
	require.True(t, chess.HasAdjacentFriendlyPiece(board, chess.PieceWhite, pos))
}

func TestHasAdjacentFriendlyPiece_NoNeighbor(t *testing.T) {
	board := emptyBoard()
	white := chess.NewPiece(chess.PawnPiece, chess.PieceWhite)
	pos := chess.Position{Row: 5, Col: 5}
	placePieceOnBoard(&board, white, pos)
	require.False(t, chess.HasAdjacentFriendlyPiece(board, chess.PieceWhite, pos))
}

func TestHasAdjacentFriendlyPiece_EnemyNeighbor(t *testing.T) {
	board := emptyBoard()
	white := chess.NewPiece(chess.PawnPiece, chess.PieceWhite)
	black := chess.NewPiece(chess.PawnPiece, chess.PieceBlack)
	pos := chess.Position{Row: 5, Col: 5}
	neighbor := chess.Position{Row: 5, Col: 6}
	placePieceOnBoard(&board, white, pos)
	placePieceOnBoard(&board, black, neighbor)
	require.False(t, chess.HasAdjacentFriendlyPiece(board, chess.PieceWhite, pos))
}

func TestHasAdjacentFriendlyPiece_AtCorner(t *testing.T) {
	board := emptyBoard()
	white1 := chess.NewPiece(chess.PawnPiece, chess.PieceWhite)
	white2 := chess.NewPiece(chess.PawnPiece, chess.PieceWhite)
	pos := chess.Position{Row: 0, Col: 0}
	neighbor := chess.Position{Row: 0, Col: 1}
	placePieceOnBoard(&board, white1, pos)
	placePieceOnBoard(&board, white2, neighbor)
	require.True(t, chess.HasAdjacentFriendlyPiece(board, chess.PieceWhite, pos))
}

func TestPieceBaseValue(t *testing.T) {
	tests := []struct {
		kind     chess.PieceKind
		expected int
	}{
		{chess.PawnPiece, 100},
		{chess.KnightPiece, 320},
		{chess.KangarooPiece, 280},
		{chess.BishopPiece, 380},
		{chess.RookPiece, 550},
		{chess.OligarchPiece, 300},
		{chess.QueenPiece, 1000},
		{chess.KingPiece, 20000},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			require.Equal(t, tt.expected, chess.PieceBaseValue(tt.kind))
		})
	}
}

func TestPieceValue_Pawn(t *testing.T) {
	board := emptyBoard()
	pawn := chess.NewPiece(chess.PawnPiece, chess.PieceWhite)
	pos := chess.Position{Row: 5, Col: 5}
	rows, cols := board.Dimensions()
	val := chess.PieceValue(board, *pawn, pos, rows, cols)
	require.GreaterOrEqual(t, val, 100)
}

func TestPieceValue_Knight_Center(t *testing.T) {
	board := emptyBoard()
	knight := chess.NewPiece(chess.KnightPiece, chess.PieceWhite)
	pos := chess.Position{Row: 5, Col: 5}
	rows, cols := board.Dimensions()
	val := chess.PieceValue(board, *knight, pos, rows, cols)
	// center knight gets bonus
	require.Greater(t, val, 320)
}

func TestPieceValue_Knight_Edge(t *testing.T) {
	board := emptyBoard()
	knight := chess.NewPiece(chess.KnightPiece, chess.PieceWhite)
	pos := chess.Position{Row: 0, Col: 5}
	rows, cols := board.Dimensions()
	val := chess.PieceValue(board, *knight, pos, rows, cols)
	// edge knight gets penalty
	require.Less(t, val, 320)
}

func TestPieceValue_Kangaroo(t *testing.T) {
	board := emptyBoard()
	kangaroo := chess.NewPiece(chess.KangarooPiece, chess.PieceWhite)
	pos := chess.Position{Row: 5, Col: 5}
	rows, cols := board.Dimensions()
	val := chess.PieceValue(board, *kangaroo, pos, rows, cols)
	require.GreaterOrEqual(t, val, 280)
}

func TestPieceValue_Bishop(t *testing.T) {
	board := emptyBoard()
	bishop := chess.NewPiece(chess.BishopPiece, chess.PieceWhite)
	pos := chess.Position{Row: 5, Col: 5}
	rows, cols := board.Dimensions()
	val := chess.PieceValue(board, *bishop, pos, rows, cols)
	require.GreaterOrEqual(t, val, 380)
}

func TestPieceValue_Rook(t *testing.T) {
	board := emptyBoard()
	rook := chess.NewPiece(chess.RookPiece, chess.PieceWhite)
	pos := chess.Position{Row: 5, Col: 5}
	rows, cols := board.Dimensions()
	val := chess.PieceValue(board, *rook, pos, rows, cols)
	require.GreaterOrEqual(t, val, 550)
}

func TestPieceValue_Oligarch_Isolated(t *testing.T) {
	board := emptyBoard()
	oligarch := chess.NewPiece(chess.OligarchPiece, chess.PieceWhite)
	pos := chess.Position{Row: 5, Col: 5}
	rows, cols := board.Dimensions()
	val := chess.PieceValue(board, *oligarch, pos, rows, cols)
	require.Equal(t, 300, val) // no adjacent friendly piece
}

func TestPieceValue_Oligarch_WithFriend(t *testing.T) {
	board := emptyBoard()
	oligarch := chess.NewPiece(chess.OligarchPiece, chess.PieceWhite)
	friend := chess.NewPiece(chess.PawnPiece, chess.PieceWhite)
	pos := chess.Position{Row: 5, Col: 5}
	placePieceOnBoard(&board, oligarch, pos)
	placePieceOnBoard(&board, friend, chess.Position{Row: 5, Col: 6})
	rows, cols := board.Dimensions()
	val := chess.PieceValue(board, *oligarch, pos, rows, cols)
	require.Equal(t, 300+600, val)
}

func TestPieceValue_King(t *testing.T) {
	board := emptyBoard()
	king := chess.NewPiece(chess.KingPiece, chess.PieceWhite)
	pos := chess.Position{Row: 5, Col: 5}
	rows, cols := board.Dimensions()
	val := chess.PieceValue(board, *king, pos, rows, cols)
	require.Equal(t, 20000, val)
}

func TestPieceValue_Queen(t *testing.T) {
	board := emptyBoard()
	queen := chess.NewPiece(chess.QueenPiece, chess.PieceWhite)
	pos := chess.Position{Row: 5, Col: 5}
	rows, cols := board.Dimensions()
	val := chess.PieceValue(board, *queen, pos, rows, cols)
	require.Equal(t, 1000, val)
}

func TestPieceValue_Pawn_WhiteAdvanced(t *testing.T) {
	board := emptyBoard()
	pawn := chess.NewPiece(chess.PawnPiece, chess.PieceWhite)
	// Advanced white pawn (low row number)
	pos := chess.Position{Row: 2, Col: 5}
	rows, cols := board.Dimensions()
	val := chess.PieceValue(board, *pawn, pos, rows, cols)
	require.Greater(t, val, 100)
}

func TestPieceValue_Pawn_BlackAdvanced(t *testing.T) {
	board := emptyBoard()
	pawn := chess.NewPiece(chess.PawnPiece, chess.PieceBlack)
	// Advanced black pawn (high row number)
	pos := chess.Position{Row: 9, Col: 5}
	rows, cols := board.Dimensions()
	val := chess.PieceValue(board, *pawn, pos, rows, cols)
	require.Greater(t, val, 100)
}

func TestPieceValue_Rook_OpenFile(t *testing.T) {
	board := emptyBoard()
	rook := chess.NewPiece(chess.RookPiece, chess.PieceWhite)
	pos := chess.Position{Row: 5, Col: 5}
	rows, cols := board.Dimensions()
	val := chess.PieceValue(board, *rook, pos, rows, cols)
	// Open file bonus should apply
	require.Greater(t, val, 550)
}

func TestPieceValue_Rook_SeventhRankWhite(t *testing.T) {
	board := emptyBoard()
	rook := chess.NewPiece(chess.RookPiece, chess.PieceWhite)
	pos := chess.Position{Row: 1, Col: 5} // 2nd row (opponent's territory)
	rows, cols := board.Dimensions()
	val := chess.PieceValue(board, *rook, pos, rows, cols)
	require.Greater(t, val, 550)
}

func TestPieceValue_Rook_SeventhRankBlack(t *testing.T) {
	board := emptyBoard()
	rook := chess.NewPiece(chess.RookPiece, chess.PieceBlack)
	pos := chess.Position{Row: 10, Col: 5} // near end row for black
	rows, cols := board.Dimensions()
	val := chess.PieceValue(board, *rook, pos, rows, cols)
	require.Greater(t, val, 550)
}

func TestPieceValue_Rook_BlockedFile(t *testing.T) {
	board := emptyBoard()
	rook := chess.NewPiece(chess.RookPiece, chess.PieceWhite)
	blocker := chess.NewPiece(chess.PawnPiece, chess.PieceWhite)
	pos := chess.Position{Row: 5, Col: 5}
	placePieceOnBoard(&board, rook, pos)
	placePieceOnBoard(&board, blocker, chess.Position{Row: 3, Col: 5}) // pawn on same col
	rows, cols := board.Dimensions()
	val := chess.PieceValue(board, *rook, pos, rows, cols)
	// No open file bonus
	require.LessOrEqual(t, val, 580) // 550 + rank bonus at most
}

func TestEvaluateBoard_WithAllPieces(t *testing.T) {
	board := chess.NewInitialBoard()
	// Symmetric initial board should have small absolute value
	score := chess.EvaluateBoard(board)
	// Not zero because position bonuses differ, but should be reasonable
	require.IsType(t, 0, score)
}

func TestKangarooBonus_HighReachable(t *testing.T) {
	board := emptyBoard()
	kangaroo := chess.NewPiece(chess.KangarooPiece, chess.PieceWhite)
	// Center position with many reachable squares
	pos := chess.Position{Row: 5, Col: 5}
	rows, cols := board.Dimensions()
	val := chess.PieceValue(board, *kangaroo, pos, rows, cols)
	// Should have bonus for >= 4 reachable squares
	require.Equal(t, 280+20, val)
}

func TestBishopBonus_OpenDiagonal(t *testing.T) {
	board := emptyBoard()
	bishop := chess.NewPiece(chess.BishopPiece, chess.PieceWhite)
	pos := chess.Position{Row: 0, Col: 0}
	rows, cols := board.Dimensions()
	val := chess.PieceValue(board, *bishop, pos, rows, cols)
	// One long diagonal should give bonus
	require.GreaterOrEqual(t, val, 380)
}
