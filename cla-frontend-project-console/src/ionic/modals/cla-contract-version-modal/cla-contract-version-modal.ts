// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component } from '@angular/core';
import {
  Events,
  IonicPage,
  ModalController,
  NavController,
  NavParams,
  ViewController
} from 'ionic-angular';

import { PlatformLocation } from '@angular/common';

@IonicPage({
  segment: 'cla-contract-version-modal'
})
@Component({
  selector: 'cla-contract-version-modal',
  templateUrl: 'cla-contract-version-modal.html'
})
export class ClaContractVersionModal {
  documentType: string; // individual | corporate
  documents: any;
  currentDocument: any;
  previousDocuments: any;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    public viewCtrl: ViewController,
    public modalCtrl: ModalController,
    public events: Events,
    private location: PlatformLocation,
  ) {
    this.location.onPopState(() => {
      this.viewCtrl.dismiss(false);
    });
    this.documentType = this.navParams.get('documentType');
    this.documents = this.navParams.get('documents').reverse();
    if (this.documents.length > 0) {
      this.currentDocument = this.documents.slice(0, 1);
      if (this.documents.length > 1) {
        this.previousDocuments = this.documents.slice(1);
      }
    }

    events.subscribe('modal:close', () => {
      this.dismiss();
    });
  }

  ngOnInit() { }

  /**
   * Called if popover dismissed with data. Passes data to a callback function
   * @param  {object} popoverData should contain .callback and .callbackData
   */
  popoverResponse(popoverData) {
    let callback = popoverData.callback;
    if (this[callback]) {
      this[callback](popoverData.callbackData);
    }
  }

  async openPdfModal(document) {
    if (document.documentS3URL == null || document.documentS3URL == '') {
      const msg =
        'Unable to open PDF document for ' +
        document.documentName +
        ' - document URL is empty.' +
        ' For further assistance with EasyCLA, please ' +
        '<a href="https://jira.linuxfoundation.org/servicedesk/customer/portal/4">submit a support request ticket</a>.';
      const height = Math.floor(window.innerHeight * 0.8);
      const width = Math.floor(window.innerWidth * 0.8);
      const myWindow = window.open('', '_blank', 'titlebar=yes,width=' + width + ',height=' + height);
      myWindow.document.write(
        '<html><head><title>PDF Document Unavailable</title></head>' +
        '<body>' +
        '<p>' +
        msg +
        '</p>' +
        '</body></html>'
      );
    } else {
      const modal = this.modalCtrl.create('PdfViewerModal', { doc: document.documentS3URL, documentType: this.documentType });
      modal.present();
    }
  }

  dismiss() {
    this.viewCtrl.dismiss();
  }
}
