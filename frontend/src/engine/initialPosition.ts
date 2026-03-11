import type { Board, Piece, PieceColor, PieceKind } from "./types"

function CreatePiece(kind: PieceKind, color: PieceColor): Piece {
    return { kind, color };
}

function CreateEmptyBoard(): Board {
    return Array.from({ length: 12 }, () => 
        Array.from({ length: 12 }, () => null)
    );
}

export function CreateInitialBoard(): Board {
    const board = CreateEmptyBoard();

    // Black pieces
    board[0][0] = CreatePiece("rook", "black");
    board[0][1] = CreatePiece("oligarch", "black");
    board[0][2] = CreatePiece("knight", "black");
    board[0][3] = CreatePiece("kangaroo", "black");
    board[0][4] = CreatePiece("bishop", "black");
    board[0][5] = CreatePiece("queen", "black");
    board[0][6] = CreatePiece("king", "black");
    board[0][7] = CreatePiece("bishop", "black");
    board[0][8] = CreatePiece("kangaroo", "black");
    board[0][9] = CreatePiece("knight", "black");
    board[0][10] = CreatePiece("oligarch", "black");
    board[0][11] = CreatePiece("rook", "black");

    for (let i = 0; i < 12; i++) {
        board[1][i] = CreatePiece("pawn", "black");
    }

    // White pieces
    board[11][0] = CreatePiece("rook", "white");
    board[11][1] = CreatePiece("oligarch", "white");
    board[11][2] = CreatePiece("knight", "white");
    board[11][3] = CreatePiece("kangaroo", "white");
    board[11][4] = CreatePiece("bishop", "white");
    board[11][5] = CreatePiece("queen", "white");
    board[11][6] = CreatePiece("king", "white");
    board[11][7] = CreatePiece("bishop", "white");
    board[11][8] = CreatePiece("kangaroo", "white");
    board[11][9] = CreatePiece("knight", "white");
    board[11][10] = CreatePiece("oligarch", "white");
    board[11][11] = CreatePiece("rook", "white");

    for (let i = 0; i < 12; i++) {
        board[10][i] = CreatePiece("pawn", "white");
    }

    return board;
}