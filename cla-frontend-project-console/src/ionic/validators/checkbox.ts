// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { FormControl } from '@angular/forms';

export class CheckboxValidator {
  static isChecked(control: FormControl): any {
    if (control.value != true) {
      return {
        notChecked: true
      };
    }

    return null;
  }
}
