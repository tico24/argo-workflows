import * as React from 'react';

import * as nsUtils from '../namespaces';
import {InputFilter} from './input-filter';

export const NamespaceFilter = (props: {value: string; onChange: (namespace: string) => void; extraNamespaces?: string[]}) =>
    nsUtils.getManagedNamespace() ? (
        <span className='namespace-filter__value'>{nsUtils.getManagedNamespace()}</span>
    ) : (
        <InputFilter value={props.value} name='ns' onChange={ns => props.onChange(ns)} extraSuggestions={props.extraNamespaces} filterSuggestions />
    );
