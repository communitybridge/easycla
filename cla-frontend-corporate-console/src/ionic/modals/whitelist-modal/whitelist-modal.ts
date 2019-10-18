// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import {Component} from "@angular/core";
import {IonicPage, NavParams, ViewController} from "ionic-angular";
import {
  AbstractControl,
  FormArray,
  FormBuilder,
  FormGroup,
  ValidationErrors,
  ValidatorFn,
  Validators
} from "@angular/forms";
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
  companyName: string;
  projectName: string;
  projectId: string;
  companyId: string;
  signatureId: string;
  whitelist: string[];
  errorMessage: string;

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
    this.companyName = this.navParams.get('companyName');
    this.projectName = this.navParams.get('projectName');
    this.projectId = this.navParams.get('projectId');
    this.companyId = this.navParams.get('companyId');
    this.signatureId = this.navParams.get('signatureId');
    this.whitelist = this.navParams.get('whitelist') || [];

    this.form = this.formBuilder.group({
      whitelist: this.formBuilder.array([])
    });
    this.submitAttempt = false;
    this.currentlySubmitting = false;
    this.errorMessage = null;
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

  /**
   * A duplicate entry validator. Returns a validation error if a duplicate entry is detected, returns null otherwise.
   *
   */
  duplicateEntryValidator(entries: string[]): ValidatorFn {
    // Return a function with the implementation - we do this so can can pass in args, otherwise we don't have access
    // to 'this'
    return (control: AbstractControl): ValidationErrors | null => {
      // Create a new array from the original - add the new value being entered on the form
      const newArray: string[] = Array.from(entries);
      newArray.push(control.value);

      // Convert the list of values to a set and back again - this will remove duplicates
      const noDuplicates: string[] = Array.from(new Set(newArray));

      // If the array lengths do not match, then we have at least one duplicate entry
      if (newArray.length !== noDuplicates.length) {
        this.errorMessage = 'Duplicate Entry: ' + control.value;
        return {
          duplicate: control.value
        }
      }

      this.errorMessage = '';
      return null;
    }
  }

  addWhitelistItem(item) {
    this.errorMessage = null;
    let ctrl = <FormArray>this.form.controls.whitelist;
    ctrl.push(
      this.formBuilder.group({
        whitelistItem: [item, Validators.compose([
          Validators.required,
          Validators.pattern(this.getValidationRegExp(this.type)),
          this.duplicateEntryValidator(this.extractWhitelist())
        ])]
      })
    );
  }

  /**
   * Called on new items added to the list
   */
  addNewWhitelistItem() {
    this.errorMessage = null;
    let ctrl = <FormArray>this.form.controls.whitelist;
    ctrl.insert(
      0,
      this.formBuilder.group({
        whitelistItem: ["", Validators.compose([
          Validators.required,
          Validators.pattern(this.getValidationRegExp(this.type)),
          this.duplicateEntryValidator(this.extractWhitelist())
        ])]
      })
    );
  }

  removeWhitelistItem(index) {
    this.errorMessage = null;
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

    this.errorMessage = null;
    let signature = new ClaSignatureModel();
    signature.signature_project_id = this.projectId;
    signature.signature_reference_id = this.companyId; // CCLA, so signature_reference_id is the company id
    signature.signature_id = this.signatureId;

    // ['email' | 'domain' | 'github']
    if (this.type === "domain") {
      // TODO: apply filter to remove duplicates
      signature.domain_whitelist = this.extractWhitelist();
    } else if (this.type === "email") {
      //email
      // TODO: apply filter to remove duplicates
      signature.email_whitelist = this.extractWhitelist();
    } else {
      //github username
      // TODO: apply filter to remove duplicates
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
    this.errorMessage = null;
    this.viewCtrl.dismiss();
  }

  saveButton(): string {
    if (this.currentlySubmitting) {
      return 'gray';
    } else {
      return 'primary';
    }
  }
  cancelButton(): string {
    if (this.currentlySubmitting) {
      return 'gray';
    } else {
      return 'primary';
    }
  }
}
