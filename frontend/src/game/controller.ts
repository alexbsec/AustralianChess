import { IsPieceSelectable } from "../engine/board";
import { EndTurn } from "../engine/engine";
import { CreateInitialBoard } from "../engine/initialPosition";
import { TryApplyMove } from "../engine/movements";
import { GetPseudoLegalMoves, RenderBoard } from "../engine/render";
import type { GameState, Position, UIState } from "../engine/types";
import { AttachBoardInput } from "../input/input";
import { CreateInitialUIState, UpdateGameOverModal } from "../ui/ui";

export class GameController {
    private m_GameState: GameState;
    private m_Container: HTMLElement;
    private m_UIState: UIState;

    constructor(container: HTMLElement) {
        this.m_GameState = CreateInitialGameState()
        this.m_Container = container;
        this.m_UIState = CreateInitialUIState();
    }

    start(): void {
        this.render();
        AttachBoardInput(this.m_Container, {
            onMouseDown: (pos) => this.handleMouseDown(pos),
            onMouseMove: (pos) => this.handleMouseMove(pos),
            onMouseUp: (pos) => this.handleMouseUp(pos),
            onMouseLeave: () => this.handleMouseLeave(),
        });
    }

    private render(): void {
        const possibleMoves = this.m_UIState.selected === null ? [] : GetPseudoLegalMoves(this.m_GameState, this.m_UIState.selected.row,
            this.m_UIState.selected.col);

        RenderBoard(this.m_GameState.board, this.m_Container, this.m_UIState, possibleMoves);
        UpdateGameOverModal(this.m_GameState, () => this.restartGame());
    }

    private handleMouseDown(pos: Position | null): void {
        if (this.m_GameState.result !== null) {
            return;
        }

        if (pos === null) {
            return;
        }

        const piece = this.m_GameState.board[pos.row][pos.col];
        if (!IsPieceSelectable(piece, this.m_GameState)) {
            return;
        }

        this.m_UIState.selected = pos;
        this.m_UIState.draggedFrom = pos;
        this.m_UIState.hovered = pos;
        this.m_UIState.isDragging = true;

        this.render();
    }

    private handleMouseMove(pos: Position | null): void {
        if (pos === null) {
            if (this.m_UIState.hovered !== null) {
                this.m_UIState.hovered = null;
                this.render();
            }
            return;
        }

        if (
            this.m_UIState.hovered !== null &&
            this.m_UIState.hovered.row === pos.row &&
            this.m_UIState.hovered.col === pos.col
        ) {
            return;
        }

        this.m_UIState.hovered = pos;
        this.render();
    }

    private handleMouseUp(pos: Position | null): void {
        if (this.m_GameState.result !== null) {
            return;
        }

        if (!this.m_UIState.isDragging || this.m_UIState.draggedFrom === null) {
            return;
        }

        if (pos === null) {
            this.resetUIAndRender();
            return;
        }

        const from = this.m_UIState.draggedFrom;
        this.m_UIState.hovered = pos;

        const isSameSquare = from.row === pos.row && from.col === pos.col;
        if (isSameSquare) {
            this.m_UIState.selected = from;
            this.resetUIAndRender();
            return;
        }

        const nextState = TryApplyMove(this.m_GameState, from, pos);

        if (nextState === null) {
            this.m_UIState.selected = from;
            this.resetUIAndRender();
            return;
        }

        this.m_GameState = EndTurn(nextState);
        this.m_UIState.selected = pos;
        this.resetUIAndRender();    
    }

    private handleMouseLeave(): void {
        if (this.m_UIState.hovered === null && !this.m_UIState.isDragging) {
            return;
        }

        this.m_UIState.hovered = null;
        this.render();
    }

    private resetUIAndRender(): void {
        this.m_UIState.isDragging = false;
        this.m_UIState.draggedFrom = null;
        this.render();
    }

    private restartGame(): void {
        this.m_GameState = CreateInitialGameState();
        this.m_UIState = CreateInitialUIState();
        this.render();
    }
};

function CreateInitialGameState(): GameState {
    return {
        board: CreateInitialBoard(),
        turn: "white",
        result: null,
        endReason: null,
    };
}