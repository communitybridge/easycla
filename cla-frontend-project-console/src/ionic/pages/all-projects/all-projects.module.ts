// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { AllProjectsPage } from './all-projects';
import { LoadingSpinnerComponentModule } from '../../components/loading-spinner/loading-spinner.module';
import { LoadingDisplayDirectiveModule } from '../../directives/loading-display/loading-display.module';
import { LayoutModule } from '../../layout/layout.module';

@NgModule({
  declarations: [AllProjectsPage],
  imports: [
    LoadingSpinnerComponentModule,
    LoadingDisplayDirectiveModule,
    IonicPageModule.forChild(AllProjectsPage),
    LayoutModule
  ],
  entryComponents: [AllProjectsPage]
})
export class AllProjectsPageModule {}
