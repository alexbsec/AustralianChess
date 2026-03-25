import { clearAuthSession, getPlayerId, isAuthenticated, authFetch } from "../auth";
import { navigateTo } from "../router";

// Difficulty values must match backend bot.Difficulty constants
export const BotDifficulty = {
    Easy: 3,
    Medium: 4,
    Hard: 5,
} as const;

// PieceColor values must match backend chess.PieceColor
export const PieceColor = {
    White: 0,
    Black: 1,
} as const;

function createButton(
    label: string,
    className: string,
    onClick: () => void | Promise<void>,
): HTMLButtonElement {
    const button = document.createElement("button");
    button.className = className;
    button.textContent = label;

    button.addEventListener("click", async () => {
        try {
            button.disabled = true;
            await onClick();
        } catch (error) {
            console.error(error);
            alert("Something went wrong.");
        } finally {
            button.disabled = false;
        }
    });

    return button;
}

export async function playBot(
    difficulty: number,
    playerPlayingAs: number,
): Promise<void> {
    const createResponse = await authFetch("/api/v1/room/create", {
        method: "GET",
    });

    if (!createResponse.ok) {
        throw new Error("Failed to create room");
    }

    const { roomId } = await createResponse.json();

    sessionStorage.setItem("botDifficulty", String(difficulty));
    sessionStorage.setItem("botColor", String(playerPlayingAs));

    window.location.href = `/room/bot?id=${encodeURIComponent(roomId)}`;
}

function createPlayBotDropdown(): HTMLElement {
    const wrapper = document.createElement("div");
    wrapper.className = "bot-dropdown-wrapper";
    wrapper.style.cssText = "position: relative; display: inline-block;";

    const trigger = document.createElement("button");
    trigger.className = "btn btn-primary";
    trigger.textContent = "Play Bot \u25BE";

    const panel = document.createElement("div");
    panel.className = "bot-panel";
    panel.style.cssText = `
        display: none;
        position: absolute;
        top: calc(100% + 10px);
        left: 0;
        background: var(--bg-elevated);
        border: 1px solid var(--panel-border);
        border-radius: var(--radius-sm);
        padding: 16px;
        min-width: 240px;
        box-shadow: var(--shadow);
        z-index: 100;
        backdrop-filter: blur(10px);
    `;

    // Color selection
    const colorLabel = document.createElement("p");
    colorLabel.style.cssText = "margin: 0 0 8px; font-size: 0.82rem; font-weight: 600; color: var(--muted); text-transform: uppercase; letter-spacing: 0.08em;";
    colorLabel.textContent = "Play as";

    const colorRow = document.createElement("div");
    colorRow.style.cssText = "display: flex; gap: 8px; margin-bottom: 14px;";

    let selectedColor: typeof PieceColor[keyof typeof PieceColor] = PieceColor.White;

    const whiteBtn = document.createElement("button");
    whiteBtn.style.cssText = "flex: 1; padding: 8px; border-radius: 8px; border: 2px solid var(--primary); background: rgba(255,184,77,0.15); color: var(--text); font-weight: 600; cursor: pointer; font-size: 0.9rem; transition: all 160ms ease;";
    whiteBtn.textContent = "White";

    const blackBtn = document.createElement("button");
    blackBtn.style.cssText = "flex: 1; padding: 8px; border-radius: 8px; border: 2px solid transparent; background: var(--secondary); color: var(--text); font-weight: 600; cursor: pointer; font-size: 0.9rem; transition: all 160ms ease;";
    blackBtn.textContent = "Black";

    whiteBtn.addEventListener("click", () => {
        selectedColor = PieceColor.White;
        whiteBtn.style.borderColor = "var(--primary)";
        whiteBtn.style.background = "rgba(255,184,77,0.15)";
        blackBtn.style.borderColor = "transparent";
        blackBtn.style.background = "var(--secondary)";
    });

    blackBtn.addEventListener("click", () => {
        selectedColor = PieceColor.Black;
        blackBtn.style.borderColor = "var(--primary)";
        blackBtn.style.background = "rgba(255,184,77,0.15)";
        whiteBtn.style.borderColor = "transparent";
        whiteBtn.style.background = "var(--secondary)";
    });

    colorRow.appendChild(whiteBtn);
    colorRow.appendChild(blackBtn);

    // Difficulty selection
    const diffLabel = document.createElement("p");
    diffLabel.style.cssText = "margin: 0 0 8px; font-size: 0.82rem; font-weight: 600; color: var(--muted); text-transform: uppercase; letter-spacing: 0.08em;";
    diffLabel.textContent = "Difficulty";

    const diffRow = document.createElement("div");
    diffRow.style.cssText = "display: flex; gap: 8px;";

    const difficulties: Array<{ label: string; value: number }> = [
        { label: "Easy", value: BotDifficulty.Easy },
        { label: "Medium", value: BotDifficulty.Medium },
        { label: "Hard", value: BotDifficulty.Hard },
    ];

    for (const diff of difficulties) {
        const diffBtn = document.createElement("button");
        diffBtn.style.cssText = "flex: 1; padding: 8px 4px; border-radius: 8px; border: 1px solid var(--panel-border); background: var(--secondary); color: var(--text); font-weight: 600; cursor: pointer; font-size: 0.88rem; transition: all 160ms ease;";
        diffBtn.textContent = diff.label;

        diffBtn.addEventListener("mouseenter", () => {
            diffBtn.style.background = "var(--secondary-hover)";
            diffBtn.style.borderColor = "var(--primary)";
        });
        diffBtn.addEventListener("mouseleave", () => {
            diffBtn.style.background = "var(--secondary)";
            diffBtn.style.borderColor = "var(--panel-border)";
        });

        diffBtn.addEventListener("click", async () => {
            try {
                diffBtn.disabled = true;
                panel.style.display = "none";
                await playBot(diff.value, selectedColor);
            } catch (err) {
                console.error(err);
                alert("Something went wrong.");
                diffBtn.disabled = false;
            }
        });

        diffRow.appendChild(diffBtn);
    }

    panel.appendChild(colorLabel);
    panel.appendChild(colorRow);
    panel.appendChild(diffLabel);
    panel.appendChild(diffRow);

    trigger.addEventListener("click", (e) => {
        e.stopPropagation();
        const isVisible = panel.style.display !== "none";
        panel.style.display = isVisible ? "none" : "block";
    });

    document.addEventListener("click", () => {
        panel.style.display = "none";
    });

    wrapper.appendChild(trigger);
    wrapper.appendChild(panel);
    return wrapper;
}

export async function createRoom(_playerId: string): Promise<void> {
    const response = await authFetch("/api/v1/room/create", {
        method: "GET",
    });
  

    if (!response.ok) {
        throw new Error("Failed to create room");
    }

    const data = await response.json();
    const roomId = data.roomId;

    window.location.href = `/room?id=${encodeURIComponent(roomId)}`;
}

export function renderLandingPage(container: HTMLDivElement): void {
    container.innerHTML = "";

    const page = document.createElement("main");
    page.className = "landing-page";

    const hero = document.createElement("section");
    hero.className = "hero";

    const badge = document.createElement("div");
    badge.className = "hero-badge";
    badge.textContent = "Strategy • Variant • Online";

    const title = document.createElement("h1");
    title.className = "hero-title";
    title.textContent = "Australian Chess";

    const subtitle = document.createElement("p");
    subtitle.className = "hero-subtitle";

    const actions = document.createElement("div");
    actions.className = "hero-actions";

    if (isAuthenticated()) {
        const playerId = getPlayerId() ?? "player";

        subtitle.textContent =
            `Welcome back, ${playerId}. Enter the board, create a room, and see whether memory helps more than pride.`;

        const createRoomButton = createButton(
            "Create Room",
            "btn btn-primary",
            async () => {
                await createRoom(playerId);
            },
        );

        const playBotDropdown = createPlayBotDropdown();

        const logoutButton = createButton(
            "Logout",
            "btn btn-secondary",
            async () => {
                clearAuthSession();
                navigateTo("/");
            },
        );

        actions.appendChild(createRoomButton);
        actions.appendChild(playBotDropdown);
        actions.appendChild(logoutButton);
    } else {
        subtitle.textContent =
            "Enter a sharper battlefield. Challenge opponents, create rooms, and master a chess variant where the familiar is no longer safe.";

        const createAccountButton = createButton(
            "Create Account",
            "btn btn-primary",
            async () => {
                navigateTo("/register");
            },
        );

        const loginButton = createButton(
            "Login",
            "btn btn-secondary",
            async () => {
                navigateTo("/login");
            },
        );

        actions.appendChild(createAccountButton);
        actions.appendChild(loginButton);
    }

    const infoGrid = document.createElement("section");
    infoGrid.className = "info-grid";

    const cards = [
        {
            title: "Tactical Depth",
            text: "Classic foundations twisted into a more dangerous and unpredictable form.",
        },
        {
            title: "Fast Match Setup",
            text: "Create a room instantly and jump into a live match without ceremony.",
        },
        {
            title: "Built for Rivalry",
            text: "Designed for direct competition, experimentation, and brutal mistakes.",
        },
    ];

    for (const item of cards) {
        const card = document.createElement("article");
        card.className = "info-card";

        const cardTitle = document.createElement("h2");
        cardTitle.className = "info-card-title";
        cardTitle.textContent = item.title;

        const cardText = document.createElement("p");
        cardText.className = "info-card-text";
        cardText.textContent = item.text;

        card.appendChild(cardTitle);
        card.appendChild(cardText);
        infoGrid.appendChild(card);
    }

    hero.appendChild(badge);
    hero.appendChild(title);
    hero.appendChild(subtitle);
    hero.appendChild(actions);

    page.appendChild(hero);
    page.appendChild(infoGrid);

    container.appendChild(page);
}
