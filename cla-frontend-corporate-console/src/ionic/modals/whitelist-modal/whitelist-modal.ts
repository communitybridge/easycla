// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component } from '@angular/core';
import { IonicPage, NavParams, ViewController } from 'ionic-angular';
import {
  AbstractControl,
  AsyncValidatorFn,
  FormArray,
  FormBuilder,
  FormGroup,
  ValidationErrors,
  ValidatorFn
} from '@angular/forms';
import { ClaService } from '../../services/cla.service';
import { ClaSignatureModel } from '../../models/cla-signature';
import { Observable } from 'rxjs';

@IonicPage({
  segment: 'whitelist-modal'
})
@Component({
  selector: 'whitelist-modal',
  templateUrl: 'whitelist-modal.html'
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
  message: any;

  constructor(
    public navParams: NavParams,
    public viewCtrl: ViewController,
    public formBuilder: FormBuilder,
    private claService: ClaService
  ) {
    this.getDefaults();
  }

  getDefaults() {
    this.type = this.navParams.get('type'); // ['email' | 'domain' | 'github']
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
    this.message = {
      error: null
    };
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
    if (formType === 'domain') {
      // domains
      regex = new RegExp(/[a-z0-9]{1,}\.[a-z]{2,}$/i);
    } else if (formType === 'email') {
      // emails
      regex = new RegExp(/^.+@.+\..+$/i);
    } else if (formType === 'github') {
      // github usernames - allow the standard GitHub characters - plus
      // allow square brackets - which are allowed for bots - ex: 'somename[bot]'
      regex = new RegExp(/^[a-z\[\d](?:[a-z\[\]\d]|-(?=[a-z\[\]\d])){0,38}$/i);
    } else if (formType === 'githubOrg') {
      regex = new RegExp(/^[a-z\[\d](?:[a-z\-\d]|-(?=[a-z\-\d])){0,38}$/i);
    }

    return regex;
  }

  myFormValidator(entries: string[], re: RegExp): ValidatorFn {
    // Return a function with the implementation - we do this so can can pass in args, otherwise we don't have access
    // to 'this'
    return (control: AbstractControl): ValidationErrors | null => {
      if (control.value == null || control.value === '') {
        return null;
      }

      // Create a new array from the original - add the new value being entered on the form
      const newArray: string[] = Array.from(entries);
      newArray.push(control.value);

      // Convert the list of values to a set and back again - this will remove duplicates
      const noDuplicates: string[] = Array.from(new Set(newArray));

      // If the array lengths do not match, then we have at least one duplicate entry
      if (newArray.length !== noDuplicates.length) {
        this.message.error = 'Duplicate Entry: ' + control.value;
        return {
          duplicate: control.value
        };
      }

      // Test the regex pattern
      if (!re.test(control.value)) {
        this.message.error = 'Invalid Pattern: ' + control.value;
        return {
          invalidPattern: control.value
        };
      }

      // Default is to pas
      this.message.error = null;
      return null;
    };
  }

  /**
   * A GitHub organization entry validator. Returns a validation error if the
   * specified GitHub organization does not exist, otherwise returns null.
  githubOrgValidator(claService: ClaService): AsyncValidatorFn {
    return (control: AbstractControl): Promise<ValidationErrors | null> | Observable<ValidationErrors | null> => {
      return claService.doesGitHubOrgExist(control.value).map(res => {
        // Response will be either:
        // - a full JSON document (valid) { "login": "deal-test-org-2", "id": 56558033, "node_id": "MDEyOk9yZ2FuaXphdGlvbjU2NTU4MDMz", ...}
        // or
        // - a small JSON document (invalid) {"message": "Not Found", "documentation_url": "https://..."}
        console.log('Validating: ' + control.value);
        console.log('Response: ');
        console.log(res);
        if (res.message) {
          this.message.error = 'Invalid Organization: ' + control.value;
          return {invalidOrg: true};
        } else {
          this.message.error = null;
          return null;
        }
      });
    }
  }
   */

  addWhitelistItem(item) {
    let ctrl = <FormArray>this.form.controls.whitelist;
    ctrl.push(
      this.formBuilder.group({
        whitelistItem: [item, [this.myFormValidator(this.extractWhitelist(), this.getValidationRegExp(this.type))]]
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
        whitelistItem: ['', [this.myFormValidator(this.extractWhitelist(), this.getValidationRegExp(this.type))]]
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
      // Don't add empty stuff
      if (item.whitelistItem == null || item.whitelistItem === '') {
        continue;
      }
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

    //this.message.error = null;
    let signature = new ClaSignatureModel();
    signature.signature_project_id = this.projectId;
    signature.signature_reference_id = this.companyId; // CCLA, so signature_reference_id is the company id
    signature.signature_id = this.signatureId;

    // ['email' | 'domain' | 'github' | 'githubOrg']
    if (this.type === 'domain') {
      signature.domain_whitelist = this.extractWhitelist();
    } else if (this.type === 'email') {
      //email
      signature.email_whitelist = this.extractWhitelist();
    } else if (this.type === 'github') {
      //github username
      signature.github_whitelist = this.extractWhitelist();
    } else if (this.type === 'githubOrg') {
      //github organization
      signature.github_org_whitelist = this.extractWhitelist();
    }

    this.claService.putSignature(signature).subscribe(
      (response) => {
        this.currentlySubmitting = false;
        this.dismiss();
      },
      (error) => {
        this.currentlySubmitting = false;
      }
    );
  }

  dismiss() {
    //this.message.error = null;
    this.viewCtrl.dismiss();
  }

  saveButton(): string {
    if (this.currentlySubmitting) {
      return 'gray';
    } else {
      return 'secondary';
    }
  }

  cancelButton(): string {
    if (this.currentlySubmitting) {
      return 'gray';
    } else {
      return 'danger';
    }
  }

  /**
   * Returns the type as a pretty string with the proper case.
   */
  getTitleType(): string {
    if (this.type === 'domain') {
      return 'Domain';
    } else if (this.type === 'email') {
      return 'Email';
    } else if (this.type === 'github') {
      return 'GitHub Username';
    } else if (this.type === 'githubOrg') {
      return 'GitHub Org';
    } else {
      return 'Unknown';
    }
  }
}
