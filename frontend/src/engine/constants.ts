import type { PieceColor, PieceKind, PieceSide } from "./types";

export const BOARD_SIZE = 12;

export const PIECE_WHITE: PieceColor = 0;
export const PIECE_BLACK: PieceColor = 1;

export const PAWN_PIECE: PieceKind = 0;
export const ROOK_PIECE: PieceKind = 1;
export const OLIGARCH_PIECE: PieceKind = 2;
export const KNIGHT_PIECE: PieceKind = 3;
export const KANGAROO_PIECE: PieceKind = 4;
export const BISHOP_PIECE: PieceKind = 5;
export const KING_PIECE: PieceKind = 6;
export const QUEEN_PIECE: PieceKind = 7;


export const PIECE_ASSETS: Record<PieceSide, Record<string, string>> = {
    white: {
        pawn: "/pieces/white/white_pawn.png",
        rook: "/pieces/white/white_rook.png",
        oligarch: "/pieces/white/white_oligarch.png",
        knight: "/pieces/white/white_knight.png",
        kangaroo: "/pieces/white/white_kangaroo.png",
        bishop: "/pieces/white/white_bishop.png",
        king: "/pieces/white/white_king.png",
        queen: "/pieces/white/white_queen.png",
    },
    black: {
        pawn: "/pieces/black/black_pawn.png",
        rook: "/pieces/black/black_rook.png",
        oligarch: "/pieces/black/black_oligarch.png",
        knight: "/pieces/black/black_knight.png",
        kangaroo: "/pieces/black/black_kangaroo.png",
        bishop: "/pieces/black/black_bishop.png",
        king: "/pieces/black/black_king.png",
        queen: "/pieces/black/black_queen.png",
    },
};


