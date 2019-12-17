// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { FormControl } from '@angular/forms';

export class EmailValidator {
  static isValid(control: FormControl): any {
    let email = control.value;
    let i = document.createElement('input');
    i.type = 'email';
    i.value = email;
    let mismatch = i.validity.typeMismatch;

    if (mismatch) {
      return {
        'not a valid email address': true
      };
    }

    return null;
  }
}
