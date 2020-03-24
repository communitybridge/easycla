// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { NgModule } from '@angular/core';
import { IonicModule } from 'ionic-angular';
import { SortTableComponent } from './sort-table';

@NgModule({
  declarations: [SortTableComponent],
  imports: [IonicModule],
  exports: [SortTableComponent],
  entryComponents: [SortTableComponent]
})
export class SortTableComponentModule {}
