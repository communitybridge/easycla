// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: AGPL-3.0-or-later

import { FormControl } from '@angular/forms';
declare var require: any;
var phoneUtil = require('google-libphonenumber').PhoneNumberUtil.getInstance();

export class PhoneNumberValidator {

    static isValid(control: FormControl): any {

      let number = control.value;
      let phoneProto;
      if(number == null || number == '') {
        return null;
      }

      let countryCode = 'US';
      if(number.charAt(0)==='+') {
        countryCode = 'ZZ';
      }

      try {
        phoneProto = phoneUtil.parse(number, countryCode);
      }
      catch (e) {
        return {
          'not a phone nubmer': true
        }
      }

      let isValid = phoneUtil.isValidNumber(phoneProto);

      if (!isValid) {
        return {
          'not valid': true
        };
      }

      return null;
    }

}
