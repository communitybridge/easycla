import {Component} from "@angular/core";
import {
  ViewController,
  IonicPage, NavParams
} from "ionic-angular";
import {FormBuilder, FormGroup, Validators} from "@angular/forms";
import {ClaService} from "../../services/cla.service";
import {AuthService} from "../../services/auth.service";
import {ClaCompanyModel} from "../../models/cla-company";

@IonicPage({
  segment: "add-manager-modal"
})
@Component({
  selector: "add-manager-modal",
  templateUrl: "add-manager-modal.html"
})
export class AddManagerModal {
  form: FormGroup;
  submitAttempt: boolean = false;

  company: ClaCompanyModel;
  managerLFID: string;

  constructor(public viewCtrl: ViewController,
              public navParams: NavParams,
              public formBuilder: FormBuilder,
              private claService: ClaService) {
    this.company = this.navParams.get("company");

    this.form = this.formBuilder.group({
      managerLFID: [this.managerLFID, Validators.compose([Validators.required])],
    });
  }


  submit() {
    this.submitAttempt = true;
    this.addManager()
  }

  addManager() {
    this.claService.postCompanyManager(this.company.company_id, this.form.getRawValue())
      .subscribe(() => this.dismiss(true));
  }

  dismiss(data = false) {
    this.viewCtrl.dismiss(data);
  }
}
