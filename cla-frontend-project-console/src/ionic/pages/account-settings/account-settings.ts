// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component, ViewChild } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { CalendarLinkValidator } from '../../validators/calendarlink';
import { NavController, IonicPage, Content, ToastController } from 'ionic-angular';
import { CincoService } from '../../services/cinco.service';
import { KeycloakService } from '../../services/keycloak/keycloak.service';
import { RolesService } from '../../services/roles.service';
import { Restricted } from '../../decorators/restricted';

@Restricted({
  roles: ['isAuthenticated', 'isPmcUser']
})
// @IonicPage({
//   segment: 'account-settings'
// })
@Component({
  selector: 'account-settings',
  templateUrl: 'account-settings.html'
})
export class AccountSettingsPage {
  user: any;

  form: FormGroup;
  submitAttempt: boolean = false;
  currentlySubmitting: boolean = false;
  @ViewChild(Content) content: Content;

  loading: any;

  constructor(
    public navCtrl: NavController,
    private cincoService: CincoService,
    public formBuilder: FormBuilder,
    public toastCtrl: ToastController,
    private keycloak: KeycloakService,
    public rolesService: RolesService
  ) {
    this.getDefaults();

    this.form = formBuilder.group({
      calendar: [this.user.calendar, Validators.compose([CalendarLinkValidator.isValid])]
    });
  }

  getDefaults() {
    this.loading = {
      user: true
    };
    this.user = {
      userId: '',
      email: '',
      roles: [],
      calendar: null
    };
  }

  ngOnInit() {
    this.getCurrentUser();
  }

  getCurrentUser() {
    this.cincoService.getCurrentUser().subscribe((response) => {
      this.user = response;
      this.user.userId = this.user.lfId;
      this.form.patchValue({ calendar: this.user.calendar });
      this.loading.user = false;
    });
  }

  submitEditUser() {
    this.submitAttempt = true;
    this.currentlySubmitting = true;
    if (!this.form.valid) {
      this.content.scrollToTop();
      this.currentlySubmitting = false;
      // prevent submit
      return;
    }
    let calendar_url = this.calendarProcess(this.form.value.calendar);
    let user = {
      //userId: this.user.userId,
      lfId: this.user.userId,
      id: this.user.id,
      email: this.user.email,
      calendar: calendar_url
    };
    this.cincoService.updateUser(this.user.id, user).subscribe((response) => {
      this.currentlySubmitting = false;
      this.updateSuccess();
      this.navCtrl.setRoot(this.navCtrl['root']);
    });
  }

  updateSuccess() {
    let toast = this.toastCtrl.create({
      message: 'User updated successfully',
      showCloseButton: true,
      closeButtonText: 'Dismiss',
      duration: 3000
    });
    toast.present();
  }

  calendarProcess(calendar_str) {
    if (calendar_str === null) {
      return '';
    }
    let calendar_embed = '<iframe src="https://calendar.google.com/calendar/embed';
    if (calendar_str.startsWith(calendar_embed)) {
      // extract the calendar url with regex, url decode
      var re = /(?:src=)["']([^"']*)["']/g;
      var found = re.exec(calendar_str);
      if (found.length >= 1) {
        let parser = new DOMParser();
        let dom = parser.parseFromString('<!doctype html><body>' + found[1], 'text/html');
        calendar_str = dom.body.textContent;
      }
    }
    return calendar_str;
  }

  cancelEditUser() {
    this.navCtrl.setRoot(this.navCtrl['root']);
  }
}
