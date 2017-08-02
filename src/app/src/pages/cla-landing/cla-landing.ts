import { Component } from '@angular/core';
import { NavController, IonicPage, ModalController, NavParams, } from 'ionic-angular';

@IonicPage({
  segment: 'cla/project/:projectId/repository/:repositoryId/user/:userId'
})
@Component({
  selector: 'cla-landing',
  templateUrl: 'cla-landing.html'
})
export class ClaLandingPage {
  projectId: string;
  repositoryId: string;
  userId: string;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    private modalCtrl: ModalController,
    // private cincoService: CincoService,
  ) {
    this.getDefaults();
    this.projectId = navParams.get('projectId');
    this.repositoryId = navParams.get('repositoryId');
    this.userId = navParams.get('userId');
  }

  getDefaults() {

  }

  ngOnInit() {
    this.determineAppropriateAgreements();
  }

  determineAppropriateAgreements() {
    // pass along project and user info to determine appropriate agreements that
    // the user could sign (ICLA, CCLA). If this project doesn't use CLAs then
    // the user can be redirected to the ICLA
  }

  openClaIndividualPage() {
    // send to the individual cla page which will give directions and redirect
    this.navCtrl.push('ClaIndividualPage', {
      projectId: this.projectId,
      repositoryId: this.repositoryId,
      userId: this.userId,
    });
  }

  openClaIndividualEmployeeModal() {
    let modal = this.modalCtrl.create('ClaSelectCompanyModal', {
      projectId: this.projectId,
      repositoryId: this.repositoryId,
      userId: this.userId,
    });
    modal.present();
  }

}
