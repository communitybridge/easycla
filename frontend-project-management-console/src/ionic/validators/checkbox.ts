// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: AGPL-3.0-or-later

import { FormControl } from '@angular/forms';

export class CheckboxValidator {

  static isChecked(control: FormControl) : any {

    if (control.value != true) {
      return {
        'notChecked' : true
      };
    }

    return null;
  }

}
