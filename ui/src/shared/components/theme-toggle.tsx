import classNames from 'classnames';
import * as React from 'react';

import {ResolvedTheme, ThemeMode} from '../use-theme';

interface ThemeToggleProps {
    mode: ThemeMode;
    resolved: ResolvedTheme;
    onChange: (mode: ThemeMode) => void;
}

// The three explicit modes the toggle cycles through, in order. 'auto' follows
// the OS preference (the default); 'light' / 'dark' are explicit user overrides
// that persist and win over the OS.
const ORDER: ThemeMode[] = ['auto', 'light', 'dark'];

const ICON: Record<ThemeMode, string> = {
    auto: 'fa-circle-half-stroke',
    light: 'fa-sun',
    dark: 'fa-moon'
};

const LABEL: Record<ThemeMode, string> = {
    auto: 'Theme: system',
    light: 'Theme: light',
    dark: 'Theme: dark'
};

// ThemeToggle is a single control pinned to the bottom of the navy icon rail
// (it mirrors the rail's icon-button rhythm). Clicking it cycles
// system -> light -> dark and persists the choice. The icon reflects the chosen
// MODE; auto shows a half-circle, and we add a data attribute carrying the
// RESOLVED theme so the rail's active-state styling can hint the live theme.
export function ThemeToggle({mode, resolved, onChange}: ThemeToggleProps) {
    const next = ORDER[(ORDER.indexOf(mode) + 1) % ORDER.length];
    const title = mode === 'auto' ? `${LABEL[mode]} (${resolved})` : LABEL[mode];
    return (
        <button
            type='button'
            className={classNames('theme-toggle', `theme-toggle--${mode}`)}
            onClick={() => onChange(next)}
            data-resolved={resolved}
            aria-label={title}
            title={title}>
            <i className={classNames('fas', ICON[mode])} aria-hidden='true' />
        </button>
    );
}
