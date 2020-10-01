// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { AppSettings } from './app-settings';
import { EnvConfig } from './cla-env-utils';

export const AUTH_ROUTE = '#/auth';


// The URL to which Auth0 will redirect the browser after authorization has been granted for the user.
export function getAuthURLFromWindow(type) {
    const redirectConsole = (type === 'Projects') ? AppSettings.PROJECT_CONSOLE_LINK : AppSettings.CORPORATE_CONSOLE_LINK;
    console.log(EnvConfig.default[redirectConsole]);
    return EnvConfig.default[redirectConsole] + AUTH_ROUTE;
}

