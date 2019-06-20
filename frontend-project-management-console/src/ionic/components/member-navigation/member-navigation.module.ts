// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: AGPL-3.0-or-later

import { NgModule } from '@angular/core';
import { IonicModule } from 'ionic-angular';
import { MemberNavigationComponent } from './member-navigation';

@NgModule({
  declarations: [
    MemberNavigationComponent,
  ],
  imports: [
    IonicModule,
  ],
  exports: [
    MemberNavigationComponent,
  ]
})
export class MemberNavigationComponentModule {}
