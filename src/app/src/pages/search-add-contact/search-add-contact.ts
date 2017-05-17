import { Component, ChangeDetectorRef } from '@angular/core';
import { NavController, NavParams, ModalController, ViewController, AlertController, IonicPage } from 'ionic-angular';
import { CincoService } from '../../app/services/cinco.service'

@IonicPage({
  segment: 'search-add-contact'
})
@Component({
  selector: 'search-add-contact',
  templateUrl: 'search-add-contact.html',
  providers: [CincoService]
})
export class SearchAddContact {
  projectId: string;
  memberId: string;
  org: any;
  enteredEmail: string;
  organizationContacts: any;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    public modalCtrl: ModalController,
    public viewCtrl: ViewController,
    public alertCtrl: AlertController,
    private changeDetectorRef: ChangeDetectorRef,
    private cincoService: CincoService
  ) {
    this.getDefaults();
    this.projectId = this.navParams.get('projectId');
    this.memberId = this.navParams.get('memberId');
    this.org = this.navParams.get('org');
  }

  ngOnInit() {
    let orgId = this.org.id;
    this.getOrganizationContacts(orgId);
  }

  getDefaults() {

  }

  getOrganizationContacts(orgId) {
    this.cincoService.getOrganizationContacts(orgId).subscribe(response => {
      if(response) {
        this.organizationContacts = response;
      }
    });
  }

  // ContactUpdate modal dismiss
  dismiss() {
    this.viewCtrl.dismiss();
  }

  addContact(contact) {
    let modal = this.modalCtrl.create('ContactUpdate', {
      projectId: this.projectId,
      memberId: this.memberId,
      org: this.org,
      contact: contact,
    });
    modal.onDidDismiss(data => {
      // A refresh of data anytime the modal is dismissed
      let orgId = this.org.id;
      this.getOrganizationContacts(orgId);
    });
    modal.present();
  }

  filterContactsByEmail() {

  }

}
