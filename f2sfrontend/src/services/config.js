import * as common from './common';

export function getConfig() {
    return common.get(`/config`);
}
