// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component, ChangeDetectorRef, ViewChild } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { NavController, NavParams, ViewController, AlertController, IonicPage, Content } from 'ionic-angular';
import { EmailValidator } from  '../../validators/email';
import {Observable} from 'rxjs/Observable';
import 'rxjs/add/observable/forkJoin';
import { CincoService } from '../../services/cinco.service';

@IonicPage({
  segment: 'console-user-update-modal'
})
@Component({
  selector: 'console-user-update-modal',
  templateUrl: 'console-user-update-modal.html',
})
export class ConsoleUserUpdateModal {
  user: any;
  userRoles: any;
  keysGetter;
  form: FormGroup;
  submitAttempt: boolean = false;
  currentlySubmitting: boolean = false;
  @ViewChild(Content) content: Content;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    public viewCtrl: ViewController,
    public alertCtrl: AlertController,
    private changeDetectorRef: ChangeDetectorRef,
    private cincoService: CincoService,
    public formBuilder: FormBuilder
  ) {
    this.getDefaults();
    let originalUser = this.navParams.get('user');
    this.keysGetter = Object.keys;

    // Deep copy originalContact to contact
    this.user = Object.assign({}, originalUser);
    this.user.userId = this.user.lfId;
    this.form = formBuilder.group({
      userId:[this.user.userId, Validators.required],
      email:[this.user.email, Validators.compose([Validators.required, EmailValidator.isValid])],
      roles:[this.user.roles],
    });
  }

  getDefaults() {
    // Instantiate user data
    this.user = {
      id: '',
      lfId: '',
      userId: '',
      email: '',
      roles: [],
    };
    this.userRoles = {};
  }

  ngOnInit() {
    this.getUserRoles();
  }

  getUserRoles() {
    this.cincoService.getUserRoles().subscribe(response => {
      if(response) {
        this.userRoles = response;
      }
    });
  }

  // ContactUpdateModal modal dismiss
  dismiss() {
    this.viewCtrl.dismiss();
  }

  removeUserPrompt() {
    let prompt = this.alertCtrl.create({
      title: 'Remove console user?',
      message: 'This will remove this user from the console and any project they manage.',
      buttons: [
        {
          text: 'Cancel',
          handler: data => {
            // Do nothing
          }
        },
        {
          text: 'Remove',
          handler: data => {
            this.removeUser();
          }
        }
      ]
    });
    prompt.present();
  }


  saveUser() {
    this.submitAttempt = true;
    this.currentlySubmitting = true;
    if (!this.form.valid){
      this.content.scrollToTop();
      this.currentlySubmitting = false;
      // prevent submit
      return;
    }

    let userId = '';
    if (this.user.userId) {
      userId = this.user.userId;
    } else {
      userId = this.form.value.userId;
    }

    let email = '';
    if (this.user.email) {
      email = this.user.email;
    } else {
      email = this.form.value.email;
    }

    var user = {
      lfId: userId,
      email: email,
    };
    if (this.user.userId) {
      this.updateUserRoles();
    }
    else {
      // create user
      // update roles
      this.cincoService.createUser(user).subscribe(response => {
        if (response) {
          this.updateUserRoles();
        }
      });
    }

  }

  addUser() {
    this.saveUser();
  }

  updateUserRoles() {
    let userId = '';
    if (this.user.id) {
      userId = this.user.id;
    } else {
      userId = this.form.value.userId;
    }

    let prevRoles = [];
    if (this.user.roles) {
      prevRoles = this.user.roles;
    }
    var newRoles = this.form.value.roles || [];
    let observablesArray = [];
    for (let i=0; i<prevRoles.length; i++) {
      let role = prevRoles[i];
      let inNew = (newRoles.indexOf(role) !== -1);
      if(!inNew) {
        let observable = this.cincoService.removeUserRole(userId, role);
        observablesArray.push(observable);
      }
    }
    for (let i=0; i<newRoles.length; i++) {
      let role = newRoles[i];
      let inPrev = (prevRoles.indexOf(role) !== -1);
      if(!inPrev) {
        let observable = this.cincoService.addUserRole(userId, role);
        observablesArray.push(observable);
      }
    }
    if (observablesArray.length == 0) {
      this.dismiss();
    }
    Observable.forkJoin(observablesArray).subscribe(response => {
      if (response) {
        this.dismiss();
      }
    });
  }

  removeUser() {
    this.currentlySubmitting = true;
    let userId = '';
    if (this.user.id) {
      userId = this.user.id;
    } else {
      userId = this.form.value.userId;
    }
    this.cincoService.removeUser(userId).subscribe(response => {
      if(response) {
        this.dismiss();
      }
    });
  }

}
