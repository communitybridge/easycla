// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: AGPL-3.0-or-later

import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { MembershipsModal } from './memberships-modal';
import { SortingDisplayComponentModule } from '../../components/sorting-display/sorting-display.module';

@NgModule({
  declarations: [
    MembershipsModal,
  ],
  imports: [
    SortingDisplayComponentModule,
    IonicPageModule.forChild(MembershipsModal),
  ],
  entryComponents: [
    MembershipsModal,
  ]
})
export class MembershipsModalModule {}
