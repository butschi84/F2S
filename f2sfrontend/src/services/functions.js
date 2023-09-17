import * as common from './common';

export function getFunctions() {
    return common.get(`/functions`);
}

export function createFunction(f2sfunction) {
    return common.post(`/functions`, f2sfunction);
}

export function deleteFunction(f2sfunction) {
    return common.del(`/functions/${f2sfunction.uid}`);
}

export function getMetricLastValue(metricQuery) {
    return common.get(`/prometheus/query?query=${metricQuery}`);
}