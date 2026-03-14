import type { Board, GameState, Piece, PieceColor, Position } from "./types";
import { BOARD_SIZE } from "./constants";

/**
 * Checks if two positions refer to the same coordinates.
 */
export function arePositionsEqual(
    a: Position | null,
    b: Position | null
): boolean {
    if (a === null || b === null) return false;
    return a.row === b.row && a.col === b.col;
}

/**
 * Creates a deep copy of the board data.
 */
export function cloneBoard(board: Board): Board {
    return {
        data: board.data.map((row) =>
            row.map((square) => ({
                color: square.color,
                piece: square.piece ? { ...square.piece } : null,
            })),
        ),
    };
}

/**
 * Extracts row and col from DOM dataset.
 */
export function getPositionFromSquare(square: HTMLElement): Position | null {
    const { row, col } = square.dataset;

    if (row === undefined || col === undefined) {
        return null;
    }

    return {
        row: Number(row),
        col: Number(col),
    };
}

/**
 * Determines if a piece belongs to the player whose turn it currently is.
 */
export function isPieceSelectable(
    piece: Piece | null,
    gameState: GameState
): boolean {
    if (!piece) return false;
    return piece.color === gameState.turn;
}

/**
 * Safety check to ensure coordinates are within the board boundaries.
 */
export function isInsideBoard(row: number, col: number): boolean {
    return row >= 0 && row < BOARD_SIZE && col >= 0 && col < BOARD_SIZE;
}

/**
 * Checks if two sets of coordinates are the same.
 */
export function isSameSquare(
    fromRow: number,
    fromCol: number,
    toRow: number,
    toCol: number
): boolean {
    return fromRow === toRow && fromCol === toCol;
}

/**
 * Checks if the piece at the target location belongs to the specified color.
 */
export function isFriendlyPiece(
    row: number,
    col: number,
    color: PieceColor,
    board: Board
): boolean {
    const targetPiece = board.data[row]?.[col]?.piece;
    return !!targetPiece && targetPiece.color === color;
}
