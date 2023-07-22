import * as common from './common';
import config from '../config';

export function getConfig() {
    return common.get(`/config`);
}
