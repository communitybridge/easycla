// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component } from '@angular/core';
import {
  Events,
  IonicPage,
  ModalController,
  NavController,
  NavParams,
  PopoverController,
  ViewController
} from 'ionic-angular';
import { ClaService } from '../../services/cla.service';

@IonicPage({
  segment: 'cla-contract-version-modal'
})
@Component({
  selector: 'cla-contract-version-modal',
  templateUrl: 'cla-contract-version-modal.html'
})
export class ClaContractVersionModal {
  claProjectId: string;
  documentType: string; // individual | corporate
  documents: any;
  currentDocument: any;
  previousDocuments: any;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    public viewCtrl: ViewController,
    private popoverCtrl: PopoverController,
    public modalCtrl: ModalController,
    private claService: ClaService,
    public events: Events
  ) {
    this.claProjectId = this.navParams.get('claProjectId');
    this.documentType = this.navParams.get('documentType');
    this.documents = this.navParams.get('documents').reverse();
    if (this.documents.length > 0) {
      this.currentDocument = this.documents.slice(0, 1);
      console.log(this.documentType + ' document:');
      console.log(this.currentDocument);
      if (this.documents.length > 1) {
        this.previousDocuments = this.documents.slice(1);
      }
    } else {
      console.log('No documents for project: ' + this.claProjectId);
    }

    events.subscribe('modal:close', () => {
      this.dismiss();
    });

    this.getDefaults();
  }

  getDefaults() {}

  ngOnInit() {}

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

  openNewWindow(event, document) {
    if (document.document_s3_url == null || document.document_s3_url == '') {
      const msg =
        'Unable to open PDF document for ' +
        document.document_name +
        ' - document URL is empty.' +
        ' For further assistance with EasyCLA, please ' +
        '<a href="https://jira.linuxfoundation.org/servicedesk/customer/portal/4">submit a support request ticket</a>.';
      console.log(msg);
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
      console.log('Opening new window to' + document.document_s3_url);
      let win = window.open(document.document_s3_url, '_blank', 'titlebar=yes');
    }
  }

  openClaContractUploadModal() {
    let modal = this.modalCtrl.create('ClaContractUploadModal', {
      claProjectId: this.claProjectId,
      documentType: this.documentType
    });
    modal.present();
  }

  dismiss() {
    this.viewCtrl.dismiss();
  }
}
