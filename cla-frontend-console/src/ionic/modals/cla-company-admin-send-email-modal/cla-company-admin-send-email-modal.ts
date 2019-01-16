import { Component, ChangeDetectorRef } from '@angular/core';
import { NavController, NavParams, ModalController, ViewController, AlertController, IonicPage } from 'ionic-angular';
import { FormBuilder, FormGroup } from '@angular/forms';
import { Validators } from '@angular/forms';
import { EmailValidator } from  '../../validators/email';
import { ClaService } from '../../services/cla.service';
import { EnvConfig } from '../../services/cla.env.utils';

@IonicPage({
  segment: 'cla/project/:projectId/user/:userId/invite-company-admin'
})
@Component({
  selector: 'cla-company-admin-send-email-modal',
  templateUrl: 'cla-company-admin-send-email-modal.html',
})
export class ClaCompanyAdminSendEmailModal {
  projectId: string;
  repositoryId: string;
  userId: string;
  consoleLink: string; 
  authenticated: boolean; // true if coming from gerrit/corporate 

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
    this.userId = navParams.get('userId');
    this.authenticated = navParams.get('authenticated');
    this.form = formBuilder.group({
      useremail:['', Validators.compose([Validators.required, EmailValidator.isValid])],
      adminemail:['', Validators.compose([Validators.required, EmailValidator.isValid])],
      adminname:[''],
    });
  }

  getDefaults() {
    this.userEmails = [];
    this.consoleLink = EnvConfig['corp-console-link'];
  }

  ngOnInit() {
    this.getUser();
  }

  getUser() {
    if (this.authenticated) {
      this.claService.getUserWithAuthToken(this.userId).subscribe(user => {
        if(user.lf_email) {
          this.userEmails.push(user.lf_email)
        }
      })
    } else {
      this.claService.getUser(this.userId).subscribe(user => {
        this.userEmails = user.user_emails;
      });
    }
  }

  dismiss() {
    this.viewCtrl.dismiss();
  }

  emailSent() {
    let alert = this.alertCtrl.create({
      title: 'E-Mail Successfully Sent!',
      subTitle: 'An E-Mail has been successfully sent. Please wait for your corporate administrator to add you to your company whitelist',
      buttons: ['Dismiss']
    });
    alert.present();
  }

  submit() {
    this.submitAttempt = true;
    this.currentlySubmitting = true;
    if (!this.form.valid) {
      this.currentlySubmitting = false;
      // prevent submit
      return;
    }
    let data = {
      user_email: this.form.value.useremail,
      admin_email: this.form.value.adminemail,
      admin_name: this.form.value.adminname,
    };
    this.claService.postEmailToCompanyAdmin(this.userId, data).subscribe(response => {
      this.emailSent();
      this.dismiss();
    });
  }

}
