// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: AGPL-3.0-or-later

import { NgModule } from '@angular/core';
import { IonicModule } from 'ionic-angular';
import { SectionHeaderComponent } from './section-header';

@NgModule({
  declarations: [
    SectionHeaderComponent,
  ],
  imports: [
    IonicModule,
  ],
  exports: [
    SectionHeaderComponent,
  ]
})
export class SectionHeaderComponentModule {}
