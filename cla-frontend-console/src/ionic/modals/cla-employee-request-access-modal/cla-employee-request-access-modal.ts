import { Component, ChangeDetectorRef } from '@angular/core';
import { NavController, NavParams, ModalController, ViewController, AlertController, IonicPage } from 'ionic-angular';
import { FormBuilder, FormGroup } from '@angular/forms';
import { Validators } from '@angular/forms';
import { EmailValidator } from  '../../validators/email';
import { ClaService } from '../../services/cla.service';

@IonicPage({
  segment: 'cla/project/:projectId/repository/:repositoryId/user/:userId/employee/company/contact'
})
@Component({
  selector: 'cla-employee-request-access-modal',
  templateUrl: 'cla-employee-request-access-modal.html',
})
export class ClaEmployeeRequestAccessModal {
  projectId: string;
  repositoryId: string;
  userId: string;
  companyId: string;
  authenticated: boolean;

  userEmails: Array<string>;

  form: FormGroup;
  submitAttempt: boolean = false;
  currentlySubmitting: boolean = false;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    public modalCtrl: ModalController,
    public viewCtrl: ViewController,
    public alertCtrl: AlertController,
    private changeDetectorRef: ChangeDetectorRef,
    private formBuilder: FormBuilder,
    private claService: ClaService,
  ) {
    this.getDefaults();
    this.projectId = navParams.get('projectId');
    this.repositoryId = navParams.get('repositoryId');
    this.userId = navParams.get('userId');
    this.companyId = navParams.get('companyId');
    this.authenticated = navParams.get('authenticated'); 
    this.form = formBuilder.group({
      email:['', Validators.compose([Validators.required, EmailValidator.isValid])],
      message:[''], // Validators.compose([Validators.required])
    });
  }

  getDefaults() {
    this.userEmails = [];
  }

  ngOnInit() {
    this.getUser();
  }
 
  getUser() {
    if (this.authenticated) { //Gerrit User
      this.claService.getUserWithAuthToken(this.userId).subscribe(user => {
        if(user.lf_email) {
          this.userEmails.push(user.lf_email)
        }
      })
    } else { // Github User
      this.claService.getUser(this.userId).subscribe(user => {
        this.userEmails = user.user_emails;
      });
    }
  }

  // ContactUpdateModal modal dismiss
  dismiss() {
    this.viewCtrl.dismiss();
  }

  submit() {
    this.submitAttempt = true;
    this.currentlySubmitting = true;
    if (!this.form.valid) {
      this.currentlySubmitting = false;
      // prevent submit
      return;
    }
    let message = {
      user_email: this.form.value.email,
      message: this.form.value.message,
    };
    this.claService.postUserMessageToCompanyManager(this.userId, this.companyId, message).subscribe(response => {
      this.emailSent();

    });
  }

  emailSent() {
    let alert = this.alertCtrl.create({
      title: 'E-Mail Successfully Sent!',
      subTitle: 'Thank you for contacting your CLA Manager. Once you are authorized, you will have to complete the CLA process from your existing pull request.',
      buttons: ['Dismiss']
    });
    alert.onDidDismiss(() => this.dismiss());
    alert.present();
  }

}
