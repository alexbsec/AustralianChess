export type PieceColor = 0 | 1;
export type PieceKind = 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7;


export type Position = {
    row: number;
    col: number;
};

export type Piece = {
    kind: PieceKind;
    color: PieceColor;
};

export type Square = {
    color: number;
    piece: Piece | null;
};

export type Board = {
    data: Square[][];
};

export type GameState = {
    board: Board;
    turn: PieceColor;
    result: unknown | null;
    endReason?: string | null;
    end_reason?: string | null;
};

export type PieceSide = "white" | "black";
