import { getPlayerId } from "../auth";
import { RoomController } from "../room/room_controller";

/**
 * Entry point for the Room page.
 * Responsibilities:
 * 1. Extract roomId from the URL.
 * 2. Validate session/auth.
 * 3. Initialize the Controller.
 */
export function renderRoomPage(container: HTMLDivElement): void {
    const urlParams = new URLSearchParams(window.location.search);
    const roomId = urlParams.get("id");
    const playerId = getPlayerId();

    if (!roomId) {
        container.innerHTML = `
            <div class="error-container">
                <h2>Invalid Room</h2>
                <p>No room ID provided. Returning to lobby...</p>
            </div>
        `;
        setTimeout(() => window.location.hash = "/", 2000);
        return;
    }

    if (!playerId) {
        container.innerHTML = `
            <div class="error-container">
                <h2>Authentication Required</h2>
                <p>Please log in to join this room.</p>
            </div>
        `;
        return;
    }

    // The Controller takes over and handles the view, socket, and state.
    new RoomController(container, roomId, playerId);
}
