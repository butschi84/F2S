import * as common from './common';
import config from '../config';

export function getFunctions() {
    return common.get(`/functions`);
}

export function createFunction(f2sfunction) {
    return common.post(`/functions`, f2sfunction);
}