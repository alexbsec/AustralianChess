import { clearAuthSession, getPlayerId, isAuthenticated } from "../auth";
import { navigateTo } from "../router";

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

async function createRoom(_playerId: string): Promise<void> {
    const response = await fetch("/api/v1/room/create", {
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

        const logoutButton = createButton(
            "Logout",
            "btn btn-secondary",
            async () => {
                clearAuthSession();
                navigateTo("/");
            },
        );

        actions.appendChild(createRoomButton);
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
