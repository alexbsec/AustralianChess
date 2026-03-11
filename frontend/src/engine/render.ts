import { CanMovePiece } from "./movements";
import type { Piece, Board, Position, GameState, UIState } from "./types";

function GetPieceImagePath(piece: Piece): string {
    return `/pieces/${piece.color}/${piece.color}_${piece.kind}.png`;
}

function CreateSquareElement(
    row: number,
    col: number,
    piece: Piece | null,
    uiState: UIState,
    possibleMoves: Position[] = [],
): HTMLDivElement {
    const square = document.createElement("div");
    square.classList.add("square");

    const isLightSquare = (row + col) % 2 === 0;
    square.classList.add(isLightSquare ? "light" : "dark");

    square.dataset.row = String(row);
    square.dataset.col = String(col);

    if (
        uiState.selected &&
        uiState.selected.row === row &&
        uiState.selected.col === col
    ) {
        square.classList.add("selected");
    }

    const isPossibleMove = possibleMoves.some(
        (move) => move.row === row && move.col === col
    );

    if (isPossibleMove) {
        square.classList.add("possible-move");

        const moveDot = document.createElement("div");
        moveDot.classList.add("move-dot");
        square.appendChild(moveDot);
    }

    if (piece !== null) {
        const img = document.createElement("img");
        img.classList.add("piece");
        img.src = GetPieceImagePath(piece);
        img.alt = `${piece.color} ${piece.kind}`;
        square.appendChild(img);
    }

    return square;
}

export function RenderBoard(
    board: Board,
    container: HTMLElement,
    uiState: UIState,
    possibleMoves: Position[] = [],
): void {
    container.innerHTML = "";

    for (let row = 0; row < 12; row++) {
        for (let col = 0; col < 12; col++) {
            const piece = board[row][col];
            const square = CreateSquareElement(row, col, piece, uiState, possibleMoves);
            container.appendChild(square);
        }
    }
}

export function GetPseudoLegalMoves(
    gameState: GameState,
    fromRow: number,
    fromCol: number,
): Position[] {
    let moves: Position[] = [];

    for (let row = 0; row < 12; row++) {
        for (let col = 0; col < 12; col++) {
            if (CanMovePiece(fromRow, fromCol, row, col, gameState)) {
                moves.push({ row, col });
            }
        }
    }

    return moves;
}
