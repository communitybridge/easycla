// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component } from '@angular/core';
import { IonicPage, NavParams, ViewController, } from 'ionic-angular';
import { PlatformLocation } from '@angular/common';

@IonicPage({
  segment: 'pdf-viewer-modal'
})
@Component({
  selector: 'pdf-viewer-modal',
  templateUrl: 'pdf-viewer-modal.html'
})
export class PdfViewerModal {
  doc: string;
  documentType: string;

  constructor(private navParams: NavParams,
    private mdContorller: ViewController,
    private location: PlatformLocation
  ) {
    this.location.onPopState(() => {
      this.mdContorller.dismiss(false);
    });
    this.doc = this.navParams.get('doc');
    this.documentType = this.navParams.get('documentType');
  }

  closeModal() {
    this.mdContorller.dismiss();
  }

}
