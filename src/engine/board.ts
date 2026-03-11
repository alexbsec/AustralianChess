import type { Board, GameState, Piece, PieceColor, Position } from "./types";

export function ArePositionsEqual(
    a: Position | null,
    b: Position | null
): boolean {
    if (a === null || b === null) {
        return false;
    }

    return a.row === b.row && a.col === b.col;
}

export function CloneBoard(board: Board): Board {
    return board.map((row) => row.map((square) => square));
}

export function GetPositionFromSquare(square: HTMLElement): Position | null {
    const row = square.dataset.row;
    const col = square.dataset.col;

    if (row === undefined || col === undefined) {
        return null;
    }

    return {
        row: Number(row),
        col: Number(col),
    };
}

export function IsPieceSelectable(
    piece: Piece | null,
    gameState: GameState
): boolean {
    if (piece === null) {
        return false;
    }

    return piece.color === gameState.turn;
}

export function IsInsideBoard(row: number, col: number): boolean {
    return row >= 0 && row < 12 && col >= 0 && col < 12;
}

export function IsSameSquare(fromRow: number, fromCol: number, toRow: number, toCol: number): boolean {
    return fromRow == toRow && fromCol == toCol;
}

export function IsFriendlyPiece(row: number, col: number, color: PieceColor, board: Board): boolean {
    const target = board[row][col];
    return target !== null && target.color === color;
}

