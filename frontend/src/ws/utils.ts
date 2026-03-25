import type { WsMessage, GameStateMessage, RoomStatusMessage, MoveResultMessage, PromoteResultMessage, InactivityWarningMessage } from "./types";
import { getAccessToken } from "../auth";

export function isGameStateMessage(message: WsMessage): message is GameStateMessage {
    return "type" in message && message.type === "game_state";
}

export function isRoomStatusMessage(message: WsMessage): message is RoomStatusMessage {
    return "room_id" in message && "game_started" in message;
}

export function isMoveResultMessage(message: WsMessage): message is MoveResultMessage {
    return "moved" in message && "game_state" in message && "room_id" in message;
}

export function isPromotionResultMessage(message: WsMessage): message is PromoteResultMessage {
    return "promoted" in message && "game_state" in message && "room_id" in message;
}

export function isInactivityWarningMessage(message: WsMessage): message is InactivityWarningMessage {
    return "seconds" in message;
}

export function buildRoomWsUrl(roomId: string): string {
    const token = getAccessToken();
    const url = new URL(
        `${getWsProtocol()}//${window.location.host}/api/v1/room/ws/${encodeURIComponent(roomId)}?token=${encodeURIComponent(token ?? "")}`,
    );

    return url.toString();
}

export function buildBotWsUrl(roomId: string, difficulty: number, playerPlayingAs: number): string {
    const token = getAccessToken();
    const url = new URL(`${getWsProtocol()}//${window.location.host}/api/v1/room/ws/bot`);
    url.searchParams.set("token", token ?? "");
    url.searchParams.set("roomId", roomId);
    url.searchParams.set("difficulty", String(difficulty));
    url.searchParams.set("player_playing_as", String(playerPlayingAs));
    return url.toString();
}

function getWsProtocol(): "ws:" | "wss:" {
    return window.location.protocol === "https:" ? "wss:" : "ws:";
}
