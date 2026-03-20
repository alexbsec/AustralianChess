import { RoomView } from "./room_view";
import { RoomSocket } from "./room_socket";
import { RoomState } from "./room_state";
import type { GameState, PieceColor, PieceKind, Position } from "../engine/types";
import {
    displayToBoardPosition,
    getSquareAtPosition,
    positionsEqual
} from "../engine/coords";
import { getPseudoLegalMoves } from "../ui/ui";
import { navigateTo } from "../router";
import { createRoom } from "../pages/landing";
import { getPlayerId } from "../auth";

export class RoomController {
    private view: RoomView;
    private socket: RoomSocket;
    private state: RoomState;
    private modalShown: boolean = false;
    private moveSound: HTMLAudioElement;
    private captureSound: HTMLAudioElement;
    private checkSound: HTMLAudioElement;
    private optimisticState: GameState | null = null;
    private draggedElement: HTMLElement | null = null;
    private audioUnlocked: boolean = false;

    constructor(container: HTMLDivElement, roomId: string) {
        this.state = new RoomState();
        this.view = new RoomView(container);
        this.view.updateRoomId(roomId);

        this.socket = new RoomSocket(roomId, {
            onGameState: (state, color) => this.handleGameState(state, color),
            onRoomStatus: (success, started) => this.handleRoomStatus(success, started),
            onMoveResult: (moved, state) => {
                if (moved) {
                    this.playCorrectSound(state);
                }
                this.handleMoveResult(moved, state);
            },
            onPromoteResult: (promoted, state) => {
                if (promoted) {
                    this.playCorrectSound(state);
                }
                this.handleMoveResult(promoted, state);
            },
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

        this.moveSound = new Audio();
        this.moveSound.src = "/sfx/piece_move.mp3";
        this.moveSound.load(); // Force the browser to start downloading
        this.moveSound.volume = 0.5;

        this.captureSound = new Audio();
        this.captureSound.src = "/sfx/piece_capture.mp3";
        this.captureSound.load();
        this.captureSound.volume = 0.5;

        this.checkSound = new Audio();
        this.checkSound.src = "/sfx/move_check.mp3";
        this.checkSound.load();
        this.checkSound.volume = 0.5;

        this.initEvents();
        this.socket.connect();
    }

    private initEvents(): void {
        this.view.boardElement.addEventListener("touchstart", (e) => this.onTouchStart(e), { passive: false });
        window.addEventListener("touchmove", (e) => this.onTouchMove(e), { passive: false });
        window.addEventListener("touchend", (e) => this.onTouchEnd(e));
        this.view.boardElement.addEventListener("mousedown", (e) => this.onMouseDown(e));
        this.view.boardElement.addEventListener("mouseup", (e) => this.onMouseUp(e));

        window.addEventListener("mousemove", (e) => this.onMouseMove(e));
        this.view.boardElement.addEventListener("mouseup", (e) => this.onMouseUp(e));

        const backBtn = document.getElementById("back-btn");
        backBtn?.addEventListener("click", () => {
            this.socket.disconnect();
            navigateTo("/");
        });
    }

    private onTouchStart(event: TouchEvent): void {
        event.preventDefault(); // stops page scroll
        const touch = event.touches[0];
        this.onMouseDown(new MouseEvent("mousedown", {
            clientX: touch.clientX,
            clientY: touch.clientY,
            bubbles: true
        }));
    }

    private onTouchMove(event: TouchEvent): void {
        event.preventDefault();
        const touch = event.touches[0];
        this.onMouseMove(new MouseEvent("mousemove", {
            clientX: touch.clientX,
            clientY: touch.clientY
        }));
    }

    private onTouchEnd(event: TouchEvent): void {
        const touch = event.changedTouches[0]; // changedTouches, not touches (finger is lifted)
        this.onMouseUp(new MouseEvent("mouseup", {
            clientX: touch.clientX,
            clientY: touch.clientY,
            bubbles: true
        }));
    }

    private playCorrectSound(newState: GameState): void {
        const oldState = this.state.gameState;
        if (!oldState) return;

        const oldPieceCount = oldState.board.data.flat().filter(sq => sq.piece !== null).length;
        const newPieceCount = newState.board.data.flat().filter(sq => sq.piece !== null).length;
        const isCapture = newPieceCount < oldPieceCount;

        const isCheck = newState.in_check;

        if (isCheck) {
            this.playSound(this.checkSound);
        } else if (isCapture) {
            this.playSound(this.captureSound);
        } else {
            this.playSound(this.moveSound);
        }
    }

    private playSound(audio: HTMLAudioElement): void {
        audio.currentTime = 0;
        audio.play().catch(err => {
            console.warn("Audio blocked:", err);
        });
    }

    private handleGameState(state: GameState, color?: PieceColor): void {
        const isSubsequentMove = this.state.gameState !== null;
        const oldTurn = this.state.gameState?.turn;
        const oldState = this.state.gameState;

        if (isSubsequentMove && oldTurn !== state.turn && Number(state.turn) === Number(this.state.playerColor)) {
            this.playCorrectSound(state);
        }

        console.log("old state:", oldState);
        if (oldState) {
            this.state.uiState.lastMove = this.detectLastMove(oldState, state);
        }

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
        const oldState = this.state.gameState;
        const wasMyMove = this.state.movePending;

        if (moved && oldState && !wasMyMove) {
            // opponent's move — detect from diff
            this.state.uiState.lastMove = this.detectLastMove(oldState, state);
        }

        this.state.updateGameState(state);
        this.state.uiState.draggingPos = null;
        this.state.clearSelection();
        this.state.movePending = false;

        if (!moved) {
            this.view.activityText.textContent = "Illegal move rejected by server.";
            this.state.uiState.lastMove = null;
        }
        this.sync();
    }

    private onMouseDown(event: MouseEvent): void {
        if (!this.audioUnlocked) {
            [this.moveSound, this.captureSound].forEach(sound => {
                const vol = sound.volume;
                sound.volume = 0;
                sound.play().then(() => {
                    sound.pause();
                    sound.volume = vol;
                });
            });
            this.audioUnlocked = true;
        }

        if (!this.state.canInteract()) return;

        const square = (event.target as HTMLElement).closest<HTMLDivElement>(".room-board-square");
        if (!square) return;

        const pieceImg = square.querySelector<HTMLImageElement>(".room-piece");
        if (!pieceImg) return;

        const boardPos = displayToBoardPosition(
            { row: Number(square.dataset.row), col: Number(square.dataset.col) },
            this.state.playerColor!
        );

        const squareData = getSquareAtPosition(this.state.gameState!, boardPos);

        if (squareData?.piece?.color === this.state.playerColor && this.state.isMyTurn()) {
            this.state.uiState.selected = boardPos;
            this.state.uiState.possibleMoves = getPseudoLegalMoves(this.state.gameState!, boardPos);
            this.state.uiState.draggingPos = boardPos;

            this.draggedElement = pieceImg.cloneNode(true) as HTMLElement;
            this.draggedElement.classList.add("dragging-piece");

            Object.assign(this.draggedElement.style, {
                position: 'fixed',
                width: `${pieceImg.offsetWidth}px`,
                height: `${pieceImg.offsetHeight}px`,
                pointerEvents: 'none',
                zIndex: '1000',
                left: `${event.clientX - pieceImg.offsetWidth / 2}px`,
                top: `${event.clientY - pieceImg.offsetHeight / 2}px`
            });

            document.body.appendChild(this.draggedElement);
            this.state.uiState.isMouseDown = true;
            this.sync();
        }
    }

    private onMouseMove(event: MouseEvent): void {
        if (!this.draggedElement) return;
        this.draggedElement.style.left = `${event.clientX - this.draggedElement.offsetWidth / 2}px`;
        this.draggedElement.style.top = `${event.clientY - this.draggedElement.offsetHeight / 2}px`;
    }

    private onMouseUp(event: MouseEvent): void {
        const selectedPos = this.state.uiState.selected;

        if (this.draggedElement) {
            this.draggedElement.remove();
            this.draggedElement = null;
        }

        if (!selectedPos || this.state.movePending) {
            this.state.uiState.isMouseDown = false;
            return;
        }

        const square = (event.target as HTMLElement).closest<HTMLDivElement>(".room-board-square");
        if (!square) {
            this.cancelDrag();
            return;
        }

        const destination = displayToBoardPosition(
            { row: Number(square.dataset.row), col: Number(square.dataset.col) },
            this.state.playerColor!
        );

        if (positionsEqual(selectedPos, destination)) {
            this.cancelDrag();
            return;
        }

        const isLegal = this.state.uiState.possibleMoves.some(m => positionsEqual(m, destination));

        if (!isLegal) {
            this.cancelDrag();
            this.state.uiState.isMouseDown = false;
            return
        }

        const isPromotion = this.isPawnPromotion(selectedPos, destination);
        if (isPromotion) {
            this.view.showPromotionPicker(this.state.playerColor!, (chosenPiece) => {
                this.applyOptimisticMove(selectedPos, destination, null);
                this.state.movePending = true;
                this.socket.sendPromote(this.state.playerColor!, chosenPiece, selectedPos, destination);
                this.view.activityText.textContent = "Promoting...";
                this.sync();
            });
        } else {
            this.applyOptimisticMove(selectedPos, destination, null);
            this.state.movePending = true;
            this.socket.sendMove(this.state.playerColor!, selectedPos, destination);
            this.view.activityText.textContent = "Promoting...";
            this.sync();
        }


        if (this.optimisticState) { }

        this.state.uiState.isMouseDown = false;
    }

    private isPawnPromotion(from: Position, to: Position): boolean {
        const piece = getSquareAtPosition(this.state.gameState!, from)?.piece;
        if (!piece || piece?.kind !== 0) return false;
        const backRank = this.state.playerColor === 0 ? 0 : 11;
        return to.row === backRank;
    }

    private detectLastMove(oldState: GameState, newState: GameState): { from: Position; to: Position } | null {
        let from: Position | null = null;
        let to: Position | null = null;

        for (let r = 0; r < oldState.board.data.length; r++) {
            for (let c = 0; c < oldState.board.data[r].length; c++) {
                const oldPiece = oldState.board.data[r][c].piece;
                const newPiece = newState.board.data[r][c].piece;

                if (oldPiece && !newPiece) {
                    from = { row: r, col: c };
                }

                if (!oldPiece && newPiece) {
                    to = { row: r, col: c };
                }
            }
        }

        console.log("diff -> from:", from, "to:", to);
        return from && to ? { from, to } : null;
    }

    private applyOptimisticMove(from: Position, to: Position, promoteTo: PieceKind | null): void {
        this.optimisticState = structuredClone(this.state.gameState!);

        const board = this.state.gameState!.board.data;
        const piece = board[from.row][from.col].piece;
        board[to.row][to.col].piece = promoteTo !== null ? { kind: promoteTo, color: this.state.playerColor! } : piece;
        board[from.row][from.col].piece = null;

        this.state.uiState.selected = null;
        this.state.uiState.possibleMoves = [];
        this.state.uiState.draggingPos = null;
        this.state.uiState.lastMove = { from, to };
    }

    private sync(): void {
        if (!this.state.gameState || this.state.playerColor === null) return;
        this.view.renderBoard(this.state.gameState, this.state.playerColor, this.state.uiState);
        this.updateStatusDisplays();
    }

    private cancelDrag(): void {
        this.state.uiState.draggingPos = null;
        this.state.clearSelection();
        this.state.uiState.isMouseDown = false;
        this.sync();
    }

    private updateStatusDisplays(): void {
        const { gameState, movePending, gameStarted } = this.state;
        if (!gameState) return;

        const hasEnded = gameState.result !== null || gameState.end_reason !== null;

        if (hasEnded && !this.modalShown) {
            const reason = gameState.end_reason || "Match finished.";
            console.log("firing game over modal with reason:", reason);
            this.view.statusText.textContent = "Match finished.";
            this.view.activityText.textContent = reason;
            this.view.turnTitle.textContent = "Game Finished";
            this.view.turnBanner.className = "room-turn-banner room-turn-banner--finished";

            this.view.showGameOverModal(reason, async () => {
                this.socket.disconnect();
                try {
                    const playerId = getPlayerId();
                    if (playerId !== null) {
                        await createRoom(playerId);
                    } else {
                        navigateTo("/");
                    }
                } catch (error) {
                    navigateTo("/");
                }
            });

            console.log("Game ended with reason:", reason);
            this.modalShown = true;
            return;
        }

        if (!gameStarted) {
            this.view.turnTitle.textContent = "Awaiting for an opponent to joing the room...";
            this.view.activityText.textContent = "Awaiting for an opponent to join the room.";
            this.view.turnBanner.className = "room-turn-banner room-turn-banner--waiting";
            return;
        }

        const isMyTurn = this.state.isMyTurn();

        if (movePending) {
            this.view.turnTitle.textContent = "Processing...";
            this.view.turnBanner.className = "room-turn-banner room-turn-banner--waiting";
        } else if (isMyTurn) {
            this.view.turnTitle.textContent = "Your Turn";
            this.view.activityText.textContent = "Your turn to move.";
            this.view.turnBanner.className = "room-turn-banner room-turn-banner--your-turn";
        } else {
            this.view.turnTitle.textContent = "Opponent is thinking...";
            this.view.activityText.textContent = "Waiting for opponent...";
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
