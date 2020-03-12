// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { NgModule } from '@angular/core';
import { IonicModule } from 'ionic-angular';
import { PdfViewerComponent } from './pdf-viewer';
import { PdfViewerModule } from 'ng2-pdf-viewer';

@NgModule({
  declarations: [PdfViewerComponent],
  imports: [IonicModule, 
    PdfViewerModule,],
  exports: [PdfViewerComponent]
})
export class PdfViewerComponentModule {}
