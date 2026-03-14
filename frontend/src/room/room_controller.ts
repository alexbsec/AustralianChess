import { RoomView } from "./room_view";
import { RoomSocket } from "./room_socket";
import { RoomState } from "./room_state";
import type { GameState, PieceColor } from "../engine/types";
import {
    displayToBoardPosition,
    getSquareAtPosition,
    positionsEqual
} from "../engine/coords";
import { getPseudoLegalMoves } from "../ui/ui";
import { navigateTo } from "../router";

export class RoomController {
    private view: RoomView;
    private socket: RoomSocket;
    private state: RoomState;

    constructor(container: HTMLDivElement, roomId: string, playerId: string) {
        this.state = new RoomState();
        this.view = new RoomView(container);
        this.view.updateRoomId(roomId);

        this.socket = new RoomSocket(roomId, playerId, {
            onGameState: (state, color) => this.handleGameState(state, color),
            onRoomStatus: (success, started) => this.handleRoomStatus(success, started),
            onMoveResult: (moved, state) => this.handleMoveResult(moved, state),
            onOpen: () => {
                this.view.updateStatus(
                    "Connected",
                    "",
                    "Game started",
                    "room-turn-banner--waiting"
                );
                this.view.statusText.textContent = "Connected to server";
            },
            onError: () => this.handleError(),
            onClose: () => this.handleClose(),
        });

        this.initEvents();
        this.socket.connect();
    }

    private initEvents(): void {
        this.view.boardElement.addEventListener("mousedown", (e) => this.onMouseDown(e));
        this.view.boardElement.addEventListener("mouseup", (e) => this.onMouseUp(e));

        const backBtn = document.getElementById("back-btn");
        backBtn?.addEventListener("click", () => {
            this.socket.disconnect();
            navigateTo("/");
        });
    }

    private handleGameState(state: GameState, color?: PieceColor): void {
        this.state.updateGameState(state, color);

        if (this.state.playerColor !== null) {
            const side = this.state.playerColor === 0 ? "White" : "Black";
            this.view.playerText.textContent = `Playing as ${side}`;
        }

        this.sync();
    }

    private handleRoomStatus(success: boolean, started: boolean): void {
        if (!success) {
            this.view.statusText.textContent = "Failed to join room.";
            return;
        }
        this.state.gameStarted = started;
        this.sync();
    }

    private handleMoveResult(moved: boolean, state: GameState): void {
        this.state.updateGameState(state);
        this.state.clearSelection();

        if (!moved) {
            this.view.activityText.textContent = "Illegal move rejected by server.";
        }
        this.sync();
    }

    private onMouseDown(event: MouseEvent): void {
        if (!this.state.canInteract()) return;

        const square = (event.target as HTMLElement).closest<HTMLDivElement>(".room-board-square");
        if (!square) return;

        const boardPos = displayToBoardPosition(
            { row: Number(square.dataset.row), col: Number(square.dataset.col) },
            this.state.playerColor!
        );

        const squareData = getSquareAtPosition(this.state.gameState!, boardPos);

        // Only select if it's our piece and our turn
        if (squareData?.piece?.color === this.state.playerColor && this.state.isMyTurn()) {
            this.state.uiState.selected = boardPos;
            this.state.uiState.possibleMoves = getPseudoLegalMoves(this.state.gameState!, boardPos);
            this.state.uiState.isMouseDown = true;
            this.sync();
        } else {
            this.state.clearSelection();
            this.sync();
        }
    }

    private onMouseUp(event: MouseEvent): void {
        if (!this.state.uiState.selected || this.state.movePending) return;

        const square = (event.target as HTMLElement).closest<HTMLDivElement>(".room-board-square");
        if (!square) {
            this.state.clearSelection();
            this.sync();
            return;
        }

        const destination = displayToBoardPosition(
            { row: Number(square.dataset.row), col: Number(square.dataset.col) },
            this.state.playerColor!
        );

        // Clicked the same square (toggle/cancel)
        if (positionsEqual(this.state.uiState.selected, destination)) {
            this.state.uiState.isMouseDown = false;
            return;
        }

        const isLegal = this.state.uiState.possibleMoves.some(m => positionsEqual(m, destination));
        if (isLegal) {
            this.state.movePending = true;
            this.socket.sendMove(this.state.playerColor!, this.state.uiState.selected, destination);
            this.view.activityText.textContent = "Move accepted. Waiting for opponent's response.";
        }

        this.state.clearSelection();
        this.sync();
    }

    private sync(): void {
        if (!this.state.gameState || this.state.playerColor === null) return;

        this.view.renderBoard(this.state.gameState, this.state.playerColor, this.state.uiState);
        this.updateStatusDisplays();
    }

    private updateStatusDisplays(): void {
        const { gameState, movePending } = this.state;
        const isMyTurn = this.state.isMyTurn();

        if (gameState?.result) {
            this.view.turnTitle.textContent = "Game Over";
            this.view.turnBanner.className = "room-turn-banner room-turn-banner--finished";
        } else if (movePending) {
            this.view.turnTitle.textContent = "Processing...";
            this.view.turnBanner.className = "room-turn-banner room-turn-banner--waiting";
        } else if (isMyTurn) {
            this.view.turnTitle.textContent = "Your Turn";
            this.view.activityText.textContent = "Your turn to move.";
            this.view.turnBanner.className = "room-turn-banner room-turn-banner--your-turn";
        } else {
            this.view.turnTitle.textContent = "Opponent is thinking...";
            this.view.turnBanner.className = "room-turn-banner room-turn-banner--opponent-turn";
        }
    }

    private handleError(): void {
        this.view.turnTitle.textContent = "Connection Error";
        this.view.turnBanner.className = "room-turn-banner room-turn-banner--finished";
    }

    private handleClose(): void {
        this.view.turnTitle.textContent = "Disconnected";
        this.view.turnBanner.className = "room-turn-banner room-turn-banner--finished";
    }
}
