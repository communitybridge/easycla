// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { FormControl } from '@angular/forms';

export class EmailValidator {
  static isValid(control: FormControl): any {
    const EMAIL_PATTERN = new RegExp(['[a-z|A-Z|0-9]+[@]+[a-z|A-Z|0-9]+[.]+([a-z|A-Z|0-9]){2}'].join(''));
    let email = control.value;
    let isValid = EMAIL_PATTERN.test(email);
    if (!isValid) {
      return {
        'not a valid email address': true
      };
    }

    return null;
  }
}
