import { IsInsideBoard, IsSameSquare } from "./board";
import { CanMoveBishop, CanMoveKangaroo, CanMoveKing, CanMoveKnight, CanMoveOligarch, CanMoveQueen, CanMoveRook, HasAnyLegalMove } from "./movements";
import type { Board, GameState, PieceColor } from "./types";

export function EndTurn(state: GameState): GameState {
    const nextState: GameState = {
        ...state,
    };

    if (IsCheckmate(state)) {
        const winner = GetOppositeColor(state.turn);
        nextState.result = `${winner}-win`;
        nextState.endReason = "checkmate";
        return nextState;
    }

    if (IsStalemate(state)) {
        nextState.result = "draw";
        nextState.endReason = "stalemate";
        return nextState
    }

    return state;
}

export function IsKingInCheck(board: Board, color: PieceColor): boolean {
    const kingPosition = FindKing(board, color);
    if (kingPosition === null) {
        return true;
    }

    const enemyColor: PieceColor = color === "white" ? "black" : "white";

    return IsSquareAttacked(
        kingPosition.row,
        kingPosition.col,
        enemyColor,
        board
    );
}

export function GetOppositeColor(color: PieceColor): PieceColor {
    return color === "white" ? "black" : "white";
}

function IsCheckmate(gameState: GameState): boolean {
    return IsKingInCheck(gameState.board, gameState.turn) && !HasAnyLegalMove(gameState);
}

function IsStalemate(gameState: GameState): boolean {
    return !IsKingInCheck(gameState.board, gameState.turn) && !HasAnyLegalMove(gameState);
}

function FindKing(board: Board, color: PieceColor): { row: number; col: number } | null {
    for (let row = 0; row < 12; row++) {
        for (let col = 0; col < 12; col++) {
            const piece = board[row][col];
            if (piece !== null && piece.kind === "king" && piece.color === color) {
                return { row, col };
            }
        }
    }

    return null;
}

function IsSquareAttacked(
    row: number,
    col: number,
    byColor: PieceColor,
    board: Board,
): boolean {
    for (let fromRow = 0; fromRow < 12; fromRow++) {
        for (let fromCol = 0; fromCol < 12; fromCol++) {
            const piece = board[fromRow][fromCol];

            if (piece === null || piece.color !== byColor) {
                continue;
            }

            if (CanPieceAttackSquare(fromRow, fromCol, row, col, board)) {
                return true;
            }
        }
    }

    return false;
}

function CanPieceAttackSquare(
    fromRow: number,
    fromCol: number,
    toRow: number,
    toCol: number,
    board: Board,
): boolean {
    const piece = board[fromRow][fromCol];
    if (piece === null) {
        return false;
    }

    if (!IsInsideBoard(toRow, toCol)) {
        return false;
    }

    if (IsSameSquare(fromRow, fromCol, toRow, toCol)) {
        return false;
    }

    switch (piece.kind) {
        case "pawn": {
            const direction = piece.color === "black" ? 1 : -1;
            const rowDiff = toRow - fromRow;
            const colDiff = Math.abs(toCol - fromCol);

            return rowDiff === direction && colDiff === 1;
        }

        case "rook":
            return CanMoveRook(fromRow, fromCol, toRow, toCol, board);

        case "bishop":
            return CanMoveBishop(fromRow, fromCol, toRow, toCol, board);

        case "queen":
            return CanMoveQueen(fromRow, fromCol, toRow, toCol, board);

        case "knight":
            return CanMoveKnight(fromRow, fromCol, toRow, toCol);

        case "kangaroo":
            return CanMoveKangaroo(fromRow, fromCol, toRow, toCol);

        case "king":
            return CanMoveKing(fromRow, fromCol, toRow, toCol, board);

        case "oligarch":
            return CanMoveOligarch(fromRow, fromCol, toRow, toCol, board);

        default:
            return false;
    }
}