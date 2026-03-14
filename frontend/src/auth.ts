const ACCESS_TOKEN_KEY = "auth.accessToken";
const PLAYER_ID_KEY = "auth.playerId";

export type AuthUser = {
    id: number;
    player_id: string;
};

export type LoginSuccessResponse = {
    access_token: string;
    user: AuthUser;
};

export function saveAuthSession(data: LoginSuccessResponse): void {
    localStorage.setItem(ACCESS_TOKEN_KEY, data.access_token);
    localStorage.setItem(PLAYER_ID_KEY, data.user.player_id);
}

export function clearAuthSession(): void {
    localStorage.removeItem(ACCESS_TOKEN_KEY);
    localStorage.removeItem(PLAYER_ID_KEY);
}

export function getAccessToken(): string | null {
    return localStorage.getItem(ACCESS_TOKEN_KEY);
}

export function getPlayerId(): string | null {
    return localStorage.getItem(PLAYER_ID_KEY);
}

export function isAuthenticated(): boolean {
    return getAccessToken() !== null;
}

export function getAuthHeader(): Record<string, string> {
    const token = getAccessToken();

    if (!token) {
        return {};
    }

    return {
        Authorization: `Bearer ${token}`,
    };
}
