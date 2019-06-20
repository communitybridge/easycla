// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: AGPL-3.0-or-later

import { NgModule } from '@angular/core';
import { IonicModule } from 'ionic-angular';
import { ModalHeaderComponent } from './modal-header';

@NgModule({
  declarations: [
    ModalHeaderComponent,
  ],
  imports: [
    IonicModule,
  ],
  exports: [
    ModalHeaderComponent,
  ]
})
export class ModalHeaderComponentModule {}
