// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Input, Component } from '@angular/core';

@Component({
  selector: 'data-table',
  templateUrl: 'data-table.html'
})
export class DataTableComponent {
  
  @Input('tableData')
  private tableData: any[];

  constructor() {
    this.tableData;
  }

}
