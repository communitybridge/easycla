// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: AGPL-3.0-or-later

import { NgModule } from '@angular/core';
import { IonicModule } from 'ionic-angular';
import { SortingDisplayComponent } from './sorting-display';

@NgModule({
  declarations: [
    SortingDisplayComponent,
  ],
  imports: [
    IonicModule,
  ],
  exports: [
    SortingDisplayComponent
  ]
})
export class SortingDisplayComponentModule {}
