// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: AGPL-3.0-or-later

import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { AssignUserModal } from './assign-user-modal';
import { LoadingSpinnerComponentModule } from '../../components/loading-spinner/loading-spinner.module';
import { LoadingDisplayDirectiveModule } from '../../directives/loading-display/loading-display.module';

@NgModule({
  declarations: [
    AssignUserModal
  ],
  imports: [
    LoadingSpinnerComponentModule,
    LoadingDisplayDirectiveModule,
    IonicPageModule.forChild(AssignUserModal)
  ],
  entryComponents: [
    AssignUserModal
  ]
})
export class AssignUserModalModule {}
