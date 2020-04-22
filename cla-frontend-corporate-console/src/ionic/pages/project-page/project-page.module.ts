// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { ProjectPage } from './project-page';
import { LoadingSpinnerComponentModule } from '../../components/loading-spinner/loading-spinner.module';
import { LoadingDisplayDirectiveModule } from '../../directives/loading-display/loading-display.module';
import { SortingDisplayComponentModule } from '../../components/sorting-display/sorting-display.module';
import { LayoutModule } from '../../layout/layout.module';
import { GetHelpComponentModule } from '../../components/get-help/get-help.module';
import { SharedModule } from '../../shared.module';

@NgModule({
  declarations: [ProjectPage],
  imports: [
    LoadingSpinnerComponentModule,
    LoadingDisplayDirectiveModule,
    SortingDisplayComponentModule,
    SharedModule,
    IonicPageModule.forChild(ProjectPage),
    LayoutModule,
    GetHelpComponentModule
  ],
  entryComponents: [ProjectPage]
})
export class ProjectPageModule { }
