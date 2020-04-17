// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Input, Component } from '@angular/core';

@Component({
  selector: 'sorting-display',
  templateUrl: 'sorting-display.html'
})
export class SortingDisplayComponent {
  @Input('sorting')
  private sorting: string;

  constructor() {
    this.sorting;
  }
}
