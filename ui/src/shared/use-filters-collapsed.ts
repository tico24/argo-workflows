import {useLocalStorage} from './use-local-storage';

// ONE shared localStorage key so the filter-sidebar collapsed preference
// follows the user across all four list views (workflows / cron / templates /
// reports). Default is `false` (expanded) so the baseline look is preserved on
// first load — zero regression for existing users.
export const FILTERS_COLLAPSED_KEY = 'filtersCollapsed';

export function useFiltersCollapsed(): [boolean, () => void] {
    const [collapsed, setCollapsed] = useLocalStorage<boolean>(FILTERS_COLLAPSED_KEY, false);
    const toggle = () => setCollapsed(!collapsed);
    return [collapsed, toggle];
}
