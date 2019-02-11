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
    private authService: AuthService
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
      error => {
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
      error => {
        this.currentlySubmitting = false;
      }
    );
  }

  dismiss() {
    this.viewCtrl.dismiss();
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
