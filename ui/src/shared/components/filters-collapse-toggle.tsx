import classNames from 'classnames';
import * as React from 'react';

interface FiltersCollapseToggleProps {
    collapsed: boolean;
    onToggle: () => void;
}

// FiltersCollapseToggle renders the affordance that collapses / expands the
// (right-hand) filter sidebar. It is shared by all four list views
// (workflows / cron / templates / reports) so the behaviour and the markup
// stay identical and a single test covers them all.
//
// One NEUTRAL control system, two states (no teal fill — see _modern-filters-layout.scss):
//  - expanded: a compact, anchored ghost icon button capping the top-right of
//    the filter panel. A "collapse panel to the right" double-chevron (»)
//    reads as "push the panel away to the edge".
//  - collapsed: a tasteful neutral panel-edge handle tucked against the right
//    edge of the now full-width table — a filter glyph over a quiet vertical
//    "Filters" label, with a left-pointing double-chevron hint, so it reads as
//    a deliberate "pull the filters back" tab rather than leftover chrome.
export function FiltersCollapseToggle({collapsed, onToggle}: FiltersCollapseToggleProps) {
    return (
        <button
            type='button'
            className={classNames('filters-collapse-toggle', {'filters-collapse-toggle--collapsed': collapsed})}
            onClick={onToggle}
            aria-expanded={!collapsed}
            aria-label={collapsed ? 'Show filters' : 'Collapse filters'}
            title={collapsed ? 'Show filters' : 'Collapse filters'}>
            {collapsed ? (
                <>
                    <i className='fa fa-filter filters-collapse-toggle__glyph' aria-hidden='true' />
                    <span className='filters-collapse-toggle__label'>Filters</span>
                    <i className='fa fa-angle-double-left filters-collapse-toggle__chevron' aria-hidden='true' />
                </>
            ) : (
                <i className='fa fa-angle-double-right' aria-hidden='true' />
            )}
        </button>
    );
}
