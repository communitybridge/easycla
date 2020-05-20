// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { FormControl } from '@angular/forms';

export class LFIDValidator {
    static isValid(control: FormControl): any {
        const LFID_PATTERN = new RegExp(['^[a-zA-Z0-9]+([._]?[a-zA-Z0-9]+)*$'].join(''));
        const lfid = control.value;
        let isValid = LFID_PATTERN.test(lfid);
        if (!isValid) {
            return {
                'not a valid LFID': true
            };
        }

        return null;
    }
}
