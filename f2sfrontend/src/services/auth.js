import * as common from './common';
import config from '../config';

export function getAuthType() {
    return common.get(`/auth/type`);
}

export function signin() {
    return common.get(`/auth/signin`);
}