// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: AGPL-3.0-or-later

import { NgModule } from '@angular/core';
import { IonicModule } from 'ionic-angular';
import { LoadingSpinnerComponent } from './loading-spinner';

@NgModule({
  declarations: [
    LoadingSpinnerComponent,
  ],
  imports: [
    IonicModule,
  ],
  exports: [
    LoadingSpinnerComponent
  ]
})
export class LoadingSpinnerComponentModule {}
