import { Component, ElementRef, ViewChild, } from '@angular/core';
import { NavController, NavParams, ViewController, ModalController, IonicPage } from 'ionic-angular';
import { EnvConfig } from '../../services/cla.env.utils';

@IonicPage({
  segment: 'cla/project/:projectId/user/:userId/admin-yesno'

})
@Component({
  selector: 'cla-company-admin-yesno-modal',
  templateUrl: 'cla-company-admin-yesno-modal.html',
})
export class ClaCompanyAdminYesnoModal {
  projectId: string;
  repositoryId: string;
  userId: string;

  consoleLink: string;

  @ViewChild('textArea') textArea: ElementRef;


  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    public viewCtrl: ViewController,
    public modalCtrl: ModalController,
  ) {
    this.projectId = navParams.get('projectId');
    this.userId = navParams.get('userId');
    this.getDefaults();
  }

  getDefaults() {
    this.consoleLink = EnvConfig['corp-console-link'];
  }

  ngOnInit() {

  }
  
  dismiss() {
    this.viewCtrl.dismiss();
  }

  openCompanyAdminConsoleLink() {
    window.open(this.consoleLink, '_blank');
  }

  openCompanyAdminSendEmail() {
    let modal = this.modalCtrl.create('ClaCompanyAdminSendEmailModal', {
      projectId: this.projectId,
      userId: this.userId,
    });
    modal.present();
    this.dismiss();
  }

}
