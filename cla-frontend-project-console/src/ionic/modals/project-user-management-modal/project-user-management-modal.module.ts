// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { ProjectUserManagementModal } from './project-user-management-modal';

@NgModule({
  declarations: [ProjectUserManagementModal],
  imports: [IonicPageModule.forChild(ProjectUserManagementModal)],
  entryComponents: [ProjectUserManagementModal]
})
export class ProjectUserManagementModalModule {}
