// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component } from "@angular/core";
import {
  Events,
  IonicPage,
  ModalController,
  NavController,
  NavParams,
  ViewController
} from "ionic-angular";
import { PopoverController } from "ionic-angular";
import { ClaService } from "../../services/cla.service";

@IonicPage({
  segment: "cla-contract-version-modal"
})
@Component({
  selector: "cla-contract-version-modal",
  templateUrl: "cla-contract-version-modal.html"
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
    this.claProjectId = this.navParams.get("claProjectId");
    this.documentType = this.navParams.get("documentType");
    this.documents = this.navParams.get("documents").reverse();
    if (this.documents.length > 0) {
      this.currentDocument = this.documents.slice(0, 1);
      console.log("currentDoc");
      console.log(this.currentDocument);
      if (this.documents.length > 1) {
        this.previousDocuments = this.documents.slice(1);
      }
    }

    events.subscribe('modal:close', () => {
      this.dismiss();
    });

    this.getDefaults();
  }

  getDefaults() {}

  ngOnInit() {}

  contractPopover(ev, document) {
    let currentContractActions = {
      items: [
        {
          label: "View PDF",
          callback: "contractView",
          callbackData: {
            document: document
          }
        }
      ]
    };
    let popover = this.popoverCtrl.create(
      "ActionPopoverComponent",
      currentContractActions
    );

    popover.present({
      ev: ev
    });

    popover.onDidDismiss(popoverData => {
      if (popoverData) {
        this.popoverResponse(popoverData);
      }
    });
  }

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

  contractView(data) {
    let majorVersion = data.document.document_major_version;
    let minorVersion = data.document.document_minor_version;
    this.claService.getProjectDocumentRevisionPdf(
      this.claProjectId,
      this.documentType,
      majorVersion,
      minorVersion
    ).subscribe(base64PDF => {
      let objbuilder = '';
        objbuilder += ('<object width="100%" height="100%" data="data:application/pdf;base64,');
        objbuilder += (base64PDF);
        objbuilder += ('" type="application/pdf" class="internal">');
        objbuilder += ('</object>');

      let win = window.open("","_blank","titlebar=yes");
      win.document.title = data.document_name;
      win.document.write('<html><body>');
      win.document.write(objbuilder);
      win.document.write('</body></html>');
    });
  }

  openClaContractUploadModal() {
    let modal = this.modalCtrl.create("ClaContractUploadModal", {
      claProjectId: this.claProjectId,
      documentType: this.documentType
    });
    modal.present();
  }

  dismiss() {
    this.viewCtrl.dismiss();
  }
}
