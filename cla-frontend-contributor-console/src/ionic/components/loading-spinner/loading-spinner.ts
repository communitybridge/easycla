// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Input, Component } from '@angular/core';

@Component({
  selector: 'loading-spinner',
  templateUrl: 'loading-spinner.html'
})
export class LoadingSpinnerComponent {
  /**
   * The text used for the upload label
   */
  @Input('loading')
  private loading: boolean;

  constructor() {
    this.loading = true;
  }
}
