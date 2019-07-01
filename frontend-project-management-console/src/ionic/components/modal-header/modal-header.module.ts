// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

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
