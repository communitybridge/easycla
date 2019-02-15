import { Component, ChangeDetectorRef } from "@angular/core";
import {
  NavController,
  NavParams,
  ModalController,
  ViewController,
  AlertController,
  IonicPage
} from "ionic-angular";
import { FormBuilder, FormGroup, Validators, FormArray } from "@angular/forms";
import { ClaService } from "../../services/cla.service";
import { ClaSignatureModel } from "../../models/cla-signature";

@IonicPage({
  segment: "whitelist-modal"
})
@Component({
  selector: "whitelist-modal",
  templateUrl: "whitelist-modal.html"
})
export class WhitelistModal {
  form: FormGroup;
  submitAttempt: boolean;
  currentlySubmitting: boolean;

  type: string;
  signatureId: string;
  whitelist: string[];

  constructor(
    public navParams: NavParams,
    public viewCtrl: ViewController,
    public formBuilder: FormBuilder,
    private claService: ClaService
  ) {
    this.getDefaults();
  }

  getDefaults() {
    this.type = this.navParams.get("type"); // ['email' | 'domain']
    this.signatureId = this.navParams.get('signatureId');
    this.whitelist = this.navParams.get('whitelist') || [];

    this.form = this.formBuilder.group({
      whitelist: this.formBuilder.array([])
    });
    this.submitAttempt = false;
    this.currentlySubmitting = false;
  }

  ngOnInit() {
    this.initializeWhitelist();
  }

  initializeWhitelist() {
    for (let item of this.whitelist) {
      this.addWhitelistItem(item);
    }
    if (this.whitelist.length === 0) {
      this.addNewWhitelistItem(); // auto start with one item
    }
  }

  addWhitelistItem(item) {
    let regexForItem = this.type === "domain" ? /[a-z0-9]{1,}\.[a-z]{2,}$/i : /^.+@.+\..+$/i;

    let ctrl = <FormArray>this.form.controls.whitelist;
    ctrl.push(
      this.formBuilder.group({
        whitelistItem: [item, Validators.compose([
          Validators.required,
          Validators.pattern(regexForItem)
        ])]
      })
    );
  }

  addNewWhitelistItem() {
    let regexForItem = this.type === "domain" ? /[a-z0-9]{1,}\.[a-z]{2,}$/i : /^.+@.+\..+$/i;

    let ctrl = <FormArray>this.form.controls.whitelist;
    ctrl.insert(
      0,
      this.formBuilder.group({
        whitelistItem: ["", Validators.compose([
          Validators.required,
          Validators.pattern(regexForItem)
        ])]
      })
    );
  }

  removeWhitelistItem(index) {
    let ctrl = <FormArray>this.form.controls.whitelist;
    ctrl.removeAt(index);
  }

  extractWhitelist(): string[] {
    let whitelist = [];
    for (let item of this.form.value.whitelist) {
      whitelist.push(item.whitelistItem);
    }
    return whitelist;
  }

  submit() {
    this.submitAttempt = true;
    this.currentlySubmitting = true;
    if (!this.form.valid) {
      this.currentlySubmitting = false;
      // prevent submit
      return;
    }

    let signature = new ClaSignatureModel()
    signature.signature_id = this.signatureId;

    if (this.type === "domain") {
      signature.domain_whitelist = this.extractWhitelist();
    } else {
      //email
      signature.email_whitelist = this.extractWhitelist();
    }

    this.claService.putSignature(signature).subscribe(
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
}
