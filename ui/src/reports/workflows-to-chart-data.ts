import type {AnnotationOptions} from 'chartjs-plugin-annotation';

import {denominator} from '../shared/duration';
import {Workflow} from '../shared/models';
import type {Chart} from './reports';

// Neutral/teal duration bar color. Duration is a neutral metric, so it must NOT
// be colored by workflow phase (which painted every bar the saturated alarm-red
// of getColorForNodePhase for failed/errored runs, reading as an error for a
// plain duration). Reuse the same teal accent the CPU resource series uses
// below (#00a2b3 === --awf-accent), reserving red for actual failure. AA-clean.
const DURATION_BAR_COLOR = '#00a2b3';

export function workflowsToChartData(workflows: Workflow[], limit: number): Chart[] {
    const filteredWorkflows = workflows
        .filter(wf => !!wf.status.finishedAt)
        .map(wf => ({
            name: wf.metadata.name,
            finishedAt: new Date(wf.status.finishedAt),
            startedAt: new Date(wf.status.startedAt),
            phase: wf.status.phase,
            resourcesDuration: wf.status.resourcesDuration
        }))
        .sort((a, b) => b.finishedAt.getTime() - a.finishedAt.getTime())
        .slice(0, limit)
        .reverse();

    const labels: string[] = new Array(filteredWorkflows.length);
    const backgroundColors: string[] = new Array(filteredWorkflows.length);
    const durationData: number[] = new Array(filteredWorkflows.length);
    const resourceData = {} as {[resource: string]: number[]};

    filteredWorkflows.forEach((wf, i) => {
        labels[i] = wf.name;
        backgroundColors[i] = DURATION_BAR_COLOR;
        durationData[i] = (wf.finishedAt.getTime() - wf.startedAt.getTime()) / 1000;
        Object.entries(wf.resourcesDuration || {}).forEach(([resource, value]) => {
            if (!resourceData[resource]) {
                resourceData[resource] = new Array(filteredWorkflows.length);
            }
            resourceData[resource][i] = value;
        });
    });
    const resourceColors = {
        'cpu': '#00a2b3',
        'memory': '#3576c4',
        'storage': '#f4c030',
        'ephemeral-storage': '#806100'
    } as {[resource: string]: string};

    const avgDuration = durationData.length > 0 ? durationData.reduce((a, b) => a + b, 0) / durationData.length : 0;
    const avgLine: AnnotationOptions<'line'> = {
        type: 'line',
        scaleID: 'duration',
        yMin: avgDuration,
        yMax: avgDuration,
        borderColor: 'gray',
        borderWidth: 1,
        label: {
            display: true,
            position: 'start',
            content: 'Average'
        }
    };

    return [
        {
            data: {
                name: 'duration',
                labels,
                datasets: [
                    {
                        yAxisID: 'duration',
                        data: durationData,
                        backgroundColor: backgroundColors
                    }
                ]
            },
            options: {
                scales: {
                    x: {},
                    duration: {
                        beginAtZero: true,
                        ticks: {},
                        title: {
                            display: true,
                            text: 'Duration (seconds)'
                        }
                    }
                },
                plugins: {
                    title: {
                        display: true,
                        text: 'Duration'
                    },
                    legend: {display: false},
                    annotation: {
                        annotations: {
                            avgLine
                        }
                    }
                }
            }
        },
        {
            data: {
                name: 'resources',
                labels,
                datasets: Object.entries(resourceData).map(([resource, data]) => ({
                    yAxisID: resource,
                    label: resource,
                    data,
                    backgroundColor: resourceColors[resource] || '#363c4a'
                }))
            },
            options: {
                scales: {
                    x: {},
                    ...Object.fromEntries(
                        Object.keys(resourceData).map(resource => [
                            resource,
                            {
                                beginAtZero: true,
                                ticks: {},
                                title: {
                                    display: true,
                                    text: resource + ' (' + denominator(resource) + ')'
                                }
                            }
                        ])
                    )
                },
                plugins: {
                    title: {
                        display: true,
                        text: 'Resources (not available for archived workflows)'
                    }
                }
            }
        }
    ];
}
