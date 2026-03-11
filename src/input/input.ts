import { ArePositionsEqual, IsPieceSelectable } from "../engine/board";
import { CanMovePiece } from "../engine/movements";
import type { GameState, InputInteractionResult, Position } from "../engine/types";

type BoardInputHandlers = {
    onMouseDown: (position: Position | null) => void;
    onMouseMove: (position: Position | null) => void;
    onMouseUp: (position: Position | null) => void;
    onMouseLeave: () => void;
};

export function AttachBoardInput(
    container: HTMLElement,
    handlers: BoardInputHandlers
): void {
    container.addEventListener("mousedown", (event) => {
        handlers.onMouseDown(GetPositionFromEvent(event));
    });

    container.addEventListener("mousemove", (event) => {
        handlers.onMouseMove(GetPositionFromEvent(event));
    });

    container.addEventListener("mouseup", (event) => {
        handlers.onMouseUp(GetPositionFromEvent(event));
    });

    container.addEventListener("mouseleave", () => {
        handlers.onMouseLeave();
    });
}

function GetPositionFromEvent(event: MouseEvent): Position | null {
    const target = event.target as HTMLElement | null;
    if (!target) {
        return null;
    }

    const square = target.closest(".square") as HTMLElement | null;
    if (!square) {
        return null;
    }

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

export function HandleMouseMove2(
    row: number,
    col: number,
    gameState: GameState,
): boolean {
    const nextHovered = { row, col };

    if (ArePositionsEqual(gameState.hovered, nextHovered)) {
        return false;
    }

    gameState.hovered = nextHovered;
    return true;
}

function HandleMouseDown2(
    row: number,
    col: number,
    gameState: GameState,
): boolean {
    if (gameState.result !== null) {
        return false;
    }

    const clickedPiece = gameState.board[row][col];
    if (!IsPieceSelectable(clickedPiece, gameState)) {
        return false;
    }

    gameState.selected = { row, col };
    gameState.draggedFrom = { row, col };
    gameState.hovered = { row, col };
    gameState.isDragging = true;

    return true;
}

function HandleMouseUp2(
    dropAt: Position,
    gameState: GameState,
): InputInteractionResult {
    if (gameState.result !== null) {
        return { changed: false, moved: false, gameEnded: true };
    }

    if (!gameState.isDragging || gameState.draggedFrom === null) {
        return { changed: false, moved: false, gameEnded: false };
    }

    const from = gameState.draggedFrom;
    let changed = false;
    let moved = false;
    let gameEnded = false;

    if (dropAt !== null) {
        gameState.hovered = { row: dropAt.row, col: dropAt.col };

        const isSameSquare = from.row === dropAt.row && from.col === dropAt.col;

        if (isSameSquare) {
            gameState.selected = { row: from.row, col: from.col };
            changed = true;
        } else if (CanMovePiece(from.row, from.col, dropAt.row, dropAt.col, gameState)) {
            MovePiece(from.row, from.col, dropAt.row, dropAt.col, gameState);
            SwitchTurn(gameState);

            gameState.selected = { row: dropAt.row, col: dropAt.col };
            moved = true;
            changed = true;

            if (IsCheckmate(gameState)) {
                gameState.result =
                    gameState.turn === "white" ? "black-win" : "white-win";
                gameState.endReason = "checkmate";
                gameEnded = true;
            } else if (IsStalemate(gameState)) {
                gameState.result = "draw";
                gameState.endReason = "stalemate";
                gameEnded = true;
            }
        } else {
            // invalid drop: snap selection back to origin
            gameState.selected = { row: from.row, col: from.col };
            changed = true;
        }
    }

    gameState.isDragging = false;
    gameState.draggedFrom = null;
    changed = true;

    return { changed: changed, moved: moved, gameEnded: gameEnded };
}

function HandleMouseLeaveBoard2(gameState: GameState): boolean {
    if (gameState.hovered === null) {
        return false;
    }

    gameState.hovered = null;
    return true;
}

function HandleMouseDown(
    row: number,
    col: number,
    gameState: GameState,
    container: HTMLElement
): void {
    if (gameState.result !== null) {
        return;
    }

    const clickedPiece = gameState.board[row][col];

    if (!IsPieceSelectable(clickedPiece, gameState)) {
        return;
    }

    gameState.selected = { row, col };
    gameState.draggedFrom = { row, col };
    gameState.hovered = { row, col };
    gameState.isDragging = true;

    RenderGame(container, gameState);
}


function HandleMouseUp(
    row: number | null,
    col: number | null,
    gameState: GameState,
    container: HTMLElement
): void {
    if (gameState.result !== null) {
        return;
    }

    if (!gameState.isDragging || gameState.draggedFrom === null) {
        return;
    }

    const from = gameState.draggedFrom;


    if (row !== null && col !== null) {
        gameState.hovered = { row, col };

        const isSameSquare = from.row === row && from.col === col;

        if (!isSameSquare) {
            if (!CanMovePiece(from.row, from.col, row, col, gameState)) {
                return;
            }

            MovePiece(from.row, from.col, row, col, gameState);
            SwitchTurn(gameState);

            if (IsCheckmate(gameState)) {
                gameState.result = gameState.turn === "white" ? "black-win" : "white-win";
                gameState.endReason = "checkmate";
            } else if (IsStalemate(gameState)) {
                gameState.result = "draw";
                gameState.endReason = "stalemate";
            }
        }

        gameState.selected = { row, col };
    }

    gameState.isDragging = false;
    gameState.draggedFrom = null;

    RenderGame(container, gameState);
}

function HandleMouseLeaveBoard(
    gameState: GameState,
    container: HTMLElement
): void {
    gameState.hovered = null;

    if (!gameState.isDragging) {
        RenderGame(container, gameState);
        return;
    }

    RenderGame(container, gameState);
}

function UpdateCursor(
    container: HTMLElement,
    gameState: GameState,
    row: number,
    col: number
): void {
    if (gameState.isDragging) {
        container.style.cursor = "grabbing";
        return;
    }

    const piece = gameState.board[row][col];
    container.style.cursor = IsPieceSelectable(piece, gameState) ? "pointer" : "default";
}

function GetSquareElementFromEvent(event: MouseEvent): HTMLElement | null {
    const target = event.target as HTMLElement | null;
    if (!target) {
        return null;
    }

    return target.closest(".square") as HTMLElement | null;
}