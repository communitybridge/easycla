// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: AGPL-3.0-or-later

import { NgModule } from '@angular/core';
import { IonicModule } from 'ionic-angular';
import { ProjectNavigationComponent } from './project-navigation';

@NgModule({
  declarations: [
    ProjectNavigationComponent,
  ],
  imports: [
    IonicModule,
  ],
  exports: [
    ProjectNavigationComponent,
  ]
})
export class ProjectNavigationComponentModule {}
