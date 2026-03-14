import type { 
    GameState, 
    PieceColor, 
} from "../engine/types";

import type { UIState } from "../ui/types";

/**
 * Manages the reactive state of the room.
 * This class holds the data while the Controller decides when to update it.
 */
export class RoomState {
    public gameState: GameState | null = null;
    public playerColor: PieceColor | null = null;
    public gameStarted: boolean = false;
    public movePending: boolean = false;
    
    public uiState: UIState = {
        selected: null,
        possibleMoves: [],
        isMouseDown: false,
    };

    /**
     * Resets the selection and movement UI state.
     */
    public clearSelection(): void {
        this.uiState.selected = null;
        this.uiState.possibleMoves = [];
        this.uiState.isMouseDown = false;
    }

    /**
     * Updates the current game state and resets pending flags.
     */
    public updateGameState(state: GameState, color?: PieceColor): void {
        this.gameState = state;
        if (typeof color === "number") {
            this.playerColor = color;
        }
        this.movePending = false;
    }

    /**
     * Helper to check if the current user is allowed to make a move.
     */
    public isMyTurn(): boolean {
        if (!this.gameState || this.playerColor === null || !this.gameStarted) {
            return false;
        }
        return this.gameState.turn === this.playerColor;
    }

    /**
     * Checks if interaction with the board is currently possible.
     */
    public canInteract(): boolean {
        return (
            this.gameStarted && 
            this.gameState !== null && 
            this.playerColor !== null && 
            !this.movePending &&
            !this.gameState.result
        );
    }
}
