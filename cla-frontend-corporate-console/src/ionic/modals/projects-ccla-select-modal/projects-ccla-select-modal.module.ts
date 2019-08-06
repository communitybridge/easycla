// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { ProjectsCclaSelectModal } from './projects-ccla-select-modal';
import { LoadingSpinnerComponentModule } from '../../components/loading-spinner/loading-spinner.module'

@NgModule({
  declarations: [
    ProjectsCclaSelectModal,
  ],
  imports: [
    IonicPageModule.forChild(ProjectsCclaSelectModal),
    LoadingSpinnerComponentModule
  ],
  entryComponents: [
    ProjectsCclaSelectModal,
  ]
})
export class ProjectsCclaSelectModalModule {}
