import {Checkbox} from 'argo-ui/src/components/checkbox';
import * as React from 'react';

import './checkbox-filter.scss';

interface Props {
    items: {name: string; count: number}[];
    type: string;
    selected: string[];
    onChange: (selected: string[]) => void;
    // Display-only: when true, render the right-aligned per-item count (Argo CD
    // status-filter style). Off by default so callers whose `count` is not a
    // meaningful tally (e.g. the cron State filter's placeholder values) do not
    // surface a misleading number. Does not affect filtering behaviour.
    showCounts?: boolean;
}

export function CheckboxFilter(props: Props) {
    const unavailableSelected = props.selected.filter(selected => !props.items.some(item => item.name === selected));
    const items = props.items.concat(unavailableSelected.map(selected => ({name: selected, count: 0})));

    return (
        <ul className='checkbox-filter columns small-12'>
            {items.map(item => (
                <li key={item.name}>
                    <React.Fragment>
                        <div className='row'>
                            <div className='checkbox-filter__label columns small-12'>
                                <Checkbox
                                    checked={props.selected.indexOf(item.name) > -1}
                                    id={`filter-${props.type}-${item.name}`}
                                    onChange={() => {
                                        const newSelected = props.selected.slice();
                                        const index = newSelected.indexOf(item.name);
                                        if (index > -1) {
                                            newSelected.splice(index, 1);
                                        } else {
                                            newSelected.push(item.name);
                                        }
                                        props.onChange(newSelected);
                                    }}
                                />{' '}
                                <span className={`checkbox-filter__dot checkbox-filter__dot--${props.type} checkbox-filter__dot--${props.type}-${item.name}`} aria-hidden='true' />
                                <label title={item.name} htmlFor={`filter-${props.type}-${item.name}`}>
                                    {item.name}
                                </label>
                                {props.showCounts && <span className='checkbox-filter__count'>{item.count}</span>}
                            </div>
                        </div>
                    </React.Fragment>
                </li>
            ))}
        </ul>
    );
}
