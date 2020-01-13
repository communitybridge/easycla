// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { CompanyPage } from './company-page';
import { LoadingSpinnerComponentModule } from '../../components/loading-spinner/loading-spinner.module';
import { LoadingDisplayDirectiveModule } from '../../directives/loading-display/loading-display.module';
import { SortingDisplayComponentModule } from '../../components/sorting-display/sorting-display.module';
import { LayoutModule } from '../../layout/layout.module';
import { NgxPaginationModule } from 'ngx-pagination';
import { NgxDatatableModule } from '@swimlane/ngx-datatable';
import { GetHelpComponentModule } from '../../components/get-help/get-help.module'; 

@NgModule({
  declarations: [CompanyPage],
  imports: [
    LoadingSpinnerComponentModule,
    LoadingDisplayDirectiveModule,
    SortingDisplayComponentModule,
    LayoutModule,
    NgxPaginationModule,
    NgxDatatableModule,
    IonicPageModule.forChild(CompanyPage),
    NgxDatatableModule,
    NgxDatatableModule,
    GetHelpComponentModule
  ],
  entryComponents: [CompanyPage]
})
export class CompanyPageModule {}
