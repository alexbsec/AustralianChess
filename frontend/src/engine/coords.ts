import { BOARD_SIZE } from "./constants";
import type { PieceColor, Position, PieceKind, PieceSide, Square, GameState } from "./types";

export function positionsEqual(a: Position | null, b: Position | null): boolean {
    if (!a || !b) {
        return false;
    }

    return a.row === b.row && a.col === b.col;
}

export function pieceColorToSide(color: PieceColor): PieceSide {
    return color === 0 ? "white" : "black";
}

export function pieceKindToName(kind: PieceKind): string {
    switch (kind) {
        case 0:
            return "pawn";
        case 1:
            return "rook";
        case 2:
            return "oligarch";
        case 3:
            return "knight";
        case 4:
            return "kangaroo";
        case 5:
            return "bishop";
        case 6:
            return "king";
        case 7:
            return "queen";
        default:
            throw new Error(`Unknown piece kind: ${kind}`);
    }
}

export function displayToBoardPosition(
    displayPos: Position,
    playerColor: PieceColor,
): Position {
    if (playerColor === 0) {
        return displayPos;
    }

    return {
        row: BOARD_SIZE - 1 - displayPos.row,
        col: BOARD_SIZE - 1 - displayPos.col,
    };
}

export function getSquareAtPosition(
    gameState: GameState,
    position: Position,
): Square | null {
    if (
        position.row < 0 ||
        position.row >= BOARD_SIZE ||
        position.col < 0 ||
        position.col >= BOARD_SIZE
    ) {
        return null;
    }

    return gameState.board.data[position.row][position.col] ?? null;
}
