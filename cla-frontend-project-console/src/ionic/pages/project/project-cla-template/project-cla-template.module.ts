// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { ProjectClaTemplatePage } from './project-cla-template';
import { LoadingSpinnerComponentModule } from '../../../components/loading-spinner/loading-spinner.module';
import { LoadingDisplayDirectiveModule } from '../../../directives/loading-display/loading-display.module';
import { SortingDisplayComponentModule } from '../../../components/sorting-display/sorting-display.module';
import { SectionHeaderComponentModule } from '../../../components/section-header/section-header.module';
import { ProjectNavigationComponentModule } from '../../../components/project-navigation/project-navigation.module';
import { LayoutModule } from '../../../layout/layout.module';
import { PdfViewerComponentModule } from '../../../components/pdf-viewer/pdf-viewer.module';
import { GetHelpComponentModule } from '../../../components/get-help/get-help.module';

@NgModule({
  declarations: [ProjectClaTemplatePage],
  imports: [
    GetHelpComponentModule,
    LoadingSpinnerComponentModule,
    LoadingDisplayDirectiveModule,
    SortingDisplayComponentModule,
    SectionHeaderComponentModule,
    ProjectNavigationComponentModule,
    PdfViewerComponentModule,
    IonicPageModule.forChild(ProjectClaTemplatePage),
    LayoutModule
  ],
  entryComponents: [ProjectClaTemplatePage]
})
export class ProjectClaTemplatePageModule {}
