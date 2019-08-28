// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { NgModule } from '@angular/core';
import { IonicModule } from 'ionic-angular';
import { DataTableComponent } from './data-table';

@NgModule({
  declarations: [
    DataTableComponent,
  ],
  imports: [
    IonicModule,
  ],
  exports: [
    DataTableComponent
  ]
})
export class DataTableComponentModule {}
