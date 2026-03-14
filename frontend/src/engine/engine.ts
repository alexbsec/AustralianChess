import { isInsideBoard, isSameSquare } from "./board";
import {
    canMoveBishop,
    canMoveKangaroo,
    canMoveKing,
    canMoveKnight,
    canMoveOligarch,
    canMoveQueen,
    canMoveRook,
    hasAnyLegalMove,
} from "./movements";
import type { Board, GameState, PieceColor, Position } from "./types";
import {
    PIECE_WHITE,
    PIECE_BLACK,
    PAWN_PIECE,
    ROOK_PIECE,
    OLIGARCH_PIECE,
    KNIGHT_PIECE,
    KANGAROO_PIECE,
    BISHOP_PIECE,
    KING_PIECE,
    QUEEN_PIECE,
    BOARD_SIZE,
} from "./constants";

/**
 * Utility to safely grab a piece from the board
 */
function getPiece(board: Board, row: number, col: number) {
    return board.data[row]?.[col]?.piece ?? null;
}

/**
 * Checks for game-over conditions and updates the state results
 */
export function endTurn(state: GameState): GameState {
    const nextState: GameState = { ...state };

    if (isCheckmate(state)) {
        const winner = getOppositeColor(state.turn);
        nextState.result = `${winner}-win`;
        nextState.endReason = "checkmate";
        return nextState;
    }

    if (isStalemate(state)) {
        nextState.result = "draw";
        nextState.endReason = "stalemate";
        return nextState;
    }

    return state;
}

export function isKingInCheck(board: Board, color: PieceColor): boolean {
    const kingPosition = findKing(board, color);
    
    // If king is missing (shouldn't happen in legal play), treat as check
    if (kingPosition === null) return true;

    const enemyColor = getOppositeColor(color);

    return isSquareAttacked(
        kingPosition.row,
        kingPosition.col,
        enemyColor,
        board,
    );
}

export function getOppositeColor(color: PieceColor): PieceColor {
    return color === PIECE_WHITE ? PIECE_BLACK : PIECE_WHITE;
}

export function isCheckmate(gameState: GameState): boolean {
    return isKingInCheck(gameState.board, gameState.turn) && !hasAnyLegalMove(gameState);
}

export function isStalemate(gameState: GameState): boolean {
    return !isKingInCheck(gameState.board, gameState.turn) && !hasAnyLegalMove(gameState);
}

function findKing(board: Board, color: PieceColor): Position | null {
    for (let row = 0; row < BOARD_SIZE; row++) {
        for (let col = 0; col < BOARD_SIZE; col++) {
            const piece = getPiece(board, row, col);
            if (piece?.kind === KING_PIECE && piece.color === color) {
                return { row, col };
            }
        }
    }
    return null;
}

export function isSquareAttacked(
    row: number,
    col: number,
    byColor: PieceColor,
    board: Board,
): boolean {
    for (let fromRow = 0; fromRow < BOARD_SIZE; fromRow++) {
        for (let fromCol = 0; fromCol < BOARD_SIZE; fromCol++) {
            const piece = getPiece(board, fromRow, fromCol);

            if (piece === null || piece.color !== byColor) {
                continue;
            }

            if (canPieceAttackSquare(fromRow, fromCol, row, col, board)) {
                return true;
            }
        }
    }
    return false;
}

function canPieceAttackSquare(
    fromRow: number,
    fromCol: number,
    toRow: number,
    toCol: number,
    board: Board,
): boolean {
    const piece = getPiece(board, fromRow, fromCol);
    if (!piece || !isInsideBoard(toRow, toCol) || isSameSquare(fromRow, fromCol, toRow, toCol)) {
        return false;
    }

    switch (piece.kind) {
        case PAWN_PIECE: {
            const direction = piece.color === PIECE_BLACK ? 1 : -1;
            const rowDiff = toRow - fromRow;
            const colDiff = Math.abs(toCol - fromCol);
            return rowDiff === direction && colDiff === 1;
        }
        case ROOK_PIECE:
            return canMoveRook(fromRow, fromCol, toRow, toCol, board);
        case BISHOP_PIECE:
            return canMoveBishop(fromRow, fromCol, toRow, toCol, board);
        case QUEEN_PIECE:
            return canMoveQueen(fromRow, fromCol, toRow, toCol, board);
        case KNIGHT_PIECE:
            return canMoveKnight(fromRow, fromCol, toRow, toCol);
        case KANGAROO_PIECE:
            return canMoveKangaroo(fromRow, fromCol, toRow, toCol);
        case KING_PIECE:
            return canMoveKing(fromRow, fromCol, toRow, toCol, board);
        case OLIGARCH_PIECE:
            return canMoveOligarch(fromRow, fromCol, toRow, toCol, board);
        default:
            return false;
    }
}
