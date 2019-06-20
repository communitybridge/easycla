// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: AGPL-3.0-or-later

import { Component, ChangeDetectorRef } from "@angular/core";
import {
  NavController,
  NavParams,
  ModalController,
  ViewController,
  AlertController,
  IonicPage
} from "ionic-angular";
import { FormBuilder, FormGroup, Validators } from "@angular/forms";
import { ClaService } from "../../services/cla.service";
import { ClaCompanyModel } from "../../models/cla-company";
import { AuthService } from "../../services/auth.service";
import { HttpErrorResponse } from "@angular/common/http"
import { ErrorObservable } from "rxjs/observable/ErrorObservable";

@IonicPage({
  segment: "add-company-modal"
})
@Component({
  selector: "add-company-modal",
  templateUrl: "add-company-modal.html"
})
export class AddCompanyModal {
  form: FormGroup;
  submitAttempt: boolean = false;
  currentlySubmitting: boolean = false;

  mode: string;
  company: ClaCompanyModel;
  companyName: string;
  userEmail: string;
  userName: string;

  constructor(
    public navParams: NavParams,
    public viewCtrl: ViewController,
    public formBuilder: FormBuilder,
    private claService: ClaService,
    private authService: AuthService,
    public alertCtrl: AlertController
  ) {
    this.getDefaults();
  }

  getDefaults() {
    this.company = this.navParams.get("company");
    this.mode = this.company ? "edit" : "add";

    this.form = this.formBuilder.group({
      companyName: [this.companyName, Validators.compose([Validators.required])],
      userEmail: [this.userEmail, Validators.compose([Validators.required])],
      userName: [this.userName, Validators.compose([Validators.required])]
    });
  }

  ngOnInit() {
  }

  ionViewDidEnter() {
    this.updateUserInfoBasedLFID();
  }

  submit() {
    this.submitAttempt = true;
    this.currentlySubmitting = true;
    if (this.mode === "add") {
      this.addCompany();
    } else {
      this.updateCompany();
    }
  }

  addCompany() {
    let company = {
      company_name: this.companyName,
      company_manager_user_email: this.userEmail,
      company_manager_user_name: this.userName
    };
    this.claService.postCompany(company).subscribe(
      response => {
        this.currentlySubmitting = false;
        this.dismiss();
      },
      (err: any) => {
        if (err.status === 409){
          let errJSON = err.json();
          this.companyExistAlert(errJSON.company_id)
        }
        this.currentlySubmitting = false;
      }
    );
  }

  updateCompany() {
    let company = {
      company_id: this.company.company_id,
      company_name: this.companyName
    };
    this.claService.putCompany(company).subscribe(
      response => {
        this.currentlySubmitting = false;
        this.dismiss();
      },
      (err: any) => {
        if (err.status === 409){
          let errJSON = err.json();
          this.companyExistAlert(errJSON.company_id)
        }
        this.currentlySubmitting = false;
      }
    );
  }

  dismiss() {
    this.viewCtrl.dismiss();
  }

  companyExistAlert(company_id) {
    let alert = this.alertCtrl.create({
      title: 'Company ' + this.companyName + ' already exists',
      message: 'The company you tried to create already exists in the CLA system. Would you like to request access?',
      buttons: [
        {
          text: 'Request',
          handler: () => {
            this.claService.sendInviteRequestEmail(company_id)
             .subscribe(() => this.dismiss());
          }
        },
        {
          text: 'Cancel',
          role: 'cancel',
          handler: () => {
            console.log('No clicked');
          }
        }
      ]
    });
    alert.present();
  }

  private updateUserInfoBasedLFID() {
    if (this.authService.isAuthenticated()) {
      this.authService.getIdToken()
        .then(token => {
          return this.authService.parseIdToken(token);
        })
        .then(tokenParsed => {
          if (tokenParsed && tokenParsed["email"]) {
            this.userEmail = tokenParsed["email"];
          }
          if (tokenParsed && tokenParsed["name"]) {
            this.userName = tokenParsed["name"];
          }
        })
        .catch(error => {
          console.log(JSON.stringify(error));
          return;
        });
    }
    return;
  }
}
