// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { FormControl } from '@angular/forms';

export class UserNameValidator {
    static isValid(control: FormControl): any {
        const NAME_PATTERN = new RegExp(['^[A-Za-z]+([\ A-Za-z]+)*'].join(''));
        const name = control.value;
        let isValid = NAME_PATTERN.test(name);
        if (!isValid) {
            return {
                'not a valid name': true
            };
        }

        return null;
    }
}
