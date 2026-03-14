import { cloneBoard, isInsideBoard, isSameSquare } from "./board";
import { getOppositeColor, isKingInCheck } from "./engine";
import type { Board, GameState, PieceColor, PieceKind, Position } from "./types";
import {
    PIECE_BLACK,
    PAWN_PIECE,
    ROOK_PIECE,
    OLIGARCH_PIECE,
    KNIGHT_PIECE,
    KANGAROO_PIECE,
    BISHOP_PIECE,
    KING_PIECE,
    QUEEN_PIECE,
    BOARD_SIZE
} from "./constants";

/**
 * Safely retrieves a piece at given coordinates.
 */
function getPiece(board: Board, row: number, col: number) {
    return board.data[row]?.[col]?.piece ?? null;
}

/**
 * Updates a square with a piece or null.
 */
function setPiece(
    board: Board,
    row: number,
    col: number,
    piece: Board["data"][number][number]["piece"],
): void {
    if (board.data[row]?.[col]) {
        board.data[row][col].piece = piece;
    }
}

function isFriendlyPieceAt(
    board: Board,
    row: number,
    col: number,
    color: PieceColor,
): boolean {
    const piece = getPiece(board, row, col);
    return piece !== null && piece.color === color;
}

/**
 * Entry point for executing a move. Returns a new State or null if illegal.
 */
export function tryApplyMove(
    gameState: GameState,
    fromPos: Position,
    toPos: Position,
): GameState | null {
    if (gameState.result !== null) return null;

    if (!canMovePiece(fromPos.row, fromPos.col, toPos.row, toPos.col, gameState)) {
        return null;
    }

    const movingColor = gameState.turn;
    const nextTurn = getOppositeColor(movingColor);
    const nextBoard = applyMoveToBoard(
        gameState.board,
        fromPos.row,
        fromPos.col,
        toPos.row,
        toPos.col,
    );

    return {
        ...gameState,
        board: nextBoard,
        turn: nextTurn,
        result: null,
        endReason: null,
    };
}

/**
 * Full legality check including King safety.
 */
export function canMovePiece(
    fromRow: number,
    fromCol: number,
    toRow: number,
    toCol: number,
    gameState: GameState,
): boolean {
    const piece = getPiece(gameState.board, fromRow, fromCol);

    if (!piece || piece.color !== gameState.turn) return false;

    if (!canMovePiecePseudoLegal(fromRow, fromCol, toRow, toCol, gameState)) {
        return false;
    }

    const nextBoard = applyMoveToBoard(
        gameState.board,
        fromRow,
        fromCol,
        toRow,
        toCol,
    );

    return !isKingInCheck(nextBoard, piece.color);
}

/**
 * Checks if the current player has any move that doesn't leave King in check.
 */
export function hasAnyLegalMove(gameState: GameState): boolean {
    for (let fromRow = 0; fromRow < BOARD_SIZE; fromRow++) {
        for (let fromCol = 0; fromCol < BOARD_SIZE; fromCol++) {
            const piece = getPiece(gameState.board, fromRow, fromCol);

            if (piece === null || piece.color !== gameState.turn) continue;

            for (let toRow = 0; toRow < BOARD_SIZE; toRow++) {
                for (let toCol = 0; toCol < BOARD_SIZE; toCol++) {
                    if (canMovePiece(fromRow, fromCol, toRow, toCol, gameState)) {
                        return true;
                    }
                }
            }
        }
    }
    return false;
}

function canMovePiecePseudoLegal(
    fromRow: number,
    fromCol: number,
    toRow: number,
    toCol: number,
    gameState: GameState,
): boolean {
    if (!isInsideBoard(fromRow, fromCol) || !isInsideBoard(toRow, toCol)) return false;
    if (isSameSquare(fromRow, fromCol, toRow, toCol)) return false;

    const piece = getPiece(gameState.board, fromRow, fromCol);
    if (!piece || isFriendlyPieceAt(gameState.board, toRow, toCol, piece.color)) return false;

    switch (piece.kind) {
        case PAWN_PIECE:
            return canMovePawn(fromRow, fromCol, toRow, toCol, gameState.board, piece.color);
        case ROOK_PIECE:
            return canMoveRook(fromRow, fromCol, toRow, toCol, gameState.board);
        case OLIGARCH_PIECE:
            return canMoveOligarch(fromRow, fromCol, toRow, toCol, gameState.board);
        case KNIGHT_PIECE:
            return canMoveKnight(fromRow, fromCol, toRow, toCol);
        case KANGAROO_PIECE:
            return canMoveKangaroo(fromRow, fromCol, toRow, toCol);
        case BISHOP_PIECE:
            return canMoveBishop(fromRow, fromCol, toRow, toCol, gameState.board);
        case QUEEN_PIECE:
            return canMoveQueen(fromRow, fromCol, toRow, toCol, gameState.board);
        case KING_PIECE:
            return canMoveKing(fromRow, fromCol, toRow, toCol, gameState.board);
        default:
            return false;
    }
}

export function canMoveKing(fromRow: number, fromCol: number, toRow: number, toCol: number, board: Board): boolean {
    if (!getPiece(board, fromRow, fromCol)) return false;
    return Math.abs(toRow - fromRow) <= 1 && Math.abs(toCol - fromCol) <= 1;
}

export function canMoveQueen(fromRow: number, fromCol: number, toRow: number, toCol: number, board: Board): boolean {
    return canMoveRook(fromRow, fromCol, toRow, toCol, board) || canMoveBishop(fromRow, fromCol, toRow, toCol, board);
}

export function canMoveBishop(fromRow: number, fromCol: number, toRow: number, toCol: number, board: Board): boolean {
    if (!getPiece(board, fromRow, fromCol)) return false;

    const rowDiff = Math.abs(toRow - fromRow);
    const colDiff = Math.abs(toCol - fromCol);
    if (rowDiff !== colDiff) return false;

    const rowStep = toRow > fromRow ? 1 : -1;
    const colStep = toCol > fromCol ? 1 : -1;

    let row = fromRow + rowStep;
    let col = fromCol + colStep;

    while (row !== toRow && col !== toCol) {
        if (getPiece(board, row, col) !== null) return false;
        row += rowStep;
        col += colStep;
    }
    return true;
}

export function canMoveKangaroo(fromRow: number, fromCol: number, toRow: number, toCol: number): boolean {
    const rowDiff = Math.abs(toRow - fromRow);
    const colDiff = Math.abs(toCol - fromCol);
    return (rowDiff === 2 && colDiff === 0) || (rowDiff === 0 && colDiff === 2) ||
           (rowDiff === 3 && colDiff === 0) || (rowDiff === 0 && colDiff === 3);
}

export function canMoveKnight(fromRow: number, fromCol: number, toRow: number, toCol: number): boolean {
    const rowDiff = Math.abs(toRow - fromRow);
    const colDiff = Math.abs(toCol - fromCol);
    return (rowDiff === 2 && colDiff === 1) || (rowDiff === 1 && colDiff === 2);
}

export function canMoveOligarch(fromRow: number, fromCol: number, toRow: number, toCol: number, board: Board): boolean {
    const piece = getPiece(board, fromRow, fromCol);
    if (!piece) return false;

    return hasAdjacentFriendlyPiece(fromRow, fromCol, piece.color, board)
        ? canMoveQueen(fromRow, fromCol, toRow, toCol, board)
        : canMoveKing(fromRow, fromCol, toRow, toCol, board);
}

export function canMoveRook(fromRow: number, fromCol: number, toRow: number, toCol: number, board: Board): boolean {
    if (fromRow !== toRow && fromCol !== toCol) return false;

    const rowStep = toRow === fromRow ? 0 : (toRow > fromRow ? 1 : -1);
    const colStep = toCol === fromCol ? 0 : (toCol > fromCol ? 1 : -1);

    let row = fromRow + rowStep;
    let col = fromCol + colStep;

    while (row !== toRow || col !== toCol) {
        if (getPiece(board, row, col) !== null) return false;
        row += rowStep;
        col += colStep;
    }
    return true;
}

export function canMovePawn(fromRow: number, fromCol: number, toRow: number, toCol: number, board: Board, color: PieceColor): boolean {
    const direction = color === PIECE_BLACK ? 1 : -1;
    const rowDiff = toRow - fromRow;
    const colDiff = toCol - fromCol;
    const targetPiece = getPiece(board, toRow, toCol);

    if (fromCol === toCol) {
        if (rowDiff === direction && !targetPiece) return true;
        if (isPieceAtStartingPosition(fromRow, fromCol, PAWN_PIECE, color) && 
            rowDiff === 2 * direction && !targetPiece && !getPiece(board, fromRow + direction, fromCol)) {
            return true;
        }
        return false;
    }

    return Math.abs(colDiff) === 1 && rowDiff === direction && !!targetPiece && targetPiece.color !== color;
}

function applyMoveToBoard(board: Board, fromRow: number, fromCol: number, toRow: number, toCol: number): Board {
    const nextBoard = cloneBoard(board);
    const piece = getPiece(nextBoard, fromRow, fromCol);
    setPiece(nextBoard, toRow, toCol, piece);
    setPiece(nextBoard, fromRow, fromCol, null);
    return nextBoard;
}

function hasAdjacentFriendlyPiece(row: number, col: number, color: PieceColor, board: Board): boolean {
    for (let rd = -1; rd <= 1; rd++) {
        for (let cd = -1; cd <= 1; cd++) {
            if (rd === 0 && cd === 0) continue;
            const nr = row + rd;
            const nc = col + cd;
            if (isInsideBoard(nr, nc) && isFriendlyPieceAt(board, nr, nc, color)) return true;
        }
    }
    return false;
}

function isPieceAtStartingPosition(row: number, col: number, kind: PieceKind, color: PieceColor): boolean {
    if (kind === PAWN_PIECE) return color === PIECE_BLACK ? row === 1 : row === 10;
    
    const backRow = color === PIECE_BLACK ? 0 : 11;
    if (row !== backRow) return false;

    const map: Record<number, number[]> = {
        [ROOK_PIECE]: [0, 11], [OLIGARCH_PIECE]: [1, 10], [KNIGHT_PIECE]: [2, 9],
        [KANGAROO_PIECE]: [3, 8], [BISHOP_PIECE]: [4, 7], [KING_PIECE]: [6], [QUEEN_PIECE]: [5]
    };
    return map[kind]?.includes(col) ?? false;
}
