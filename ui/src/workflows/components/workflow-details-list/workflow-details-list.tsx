import * as React from 'react';

import {TimestampSwitch} from '../../../shared/components/timestamp';
import * as models from '../../../shared/models';
import useTimestamp, {TIMESTAMP_KEYS} from '../../../shared/use-timestamp';
import {WorkflowsRow} from '../../../workflows/components/workflows-row/workflows-row';

import './workflow-details-list.scss';

interface WorkflowDetailsList {
    workflows: models.Workflow[];
    columns: models.Column[];
}

export function WorkflowDetailsList(props: WorkflowDetailsList) {
    const [storedDisplayISOFormatStart, setStoredDisplayISOFormatStart] = useTimestamp(TIMESTAMP_KEYS.CRON_ROW_STARTED);
    const [storedDisplayISOFormatFinished, setStoredDisplayISOFormatFinished] = useTimestamp(TIMESTAMP_KEYS.CRON_ROW_FINISHED);
    return (
        <div className='argo-table-list workflows-details-list'>
            <div className='row argo-table-list__head'>
                <div className='columns small-1 workflows-list__status' />
                <div className='row small-11'>
                    <div className='columns small-2'>Name</div>
                    <div className='columns small-1'>Namespace</div>
                    <div className='columns small-1'>
                        Started <TimestampSwitch storedDisplayISOFormat={storedDisplayISOFormatStart} setStoredDisplayISOFormat={setStoredDisplayISOFormatStart} />
                    </div>
                    <div className='columns small-1'>
                        Finished <TimestampSwitch storedDisplayISOFormat={storedDisplayISOFormatFinished} setStoredDisplayISOFormat={setStoredDisplayISOFormatFinished} />
                    </div>
                    <div className='columns small-1'>Duration</div>
                    <div className='columns small-1'>Progress</div>
                    <div className='columns small-2'>Message</div>
                    <div className='columns small-1'>Details</div>
                    <div className='columns small-1'>Archived</div>
                    {(props.columns || []).map(col => {
                        return (
                            <div className='columns small-1' key={col.key}>
                                {col.name}
                            </div>
                        );
                    })}
                </div>
            </div>
            {/* checkboxes are not visible and are unused in details pages */}
            {props.workflows.map(wf => {
                return (
                    <WorkflowsRow
                        workflow={wf}
                        key={wf.metadata.uid}
                        checked={false}
                        columns={props.columns}
                        onChange={null}
                        select={null}
                        displayISOFormatStart={storedDisplayISOFormatStart}
                        displayISOFormatFinished={storedDisplayISOFormatFinished}
                    />
                );
            })}
        </div>
    );
}
