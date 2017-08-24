import { Component } from '@angular/core';
import { NavController, NavParams, ViewController, IonicPage, ModalController } from 'ionic-angular';
import { CincoService } from '../../services/cinco.service'

@IonicPage({
  segment: 'cla-corporate-member-modal'
})
@Component({
  selector: 'cla-corporate-member-modal',
  templateUrl: 'cla-corporate-member-modal.html',
  providers: [CincoService]
})
export class ClaCorporateMemberModal {
  members: any;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    public viewCtrl: ViewController,
    private cincoService: CincoService,
    public modalCtrl: ModalController,
  ) {
    this.getDefaults();
  }

  ngOnInit() {

  }

  getDefaults() {
    this.members = [
      {
        id: 'A000000001',
        name: 'Google',
      },
      {
        id: 'A000000002',
        name: 'Member 2',
      },
      {
        id: 'A000000003',
        name: 'Member 3',
      },
      {
        id: 'A000000004',
        name: 'Member 4',
      },
    ];
  }

  openClaCorporateMemberOptionsModal() {
    let modal = this.modalCtrl.create('ClaCorporateMemberOptionsModal');
    modal.present();
  }

  // ContactUpdateModal modal dismiss
  dismiss() {
    this.viewCtrl.dismiss();
  }



}
