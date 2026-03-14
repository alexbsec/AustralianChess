import { saveAuthSession } from "../auth";
import { navigateTo } from "../router";

type LoginErrorResponse = {
    error: string;
};

type LoginSuccessResponse = {
    access_token: string;
    user: {
        id: number;
        player_id: string;
    };
};

export function renderLoginPage(container: HTMLDivElement): void {
    const page = document.createElement("main");
    page.className = "auth-page";

    const card = document.createElement("section");
    card.className = "auth-card";

    const backButton = document.createElement("button");
    backButton.className = "link-button";
    backButton.textContent = "← Back to home";
    backButton.addEventListener("click", () => navigateTo("/"));

    const title = document.createElement("h1");
    title.className = "auth-title";
    title.textContent = "Login";

    const subtitle = document.createElement("p");
    subtitle.className = "auth-subtitle";
    subtitle.textContent =
        "Enter your credentials and return to the board.";

    const form = document.createElement("form");
    form.className = "auth-form";

    const usernameGroup = document.createElement("div");
    usernameGroup.className = "form-group";

    const usernameLabel = document.createElement("label");
    usernameLabel.htmlFor = "login-username";
    usernameLabel.textContent = "Username";

    const usernameInput = document.createElement("input");
    usernameInput.id = "login-username";
    usernameInput.name = "username";
    usernameInput.type = "text";
    usernameInput.placeholder = "Enter your username";
    usernameInput.autocomplete = "username";
    usernameInput.required = true;
    usernameInput.className = "form-input";

    usernameGroup.appendChild(usernameLabel);
    usernameGroup.appendChild(usernameInput);

    const passwordGroup = document.createElement("div");
    passwordGroup.className = "form-group";

    const passwordLabel = document.createElement("label");
    passwordLabel.htmlFor = "login-password";
    passwordLabel.textContent = "Password";

    const passwordInput = document.createElement("input");
    passwordInput.id = "login-password";
    passwordInput.name = "password";
    passwordInput.type = "password";
    passwordInput.placeholder = "Enter your password";
    passwordInput.autocomplete = "current-password";
    passwordInput.required = true;
    passwordInput.className = "form-input";

    passwordGroup.appendChild(passwordLabel);
    passwordGroup.appendChild(passwordInput);

    const feedback = document.createElement("div");
    feedback.className = "form-feedback hidden";

    const submitButton = document.createElement("button");
    submitButton.type = "submit";
    submitButton.className = "btn btn-primary auth-submit";
    submitButton.textContent = "Login";

    const registerHint = document.createElement("p");
    registerHint.className = "auth-footer-text";
    registerHint.innerHTML = `Don't have an account? <a href="/register">Create one</a>`;

    registerHint.querySelector("a")?.addEventListener("click", (event) => {
        event.preventDefault();
        navigateTo("/register");
    });

    function setFeedback(message: string, type: "error" | "success"): void {
        feedback.className = `form-feedback ${type}`;
        feedback.textContent = message;
    }

    function clearFeedback(): void {
        feedback.className = "form-feedback hidden";
        feedback.textContent = "";
    }

    form.addEventListener("submit", async (event) => {
        event.preventDefault();
        clearFeedback();

        const username = usernameInput.value.trim();
        const password = passwordInput.value;

        if (username.length === 0) {
            setFeedback("Username is required.", "error");
            return;
        }

        if (password.length === 0) {
            setFeedback("Password is required.", "error");
            return;
        }

        submitButton.disabled = true;
        submitButton.textContent = "Logging in...";

        try {
            const response = await fetch("/api/v1/user/login", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    username,
                    password,
                }),
            });

            const data = (await response.json()) as
                | LoginSuccessResponse
                | LoginErrorResponse;

            if (!response.ok) {
                const errorMessage =
                    "error" in data ? data.error : "Could not log in.";
                setFeedback(errorMessage, "error");
                return;
            }

            if (
                !("access_token" in data) ||
                !("user" in data) ||
                !data.user?.player_id
            ) {
                setFeedback("Unexpected server response.", "error");
                return;
            }

            saveAuthSession(data);

            setFeedback(
                `Welcome back, ${data.user.player_id}. Redirecting...`,
                "success",
            );

            setTimeout(() => {
                navigateTo("/");
            }, 900);
        } catch (error) {
            console.error(error);
            setFeedback("Network error. Please try again.", "error");
        } finally {
            submitButton.disabled = false;
            submitButton.textContent = "Login";
        }
    });

    form.appendChild(usernameGroup);
    form.appendChild(passwordGroup);
    form.appendChild(feedback);
    form.appendChild(submitButton);

    card.appendChild(backButton);
    card.appendChild(title);
    card.appendChild(subtitle);
    card.appendChild(form);
    card.appendChild(registerHint);

    page.appendChild(card);
    container.appendChild(page);
}
