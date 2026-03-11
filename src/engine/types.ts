export type PieceColor = "white" | "black";

export type PieceKind = "king" | "queen" | "rook" | "bishop" | "kangaroo" | "knight" | "oligarch" | "pawn";

type GameResult = "white-win" | "black-win" | "draw" | null;

export interface Piece {
    kind: PieceKind;
    color: PieceColor;
};

export type Square = Piece | null;
export type Board = Square[][];

export interface Position {
    row: number;
    col: number;
};

export interface PixelPosition {
    x: number;
    y: number;
}

export interface GameState {
    board: Board;
    turn: PieceColor;
    result: GameResult;
    endReason: "checkmate" | "stalemate" | null;
}

export interface UIState {
    selected: Position | null;
    hovered: Position | null;
    draggedFrom: Position | null;
    isDragging: boolean;
};