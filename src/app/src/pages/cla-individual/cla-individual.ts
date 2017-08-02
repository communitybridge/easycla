import { Component } from '@angular/core';
import { NavController, NavParams, IonicPage, ModalController, } from 'ionic-angular';

@IonicPage({
  segment: 'cla/project/:projectId/repository/:repositoryId/user/:userId/individual'
})
@Component({
  selector: 'cla-individual',
  templateUrl: 'cla-individual.html'
})
export class ClaIndividualPage {
  projectId: string;
  repositoryId: string;
  userId: string;

  project: any;
  agreementUrl: string;

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
    this.getProject();
    this.getAgreementUrl();
  }

  getProject() {
    this.project = {
      id: '0000000001',
      name: 'Project Name',
      logoRef: 'https://dummyimage.com/225x102/d8d8d8/242424.png&text=Project+Logo',
    };
  }

  getAgreementUrl() {
    console.log('get agreement url');
    console.log(this.projectId);
    console.log(this.repositoryId);
    console.log(this.userId);
    this.agreementUrl = 'https://docusign.com';
    this.openClaAgreement();
  }


  openClaAgreement() {
    console.log('info for docusign');
    window.open(this.agreementUrl, '_blank');
  }

}
