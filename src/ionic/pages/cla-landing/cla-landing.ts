import { Component } from '@angular/core';
import { NavController, IonicPage, ModalController, NavParams, } from 'ionic-angular';
import { ClaService } from 'cla-service';

@IonicPage({
  segment: 'cla/project/:projectId/user/:userId'
})
@Component({
  selector: 'cla-landing',
  templateUrl: 'cla-landing.html'
})
export class ClaLandingPage {
  projectId: string;
  userId: string;

  user: any;
  project: any;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    private modalCtrl: ModalController,
    private claService: ClaService,
  ) {
    this.getDefaults();
    this.projectId = navParams.get('projectId');
    this.userId = navParams.get('userId');
  }

  getDefaults() {
    this.project = {
      project_name: "",
    }

  }

  ngOnInit() {
    this.determineAppropriateAgreements();
    this.getUser(this.userId);
    this.getProject(this.projectId);
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
      userId: this.userId,
    });
  }

  openClaIndividualEmployeeModal() {
    let modal = this.modalCtrl.create('ClaSelectCompanyModal', {
      projectId: this.projectId,
      userId: this.userId,
    });
    modal.present();
  }

  getUser(userId) {
    this.claService.getUser(userId).subscribe(response => {
      this.user = response;
    });
  }

  getProject(projectId) {
    this.claService.getProject(projectId).subscribe(response => {
      this.project = response;
      if (!this.project.logoRef) {
        this.project.logoRef = "https://dummyimage.com/200x100/bbb/fff.png&text=+";
      }
    });
  }

}
