import {fireEvent, render, screen} from '@testing-library/react';
import React from 'react';

import {FiltersCollapseToggle} from './filters-collapse-toggle';

describe('FiltersCollapseToggle', () => {
    it('renders the expanded (collapse) affordance and fires onToggle', () => {
        const onToggle = jest.fn();
        render(<FiltersCollapseToggle collapsed={false} onToggle={onToggle} />);
        const button = screen.getByRole('button', {name: 'Collapse filters'});
        expect(button).toBeInTheDocument();
        expect(button).toHaveAttribute('aria-expanded', 'true');
        fireEvent.click(button);
        expect(onToggle).toHaveBeenCalledTimes(1);
    });

    it('renders the collapsed (show) handle with a label and fires onToggle', () => {
        const onToggle = jest.fn();
        render(<FiltersCollapseToggle collapsed={true} onToggle={onToggle} />);
        const button = screen.getByRole('button', {name: 'Show filters'});
        expect(button).toBeInTheDocument();
        expect(button).toHaveAttribute('aria-expanded', 'false');
        expect(screen.getByText('Filters')).toBeInTheDocument();
        fireEvent.click(button);
        expect(onToggle).toHaveBeenCalledTimes(1);
    });
});
