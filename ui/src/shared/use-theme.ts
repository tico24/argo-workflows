import {useEffect, useState} from 'react';

import {useLocalStorage} from './use-local-storage';

// FEATURE: optional dark mode.
//
// ONE shared localStorage key for the user's theme preference. JSON-encoded by
// useLocalStorage, mirroring the modal-suppression / filters-collapsed
// conventions already used in the codebase. Value is one of the ThemeMode
// strings below.
//
// Resolution order:
//   1. explicit user override ('light' | 'dark') always wins;
//   2. 'auto' (the DEFAULT) follows the OS via prefers-color-scheme, and reacts
//      live to OS changes.
// Default mode is 'auto' so first load (with no stored choice) follows the OS
// and the light baseline is byte-for-byte preserved on a light OS.
export const THEME_KEY = 'theme';

export type ThemeMode = 'light' | 'dark' | 'auto';
export type ResolvedTheme = 'light' | 'dark';

const DARK_QUERY = '(prefers-color-scheme: dark)';

// systemPrefersDark reads the OS preference synchronously. Guarded for non-DOM
// (test/SSR) environments where matchMedia may be absent.
export function systemPrefersDark(): boolean {
    if (typeof window === 'undefined' || typeof window.matchMedia !== 'function') {
        return false;
    }
    return window.matchMedia(DARK_QUERY).matches;
}

// resolveTheme maps a stored mode + the current OS preference onto the concrete
// theme that should be painted.
export function resolveTheme(mode: ThemeMode, prefersDark: boolean): ResolvedTheme {
    if (mode === 'light' || mode === 'dark') {
        return mode;
    }
    return prefersDark ? 'dark' : 'light';
}

// applyThemeClass writes the argo-ui theme ancestor class onto the document
// root so that legacy themify() rules AND the restyle's `.theme-dark` token
// overrides resolve — including content rendered OUTSIDE the argo-ui Layout
// (login, widgets, top-level Popup, notifications, portals).
export function applyThemeClass(resolved: ResolvedTheme): void {
    if (typeof document === 'undefined') {
        return;
    }
    const root = document.documentElement;
    root.classList.remove('theme-light', 'theme-dark');
    root.classList.add('theme-' + resolved);
}

// useTheme centralises mode persistence, OS-preference resolution, live OS
// change reaction (auto mode only), and applying the root ancestor class.
// Returns [mode, setMode, resolved]. Kept high in app-router so toggling the
// theme never remounts route/filter components.
export function useTheme(): [ThemeMode, (mode: ThemeMode) => void, ResolvedTheme] {
    const [mode, setMode] = useLocalStorage<ThemeMode>(THEME_KEY, 'auto');
    const [prefersDark, setPrefersDark] = useState<boolean>(() => systemPrefersDark());

    // React to OS changes so 'auto' mode follows the system live.
    useEffect(() => {
        if (typeof window === 'undefined' || typeof window.matchMedia !== 'function') {
            return undefined;
        }
        const mql = window.matchMedia(DARK_QUERY);
        const onChange = (e: MediaQueryListEvent) => setPrefersDark(e.matches);
        // addEventListener is the modern API; older Safari only has addListener.
        if (typeof mql.addEventListener === 'function') {
            mql.addEventListener('change', onChange);
            return () => mql.removeEventListener('change', onChange);
        }
        mql.addListener(onChange);
        return () => mql.removeListener(onChange);
    }, []);

    const resolved = resolveTheme(mode, prefersDark);

    // Apply the class to the document root whenever the resolved theme changes,
    // so routes rendered outside argo-ui's Layout are themed too.
    useEffect(() => {
        applyThemeClass(resolved);
    }, [resolved]);

    return [mode, setMode, resolved];
}
