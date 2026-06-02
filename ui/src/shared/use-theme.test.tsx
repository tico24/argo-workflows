import {act, renderHook} from '@testing-library/react';

import {applyThemeClass, resolveTheme, systemPrefersDark, THEME_KEY, useTheme} from './use-theme';

// Minimal matchMedia mock that lets a test flip the OS preference and fire the
// 'change' event registered by useTheme.
function mockMatchMedia(initialDark: boolean) {
    let matches = initialDark;
    const listeners: Array<(e: MediaQueryListEvent) => void> = [];
    const mql = {
        get matches() {
            return matches;
        },
        media: '(prefers-color-scheme: dark)',
        addEventListener: (_: string, cb: (e: MediaQueryListEvent) => void) => listeners.push(cb),
        removeEventListener: (_: string, cb: (e: MediaQueryListEvent) => void) => {
            const i = listeners.indexOf(cb);
            if (i >= 0) listeners.splice(i, 1);
        },
        addListener: (cb: (e: MediaQueryListEvent) => void) => listeners.push(cb),
        removeListener: (cb: (e: MediaQueryListEvent) => void) => {
            const i = listeners.indexOf(cb);
            if (i >= 0) listeners.splice(i, 1);
        }
    };
    (window as any).matchMedia = jest.fn().mockReturnValue(mql);
    return {
        setDark(v: boolean) {
            matches = v;
            listeners.forEach(cb => cb({matches: v} as MediaQueryListEvent));
        }
    };
}

describe('resolveTheme', () => {
    it('explicit light/dark override always wins regardless of OS', () => {
        expect(resolveTheme('light', true)).toBe('light');
        expect(resolveTheme('dark', false)).toBe('dark');
    });

    it('auto follows the OS preference', () => {
        expect(resolveTheme('auto', true)).toBe('dark');
        expect(resolveTheme('auto', false)).toBe('light');
    });
});

describe('applyThemeClass', () => {
    it('sets exactly one theme class on the document root', () => {
        applyThemeClass('dark');
        expect(document.documentElement.classList.contains('theme-dark')).toBe(true);
        expect(document.documentElement.classList.contains('theme-light')).toBe(false);
        applyThemeClass('light');
        expect(document.documentElement.classList.contains('theme-light')).toBe(true);
        expect(document.documentElement.classList.contains('theme-dark')).toBe(false);
    });
});

describe('useTheme', () => {
    beforeEach(() => {
        window.localStorage.clear();
        document.documentElement.classList.remove('theme-light', 'theme-dark');
    });

    it('defaults to auto and follows a dark OS on first load', () => {
        mockMatchMedia(true);
        const {result} = renderHook(() => useTheme());
        const [mode, , resolved] = result.current;
        expect(mode).toBe('auto');
        expect(resolved).toBe('dark');
        expect(document.documentElement.classList.contains('theme-dark')).toBe(true);
    });

    it('defaults to auto and follows a light OS on first load', () => {
        mockMatchMedia(false);
        const {result} = renderHook(() => useTheme());
        expect(result.current[2]).toBe('light');
        expect(document.documentElement.classList.contains('theme-light')).toBe(true);
    });

    it('persists an explicit override and restores it on reload', () => {
        mockMatchMedia(true); // OS is dark...
        const first = renderHook(() => useTheme());
        act(() => first.result.current[1]('light')); // ...but user picks light
        expect(first.result.current[2]).toBe('light');
        expect(JSON.parse(window.localStorage.getItem(THEME_KEY)!)).toBe('light');

        // Simulate reload: a fresh hook reads the stored override and ignores OS.
        const second = renderHook(() => useTheme());
        expect(second.result.current[0]).toBe('light');
        expect(second.result.current[2]).toBe('light');
    });

    it('reacts live to OS changes while in auto mode', () => {
        const os = mockMatchMedia(false);
        const {result} = renderHook(() => useTheme());
        expect(result.current[2]).toBe('light');
        act(() => os.setDark(true));
        expect(result.current[2]).toBe('dark');
        expect(document.documentElement.classList.contains('theme-dark')).toBe(true);
    });

    it('does NOT follow OS changes once an explicit override is set', () => {
        const os = mockMatchMedia(false);
        const {result} = renderHook(() => useTheme());
        act(() => result.current[1]('dark')); // explicit dark
        act(() => os.setDark(true)); // OS goes dark too
        act(() => os.setDark(false)); // OS goes light — override must hold
        expect(result.current[2]).toBe('dark');
    });
});

describe('systemPrefersDark', () => {
    it('returns the matchMedia result', () => {
        mockMatchMedia(true);
        expect(systemPrefersDark()).toBe(true);
        mockMatchMedia(false);
        expect(systemPrefersDark()).toBe(false);
    });
});
