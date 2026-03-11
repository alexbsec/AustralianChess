import { CloneBoard, IsFriendlyPiece, IsInsideBoard, IsSameSquare } from "./board";
import { GetOppositeColor, IsKingInCheck } from "./engine";
import type { Board, GameState, PieceColor, PieceKind, Position } from "./types";

export function TryApplyMove(
    gameState: GameState,
    fromPos: Position,
    toPos: Position,
): GameState | null {
    if (gameState.result !== null) {
        return null;
    }

    if (!CanMovePiece(fromPos.row, fromPos.col, toPos.row, toPos.col, gameState)) {
        return null;
    }

    const movingColor = gameState.turn;
    const nextTurn = GetOppositeColor(movingColor);
    const nextBoard = ApplyMoveToBoard(
        gameState.board,
        fromPos.row,
        fromPos.col,
        toPos.row,
        toPos.col
    );

    return {
        ...gameState,
        board: nextBoard,
        turn: nextTurn,
        result: null,
        endReason: null,
    };
}

export function CanMovePiece(
    fromRow: number,
    fromCol: number,
    toRow: number,
    toCol: number,
    gameState: GameState,
): boolean {
    const piece = gameState.board[fromRow][fromCol];
    if (piece === null) {
        return false;
    }

    if (piece.color !== gameState.turn) {
        return false;
    }

    if (!CanMovePiecePseudoLegal(fromRow, fromCol, toRow, toCol, gameState)) {
        return false;
    }

    const nextBoard = ApplyMoveToBoard(
        gameState.board,
        fromRow,
        fromCol,
        toRow,
        toCol
    );

    return !IsKingInCheck(nextBoard, piece.color);
}

export function HasAnyLegalMove(gameState: GameState): boolean {
    for (let fromRow = 0; fromRow < 12; fromRow++) {
        for (let fromCol = 0; fromCol < 12; fromCol++) {
            const piece = gameState.board[fromRow][fromCol];

            if (piece === null || piece.color !== gameState.turn) {
                continue;
            }

            for (let toRow = 0; toRow < 12; toRow++) {
                for (let toCol = 0; toCol < 12; toCol++) {
                    if (CanMovePiece(fromRow, fromCol, toRow, toCol, gameState)) {
                        console.log(`has legal move. Can move ${piece.kind}`)
                        return true;
                    }
                }
            }
        }
    }

    return false;
}

function CanMovePiecePseudoLegal(
    fromRow: number,
    fromCol: number,
    toRow: number,
    toCol: number,
    gameState: GameState,
): boolean {
    if (!IsInsideBoard(fromRow, fromCol) || !IsInsideBoard(toRow, toCol)) {
        return false;
    }

    if (IsSameSquare(fromRow, fromCol, toRow, toCol)) {
        return false;
    }

    const piece = gameState.board[fromRow][fromCol];
    if (piece == null) {
        return false;
    }

    if (IsFriendlyPiece(toRow, toCol, piece.color, gameState.board)) {
        return false;
    }
    
    switch (piece.kind) {
        case "pawn":
            return CanMovePawn(fromRow, fromCol, toRow, toCol, gameState.board, piece.color);
        case "rook":
            return CanMoveRook(fromRow, fromCol, toRow, toCol, gameState.board);
        case "oligarch":
            return CanMoveOligarch(fromRow, fromCol, toRow, toCol, gameState.board);
        case "knight":
            return CanMoveKnight(fromRow, fromCol, toRow, toCol);
        case "kangaroo":
            return CanMoveKangaroo(fromRow, fromCol, toRow, toCol);
        case "bishop":
            return CanMoveBishop(fromRow, fromCol, toRow, toCol, gameState.board);
        case "queen":
            return CanMoveQueen(fromRow, fromCol, toRow, toCol, gameState.board);
        case "king":
            return CanMoveKing(fromRow, fromCol, toRow, toCol, gameState.board);
    }

    return false;
}

export function CanMoveKing(
    fromRow: number,
    fromCol: number,
    toRow: number,
    toCol: number,
    board: Board
): boolean {
    const piece = board[fromRow][fromCol];
    if (piece === null) {
        return false;
    }

    const rowDiff = Math.abs(toRow - fromRow);
    const colDiff = Math.abs(toCol - fromCol);

    return rowDiff <= 1 && colDiff <= 1;
}

export function CanMoveQueen(
    fromRow: number,
    fromCol: number,
    toRow: number,
    toCol: number,
    board: Board
): boolean {
    return (
        CanMoveRook(fromRow, fromCol, toRow, toCol, board) ||
        CanMoveBishop(fromRow, fromCol, toRow, toCol, board)
    );
}

export function CanMoveBishop(
    fromRow: number,
    fromCol: number,
    toRow: number,
    toCol: number,
    board: Board
): boolean {
    const piece = board[fromRow][fromCol];
    if (piece === null) {
        return false;
    }

    const rowDiff = Math.abs(toRow - fromRow);
    const colDiff = Math.abs(toCol - fromCol);

    if (rowDiff !== colDiff) {
        return false;
    }

    const rowStep = toRow > fromRow ? 1 : -1;
    const colStep = toCol > fromCol ? 1 : -1;

    let row = fromRow + rowStep;
    let col = fromCol + colStep;

    while (row !== toRow && col !== toCol) {
        if (board[row][col] !== null) {
            return false;
        }

        row += rowStep;
        col += colStep;
    }

    return true;
}

export function CanMoveKangaroo(
    fromRow: number,
    fromCol: number,
    toRow: number,
    toCol: number,
): boolean {
    const rowDiff = Math.abs(toRow - fromRow);
    const colDiff = Math.abs(toCol - fromCol);

    return (
        (rowDiff === 2 && colDiff === 0) ||
        (rowDiff === 0 && colDiff === 2) ||
        (rowDiff === 3 && colDiff === 0) ||
        (rowDiff === 0 && colDiff === 3)
    );
}

export function CanMoveKnight(
    fromRow: number,
    fromCol: number,
    toRow: number,
    toCol: number,
): boolean {
    const rowDiff = Math.abs(toRow - fromRow);
    const colDiff = Math.abs(toCol - fromCol);

    return (
        (rowDiff === 2 && colDiff === 1) ||
        (rowDiff === 1 && colDiff === 2)
    );
}

export function CanMoveOligarch(fromRow: number, fromCol: number, toRow: number, toCol: number, board: Board): boolean {
    const oligarchPiece = board[fromRow][fromCol];

    if (oligarchPiece === null) {
        return false;
    }

    const hasFriendlyNeighbor = HasAdjacentFriendlyPiece(fromRow, fromCol, oligarchPiece.color, board);

    if (hasFriendlyNeighbor) {
        return CanMoveQueen(fromRow, fromCol, toRow, toCol, board);
    }

    return CanMoveKing(fromRow, fromCol, toRow, toCol, board);
}

export function CanMoveRook(fromRow: number, fromCol: number, toRow: number, toCol: number, board: Board): boolean {
    const isSameRow = fromRow === toRow;
    const isSameCol = fromCol === toCol;

    if (!isSameRow && !isSameCol) {
        return false;
    }

    const rowStep = toRow === fromRow ? 0 : (toRow > fromRow ? 1 : -1);
    const colStep = toCol === fromCol ? 0 : (toCol > fromCol ? 1 : -1);

    let row = fromRow + rowStep;
    let col = fromCol + colStep;

    while (row !== toRow || col !== toCol) {
        if (board[row][col] !== null) {
            // a piece is there
            return false;
        }

        row += rowStep;
        col += colStep;
    }

    return true;
}

export function CanMovePawn(fromRow: number, fromCol: number, toRow: number, toCol: number, board: Board, color: PieceColor): boolean {
    const direction = color === "black" ? 1 : -1;
    const rowDiff = toRow - fromRow;
    const colDiff = toCol - fromCol;
    const targetPiece = board[toRow][toCol];

    // Forward movement: same col only, destination must be empty
    if (fromCol == toCol) {
        if (rowDiff === direction && targetPiece === null) {
            return true;
        }

        if (
            IsPieceAtStartingPosition(fromRow, fromCol, "pawn", color) && 
            rowDiff === 2 * direction && 
            targetPiece === null && 
            board[fromRow + direction][fromCol] === null
        ) {
            return true;
        }

        return false;
    }

    // Diagonal capture: exactly one col sideways, one row forward, target must be enemy
    if (colDiff === 1 && rowDiff === direction && targetPiece !== null && targetPiece.color !== color) {
        return true;
    }

    return false;
}



function ApplyMoveToBoard(
    board: Board,
    fromRow: number,
    fromCol: number,
    toRow: number,
    toCol: number,
): Board {
    const nextBoard = CloneBoard(board);
    const piece = nextBoard[fromRow][fromCol];

    nextBoard[toRow][toCol] = piece;
    nextBoard[fromRow][fromCol] = null;

    return nextBoard;
}



function HasAdjacentFriendlyPiece(
    row: number,
    col: number,
    color: PieceColor,
    board: Board
): boolean {
    for (let rowOffset = -1; rowOffset <= 1; rowOffset++) {
        for (let colOffset = -1; colOffset <= 1; colOffset++) {
            if (rowOffset === 0 && colOffset === 0) {
                continue;
            }

            const nextRow = row + rowOffset;
            const nextCol = col + colOffset;

            if (!IsInsideBoard(nextRow, nextCol)) {
                continue;
            }

            const neighbor = board[nextRow][nextCol];
            if (neighbor !== null && neighbor.color === color) {
                return true;
            }
        }
    }

    return false;
}

function IsPieceAtStartingPosition(
    row: number,
    col: number,
    pieceKind: PieceKind,
    pieceColor: PieceColor,
): boolean {
    if (pieceKind === "pawn") {
        return pieceColor === "black" ? row === 1 : row === 10;
    }

    const backRow = pieceColor === "black" ? 0 : 11;

    if (row !== backRow) {
        return false;
    }

    const startingColumnsByPiece: Record<Exclude<PieceKind, "pawn">, number[]> = {
        rook: [0, 11],
        oligarch: [1, 10],
        knight: [2, 9],
        kangaroo: [3, 8],
        bishop: [4, 7],
        queen: [5],
        king: [6],
    };

    return startingColumnsByPiece[pieceKind].includes(col);
}

