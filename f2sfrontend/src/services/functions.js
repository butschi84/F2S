import * as common from './common';
import config from '../config';

export function getFunctions() {
    return common.get(`/api/products`);
}
