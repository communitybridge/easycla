import { Component } from '@angular/core';
import { NavController, NavParams, IonicPage, ModalController, } from 'ionic-angular';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { CheckboxValidator } from  '../../validators/checkbox';

@IonicPage({
  segment: 'cla/project/:projectId/repository/:repositoryId/user/:userId/employee/company/:companyId/confirm'
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

  }

  ngOnInit() {
    this.getProject();
    this.getCompany();
  }

  getProject() {
    this.project = {
      id: '0000000001',
      name: 'Project Name',
      logoRef: 'https://dummyimage.com/225x102/d8d8d8/242424.png&text=Project+Logo',
    };
  }

  getCompany() {
    this.company = {
      name: 'Company Name',
      id: '0000000001',
    };
  }

  submit() {
    // Post to CLA endpoint for user confirming their company affiliation
    // endpoint should return if they validate for that affiliation
    // if they do:
    //   Get remaining signatures needed for user
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
