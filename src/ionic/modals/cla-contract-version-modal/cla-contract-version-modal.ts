import { Component } from '@angular/core';
import { NavController, ModalController, NavParams, ViewController, IonicPage, } from 'ionic-angular';
import { CincoService } from '../../services/cinco.service';
import { PopoverController } from 'ionic-angular';

@IonicPage({
  segment: 'cla-contract-version-modal'
})
@Component({
  selector: 'cla-contract-version-modal',
  templateUrl: 'cla-contract-version-modal.html',
  providers: [CincoService]
})
export class ClaContractVersionModal {

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    public viewCtrl: ViewController,
    private popoverCtrl: PopoverController,
    public modalCtrl: ModalController,
  ) {
    this.getDefaults();
  }

  ngOnInit() {

  }

  getDefaults() {

  }

  currentContractPopover(ev, index) {
    let currentContractActions = {
      items: [
        {
          label: 'Download',
          callback: 'contractDownload',
          callbackData: {
            index: index,
          }
        },
        {
          label: 'Delete',
          callback: 'contractDelete',
          callbackData: {
            index: index,
          }
        },
      ]
    };
    let popover = this.popoverCtrl.create(
      'ActionPopoverComponent',
      currentContractActions,
    );

    popover.present({
      ev: ev
    });

    popover.onDidDismiss((popoverData) => {
      if(popoverData) {
        this.popoverResponse(popoverData);
      }
    });
  }

  previousContractPopover(ev, index) {
    let previousContractActions = {
      items: [
        {
          label: 'Make Current',
          callback: 'contractMakeCurrent',
          callbackData: {
            index: index,
          }
        },
        {
          label: 'Download',
          callback: 'contractDownload',
          callbackData: {
            index: index,
          }
        },
        {
          label: 'Delete',
          callback: 'contractDelete',
          callbackData: {
            index: index,
          }
        },
      ]
    };

    let popover = this.popoverCtrl.create(
      'ActionPopoverComponent',
      previousContractActions,
    );

    popover.present({
      ev: ev
    });

    popover.onDidDismiss((popoverData) => {
      if(popoverData) {
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
    if(this[callback]) {
      this[callback](popoverData.callbackData);
    }
  }

  contractMakeCurrent(data) {
    console.log('contract make current');
  }

  contractDownload(data) {
    console.log('contract download');
  }

  contractDelete(data) {
    console.log('contract delete');
  }

  openClaContractUploadModal(uploadInfo) {
    let modal = this.modalCtrl.create('ClaContractUploadModal', {
      uploadInfo: uploadInfo,
    });
    modal.present();
  }

  dismiss() {
    this.viewCtrl.dismiss();
  }

}
