// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { FormControl } from '@angular/forms';

export class CalendarLinkValidator {
  static isValid(control: FormControl): any {
    let entered_url = control.value;
    if (entered_url == null || entered_url == '') {
      return null;
    }
    let calendar_url = 'https://calendar.google.com/calendar/embed';
    let calendar_embed = '<iframe src="https://calendar.google.com/calendar/embed';
    let valid_url = entered_url.startsWith(calendar_url) || entered_url.startsWith(calendar_embed);
    if (!valid_url) {
      return {
        'not a valid calendar url': true
      };
    }

    return null;
  }
}
