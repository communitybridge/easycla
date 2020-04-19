// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { ProjectClaPage } from './project-cla';
import { LoadingSpinnerComponentModule } from '../../../components/loading-spinner/loading-spinner.module';
import { LoadingDisplayDirectiveModule } from '../../../directives/loading-display/loading-display.module';
import { SortingDisplayComponentModule } from '../../../components/sorting-display/sorting-display.module';
import { SectionHeaderComponentModule } from '../../../components/section-header/section-header.module';
import { ProjectNavigationComponentModule } from '../../../components/project-navigation/project-navigation.module';
import { LayoutModule } from '../../../layout/layout.module';
import { GetHelpComponentModule } from '../../../components/get-help/get-help.module';

@NgModule({
  declarations: [ProjectClaPage],
  imports: [
    GetHelpComponentModule,
    LoadingSpinnerComponentModule,
    LoadingDisplayDirectiveModule,
    SortingDisplayComponentModule,
    SectionHeaderComponentModule,
    ProjectNavigationComponentModule,
    IonicPageModule.forChild(ProjectClaPage),
    LayoutModule
  ],
  entryComponents: [ProjectClaPage]
})
export class ProjectClaPageModule {}
