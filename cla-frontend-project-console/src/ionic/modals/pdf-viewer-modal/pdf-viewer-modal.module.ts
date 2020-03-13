// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { LoadingSpinnerComponentModule } from '../../components/loading-spinner/loading-spinner.module';
import { LoadingDisplayDirectiveModule } from '../../directives/loading-display/loading-display.module';
import { ModalHeaderComponentModule } from '../../components/modal-header/modal-header.module';
import { PdfViewerComponentModule } from '../../components/pdf-viewer/pdf-viewer.module';
import { PdfViewerModal } from './pdf-viewer-modal';

@NgModule({
  declarations: [PdfViewerModal],
  imports: [
    LoadingSpinnerComponentModule,
    LoadingDisplayDirectiveModule,
    ModalHeaderComponentModule,
    PdfViewerComponentModule,
    IonicPageModule.forChild(PdfViewerModal)
  ],
  entryComponents: [PdfViewerModal]
})
export class ClaContractVersionModalModule {}
