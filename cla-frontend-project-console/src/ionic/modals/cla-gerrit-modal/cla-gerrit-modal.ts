// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component } from '@angular/core';
import { NavParams, ViewController, IonicPage, Events } from 'ionic-angular';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { ExtraValidators } from '../../validators/requireSelfAnd';
import { ClaService } from '../../services/cla.service';
import { PlatformLocation } from '@angular/common';
import { generalConstants } from '../../constants/general';

@IonicPage({
  segment: 'cla-gerrit-modal'
})
@Component({
  selector: 'cla-gerrit-modal',
  templateUrl: 'cla-gerrit-modal.html'
})
export class ClaGerritModal {
  form: FormGroup;
  submitAttempt: boolean = false;
  currentlySubmitting: boolean = false;
  projectId: string;
  user: any;
  errorMessage: string;

  constructor(
    public navParams: NavParams,
    public viewCtrl: ViewController,
    private formBuilder: FormBuilder,
    public claService: ClaService,
    public events: Events,
    private location: PlatformLocation
  ) {
    this.projectId = this.navParams.get('projectId');
    this.location.onPopState(() => {
      this.viewCtrl.dismiss(false);
    });
    this.form = formBuilder.group({
      gerritName: ['', Validators.compose([
        Validators.required,
        Validators.pattern(new RegExp(generalConstants.GERRIT_NAME_REGEX)),
        Validators.maxLength(99)
      ])],
      URL: ['', Validators.compose([Validators.required])],
      groupIdIcla: [
        '',
        (control) => {
          return ExtraValidators.requireSelfOr(control, 'groupIdCcla');
        }
      ],
      groupIdCcla: [
        '',
        (control) => {
          return ExtraValidators.requireSelfOr(control, 'groupIdIcla');
        }
      ]
    });

    events.subscribe('modal:close', () => {
      this.dismiss();
    });
  }

  ngOnInit() { }

  dismiss(data?) {
    this.viewCtrl.dismiss(data);
  }

  submit() {
    this.submitAttempt = true;
    this.currentlySubmitting = true;
    if (!this.form.valid) {
      this.currentlySubmitting = false;
      // prevent submit
      return;
    }
    this.postGerritInstance();
  }

  postGerritInstance() {
    let gerrit = {
      gerritName: this.form.value.gerritName,
      gerritUrl: this.form.value.URL,
      groupIdCcla: this.form.value.groupIdCcla,
      groupIdIcla: this.form.value.groupIdIcla,
      version: "v1"
    };
    this.claService.postGerritInstance(this.projectId, gerrit).subscribe(
      (response) => {
        this.dismiss(true);
      },
      (error) => {
        console.log(error)
        if (error.status === 422) {
          this.errorMessage = 'Invalid Gerrit Instance URL.';
        } else if (error._body) {
          this.errorMessage = JSON.parse(error._body).Message;
        }
      }
    );
  }
}
