import { isAuthenticated } from "../auth";
import { navigateTo } from "../router";
import { RoomController } from "../room/room_controller";
import { buildBotWsUrl } from "../ws/utils";

export function renderRoomBotPage(container: HTMLDivElement): void {
    if (!isAuthenticated()) {
        navigateTo("/login");
        return;
    }

    const params = new URLSearchParams(window.location.search);
    const roomId = params.get("id");

    if (!roomId) {
        navigateTo("/");
        return;
    }

    const difficulty = Number(sessionStorage.getItem("botDifficulty"));
    const playerPlayingAs = Number(sessionStorage.getItem("botColor") ?? "0");

    if (!difficulty) {
        navigateTo("/");
        return;
    }

    const wsUrl = buildBotWsUrl(roomId, difficulty, playerPlayingAs);

    new RoomController(container, roomId, {
        wsUrl,
        onNewGame: () => navigateTo("/"),
        newGameLabel: "Play Again",
    });
}
