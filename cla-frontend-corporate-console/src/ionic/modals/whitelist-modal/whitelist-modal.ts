// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import {Component} from "@angular/core";
import {IonicPage, NavParams, ViewController} from "ionic-angular";
import {FormArray, FormBuilder, FormGroup, Validators} from "@angular/forms";
import {ClaService} from "../../services/cla.service";
import {ClaSignatureModel} from "../../models/cla-signature";

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
    this.type = this.navParams.get("type"); // ['email' | 'domain' | 'github']
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

  /**
   * Returns the regular expression for the form type.
   * @param formType the form type, e.g.: domain, email, github username
   */
  getValidationRegExp(formType: string): RegExp {
    let regex: RegExp;
    if (formType === "domain") {
      // domains
      regex = new RegExp(/[a-z0-9]{1,}\.[a-z]{2,}$/i);
    } else if (formType === "email") {
      // emails
      regex = new RegExp(/^.+@.+\..+$/i);
    } else {
      // github usernames - allow the standard GitHub characters - plus
      // allow square brackets - which are allowed for bots - ex: 'somename[bot]'
      regex = new RegExp(/^[a-z\[\d](?:[a-z\[\]\d]|-(?=[a-z\[\]\d])){0,38}$/i);
    }

    return regex;
  }

  addWhitelistItem(item) {
    let ctrl = <FormArray>this.form.controls.whitelist;
    ctrl.push(
      this.formBuilder.group({
        whitelistItem: [item, Validators.compose([
          Validators.required,
          Validators.pattern(this.getValidationRegExp(this.type))
        ])]
      })
    );
  }

  /**
   * Called on new items added to the list
   */
  addNewWhitelistItem() {
    let ctrl = <FormArray>this.form.controls.whitelist;
    ctrl.insert(
      0,
      this.formBuilder.group({
        whitelistItem: ["", Validators.compose([
          Validators.required,
          Validators.pattern(this.getValidationRegExp(this.type))
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

    let signature = new ClaSignatureModel();
    signature.signature_id = this.signatureId;

    if (this.type === "domain") {
      signature.domain_whitelist = this.extractWhitelist();
    } else if (this.type === "email") {
      //email
      signature.email_whitelist = this.extractWhitelist();
    } else {
      //github username
      signature.github_whitelist = this.extractWhitelist();
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
