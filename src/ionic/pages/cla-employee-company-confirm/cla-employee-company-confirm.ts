import { Component } from '@angular/core';
import { NavController, NavParams, IonicPage, ModalController, } from 'ionic-angular';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { CheckboxValidator } from  '../../validators/checkbox';
import { ClaService } from 'cla-service';

@IonicPage({
  segment: 'project/:projectId/user/:userId/employee/company/:companyId/confirm'
})
@Component({
  selector: 'cla-employee-company-confirm',
  templateUrl: 'cla-employee-company-confirm.html'
})
export class ClaEmployeeCompanyConfirmPage {
  projectId: string;
  repositoryId: string;
  userId: string;
  companyId: string;

  user: any;
  project: any;
  company: any;
  gitService: string;

  form: FormGroup;
  submitAttempt: boolean = false;
  currentlySubmitting: boolean = false;

  constructor(
    public navCtrl: NavController,
    private modalCtrl: ModalController,
    public navParams: NavParams,
    private formBuilder: FormBuilder,
    private claService: ClaService,
  ) {
    this.getDefaults();
    this.projectId = navParams.get('projectId');
    this.repositoryId = navParams.get('repositoryId');
    this.userId = navParams.get('userId');
    this.companyId = navParams.get('companyId');

    this.form = formBuilder.group({
      agree:[false, Validators.compose([CheckboxValidator.isChecked])],
    });
  }

  getDefaults() {
    this.project = {
      name: '',
    };
    this.company = {
      name: '',
    };
  }

  ngOnInit() {
    this.getUser(this.userId);
    this.getProject(this.projectId);
    this.getCompany(this.companyId);
  }

  getUser(userId) {
    this.claService.getUser(userId).subscribe(response => {
      this.user = response;
    });
  }

  getProject(projectId) {
    this.claService.getProject(projectId).subscribe(response => {
      this.project = response;
    });
  }

  getCompany(companyId) {
    this.claService.getCompany(companyId).subscribe(response => {
      this.company = response;
    });
  }

  submit() {
    // Post to CLA endpoint for user confirming their company affiliation
    // endpoint should return if they validate for that affiliation
    // if they do:
    //   Get remaining signatures needed for user
    //   getProject() -> {
    //     "project_ccla_requires_icla_signature": true|false,
    //   }
    //   if ("project_ccla_requires_icla_signature"):
    //    get
    //   if icla needs signed:
    //     send to icla page
    //   else
    //     send to return url
    // else:
    //   send user to troubleshooting page
    this.submitAttempt = true;
    this.currentlySubmitting = true;
    if (!this.form.valid) {
      this.currentlySubmitting = false;
      // prevent submit
      return;
    }
    this.openClaEmployeeCompanyTroubleshootPage();
  }

  openClaEmployeeCompanyTroubleshootPage() {
    this.navCtrl.push('ClaEmployeeCompanyTroubleshootPage', {
      projectId: this.projectId,
      repositoryId: this.repositoryId,
      userId: this.userId,
      companyId: this.company.id,
    });
  }

  goBackToPrevious() {
    this.navCtrl.pop();
  }

}
