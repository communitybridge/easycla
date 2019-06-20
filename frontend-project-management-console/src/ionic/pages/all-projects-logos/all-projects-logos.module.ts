// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: AGPL-3.0-or-later

import { NgModule } from '@angular/core';

import { IonicPageModule } from 'ionic-angular';

import { AllProjectsLogosPage } from './all-projects-logos';
import { SortingDisplayComponentModule } from '../../components/sorting-display/sorting-display.module';

@NgModule({
  declarations: [
    AllProjectsLogosPage
  ],
  imports: [
    SortingDisplayComponentModule,
    IonicPageModule.forChild(AllProjectsLogosPage)
  ],
})
export class AllProjectsLogosPageModule {}
