// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { ClaContractViewSignaturesModal } from './cla-contract-view-signatures-modal';
import { LoadingSpinnerComponentModule } from '../../components/loading-spinner/loading-spinner.module';
import { LoadingDisplayDirectiveModule } from '../../directives/loading-display/loading-display.module';
import { ModalHeaderComponentModule } from "../../components/modal-header/modal-header.module";
import {SortingDisplayComponentModule} from "../../components/sorting-display/sorting-display.module";
import {NgxPaginationModule} from 'ngx-pagination';

@NgModule({
  declarations: [
    ClaContractViewSignaturesModal
  ],
  imports: [
    LoadingSpinnerComponentModule,
    LoadingDisplayDirectiveModule,
    ModalHeaderComponentModule,
    SortingDisplayComponentModule,
    NgxPaginationModule,
    IonicPageModule.forChild(ClaContractViewSignaturesModal)
  ],
  entryComponents: [
    ClaContractViewSignaturesModal
  ]
})
export class ClaContractViewSignaturesModalModule {}
