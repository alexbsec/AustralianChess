import type { GameState, UIState } from "../engine/types";

export function CreateInitialUIState(): UIState {
    return {
        selected: null,
        hovered: null,
        draggedFrom: null,
        isDragging: false,
    }
}

export function UpdateGameOverModal(
    gameState: GameState,
    onPlayAgain: () => void
): void {
    let modal = document.getElementById("game-over-modal") as HTMLDivElement | null;

    if (!modal) {
        modal = document.createElement("div");
        modal.id = "game-over-modal";
        modal.className = "modal hidden";

        const content = document.createElement("div");
        content.className = "modal-content";

        const title = document.createElement("h2");
        title.id = "game-over-title";

        const text = document.createElement("p");
        text.id = "game-over-text";

        const button = document.createElement("button");
        button.id = "play-again-btn";
        button.type = "button";
        button.textContent = "Play Again";

        content.appendChild(title);
        content.appendChild(text);
        content.appendChild(button);
        modal.appendChild(content);
        document.body.appendChild(modal);
    }

    const title = modal.querySelector<HTMLHeadingElement>("#game-over-title");
    const text = modal.querySelector<HTMLParagraphElement>("#game-over-text");
    const button = modal.querySelector<HTMLButtonElement>("#play-again-btn");

    if (!title || !text || !button) {
        return;
    }

    button.onclick = (event) => {
        event.preventDefault();
        event.stopPropagation();
        onPlayAgain();
    };

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