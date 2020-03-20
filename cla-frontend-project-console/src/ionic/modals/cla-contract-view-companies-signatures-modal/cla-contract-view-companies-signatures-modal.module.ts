// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { ClaContractViewCompaniesSignaturesModal } from './cla-contract-view-companies-signatures-modal';
import { LoadingSpinnerComponentModule } from '../../components/loading-spinner/loading-spinner.module';
import { LoadingDisplayDirectiveModule } from '../../directives/loading-display/loading-display.module';
import { ModalHeaderComponentModule } from '../../components/modal-header/modal-header.module';
import { SortingDisplayComponentModule } from '../../components/sorting-display/sorting-display.module';
import { NgxPaginationModule } from 'ngx-pagination';
import { NgxDatatableModule } from '@swimlane/ngx-datatable';
import { SortTableComponentModule } from '../../components/sort-table/sort-table.module';

@NgModule({
  declarations: [ClaContractViewCompaniesSignaturesModal],
  imports: [
    LoadingSpinnerComponentModule,
    LoadingDisplayDirectiveModule,
    ModalHeaderComponentModule,
    SortingDisplayComponentModule,
    NgxPaginationModule,
    NgxDatatableModule,
    SortTableComponentModule,
    IonicPageModule.forChild(ClaContractViewCompaniesSignaturesModal)
  ],
  entryComponents: [ClaContractViewCompaniesSignaturesModal]
})
export class ClaContractViewCompaniesSignaturesModalModule {}
