import { Component } from '@angular/core';
import { NavController, ModalController, NavParams, IonicPage } from 'ionic-angular';
import { ClaService } from 'cla-service';
import { ClaCompanyModel } from '../../models/cla-company';
import { ClaUserModel } from '../../models/cla-user';
import { ClaSignatureModel } from '../../models/cla-signature';
import { SortService } from '../../services/sort.service';

@IonicPage({
  segment: 'company/:companyId/project/:projectId'
})
@Component({
  selector: 'project-page',
  templateUrl: 'project-page.html',
})
export class ProjectPage {
  signatures: ClaSignatureModel[];
  loading: any;
  companyId: string;
  projectId: string;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    private claService: ClaService,
    public modalCtrl: ModalController,
  ) {
    this.companyId = navParams.get('companyId');
    this.projectId = navParams.get('projectId');
    console.log(this.companyId);
    console.log(this.projectId);
    this.getDefaults();
  }

  getDefaults() {
    this.loading = {};
  }

  ngOnInit() {
    this.getProject();
    this.getProjectSignatures();
  }

  getProject() {
    this.claService.getProject(this.projectId).subscribe(response => {
      console.log('project: ' + this.projectId);
      console.log(response);
    });
  }
  // TODO: need to get all signatures for a project that are employees of a particular company
  getProjectSignatures() {
    this.claService.getProjectSignatures(this.projectId).subscribe(response => {
      console.log('signatures:');
      console.log(response);
      this.signatures = response;
    });
  }

}
