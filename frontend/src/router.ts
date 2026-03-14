export function navigateTo(path: string): void {
    window.history.pushState({}, "", path);
    window.dispatchEvent(new Event("app:navigate"));
}
