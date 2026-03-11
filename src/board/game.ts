import { CreateInitialBoard } from "./initialPosition";
import { CanMovePiece, IsCheckmate, IsStalemate } from "./movements";
import { GetPseudoLegalMoves, RenderBoard } from "./render";
import type { GameState, Position, Piece } from "./types";

export function CreateInitialGameState(): GameState {
    return {
        board: CreateInitialBoard(),
        turn: "white",
        selected: null,
        hovered: null,
        draggedFrom: null,
        isDragging: false,
        result: null,
        endReason: null,
    };
}

export function InitializeBoard(container: HTMLElement): void {
    const gameState = CreateInitialGameState();

    RenderGame(container, gameState);

    container.addEventListener("mousedown", (event: MouseEvent) => {
        const square = GetSquareElementFromEvent(event);
        if (!square) {
            return;
        }

        const position = GetPositionFromSquare(square);
        if (!position) {
            return;
        }

        HandleMouseDown(position.row, position.col, gameState, container);
    });

    container.addEventListener("mousemove", (event: MouseEvent) => {
        const square = GetSquareElementFromEvent(event);

        if (!square) {
            HandleMouseLeaveBoard(gameState, container);
            return;
        }

        const position = GetPositionFromSquare(square);
        if (!position) {
            return;
        }

        HandleMouseMove(position.row, position.col, gameState, container);
    });

    container.addEventListener("mouseup", (event: MouseEvent) => {
        const square = GetSquareElementFromEvent(event);

        if (!square) {
            HandleMouseUp(null, null, gameState, container);
            return;
        }

        const position = GetPositionFromSquare(square);
        if (!position) {
            HandleMouseUp(null, null, gameState, container);
            return;
        }

        HandleMouseUp(position.row, position.col, gameState, container);
    });

    container.addEventListener("mouseleave", () => {
        HandleMouseLeaveBoard(gameState, container);
    });
}

function RenderGame(container: HTMLElement, gameState: GameState): void {
    const possibleMoves =
        gameState.selected === null
            ? []
            : GetPseudoLegalMoves(
                  gameState,
                  gameState.selected.row,
                  gameState.selected.col
              );

    RenderBoard(gameState.board, container, gameState, possibleMoves);

    UpdateGameOverModal(gameState);
}

function UpdateGameOverModal(gameState: GameState): void {
    let modal = document.getElementById("game-over-modal") as HTMLDivElement | null;

    if (!modal) {
        modal = document.createElement("div");
        modal.id = "game-over-modal";
        modal.className = "game-over-modal hidden";

        const content = document.createElement("div");
        content.className = "game-over-modal-content";

        const title = document.createElement("h2");
        title.id = "game-over-title";

        const text = document.createElement("p");
        text.id = "game-over-text";

        const button = document.createElement("button");
        button.id = "play-again";
        button.type = "button";
        button.textContent = "Play Again";

        button.onpointerdown = (event) => {
            event.preventDefault();
            event.stopPropagation();
            location.replace(location.href);
        };

        content.appendChild(title);
        content.appendChild(text);
        content.appendChild(button);
        modal.appendChild(content);
        document.body.appendChild(modal);
    }

    const title = modal.querySelector<HTMLHeadingElement>("#game-over-title");
    const text = modal.querySelector<HTMLParagraphElement>("#game-over-text");

    if (!title || !text) {
        return;
    }

    if (gameState.result === null) {
        modal.classList.add("hidden");
        return;
    }

    if (gameState.endReason === "checkmate") {
        title.textContent = "Checkmate";
        text.textContent =
            gameState.result === "white-win"
                ? "White wins."
                : gameState.result === "black-win"
                ? "Black wins."
                : "Draw.";
    } else {
        title.textContent = "Stalemate";
        text.textContent = "Draw.";
    }

    modal.classList.remove("hidden");
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

function HandleMouseMove(
    row: number,
    col: number,
    gameState: GameState,
    container: HTMLElement
): void {
    const nextHovered = { row, col };

    if (ArePositionsEqual(gameState.hovered, nextHovered)) {
        return;
    }

    gameState.hovered = nextHovered;
    UpdateCursor(container, gameState, row, col);
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
    container.style.cursor = IsPieceSelectable(piece, gameState) ? "grab" : "default";
}

function SwitchTurn(gameState: GameState): void {
    gameState.turn = gameState.turn === "white" ? "black" : "white";
}

function IsPieceSelectable(
    piece: Piece | null,
    gameState: GameState
): boolean {
    if (piece === null) {
        return false;
    }

    return piece.color === gameState.turn;
}

function GetSquareElementFromEvent(event: MouseEvent): HTMLElement | null {
    const target = event.target as HTMLElement | null;
    if (!target) {
        return null;
    }

    return target.closest(".square") as HTMLElement | null;
}

function GetPositionFromSquare(square: HTMLElement): Position | null {
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

function ArePositionsEqual(
    a: Position | null,
    b: Position | null
): boolean {
    if (a === null || b === null) {
        return false;
    }

    return a.row === b.row && a.col === b.col;
}

function MovePiece(
    fromRow: number,
    fromCol: number,
    toRow: number,
    toCol: number,
    gameState: GameState
): void {
    const piece = gameState.board[fromRow][fromCol];

    if (piece === null) {
        return;
    }

    gameState.board[toRow][toCol] = piece;
    gameState.board[fromRow][fromCol] = null;
}