// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: AGPL-3.0-or-later

import { Component, Renderer, ElementRef, ViewChild, } from '@angular/core';
import { NavController, NavParams, ViewController, AlertController, ToastController, IonicPage  } from 'ionic-angular';
import { CincoService } from '../../services/cinco.service'

@IonicPage({
  segment: 'assign-user-modal'
})
@Component({
  selector: 'assign-user-modal',
  templateUrl: 'assign-user-modal.html',
})
export class AssignUserModal {
  projectId: string; // Always Needed
  users: any;
  folders: any;
  selectedUsers: any;
  loading: any;

  /**
   * Comma separated array of allowed user extensions
   */
  uploadTypes: string;

  /**
   * Maximum number of bytes for upload size
   */
  uploadSizeMax: number;

  /**
   * Native upload button (hidden)
   */
  @ViewChild('input')
  private nativeInputBtn: ElementRef;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    public viewCtrl: ViewController,
    private cincoService: CincoService,
    public toastCtrl: ToastController,
    private renderer: Renderer,
    public alertCtrl: AlertController,
  ) {
    this.selectedUsers = [];
    this.getDefaults();
  }

  ngOnInit() {
    // this.getAllUsers();
  }

  getDefaults() {
    this.loading = {
      users: true,
    };
    this.users = [
      {
        userId:'',
        email:'',
        roles:[],
        calendar:null
      },
    ];
  }

  // getAllUsers() {
  //   this.cincoService.getAllUsers().subscribe(response => {
  //     if(response) {
  //       this.users = response;
  //       this.loading.users = false;
  //     }
  //   });
  // }

  // ContactUpdateModal modal dismiss
  dismiss() {
    this.viewCtrl.dismiss();
  }

  dismissWithData(event) {
    event.stopPropagation();
    this.viewCtrl.dismiss({
      selectedUsers: this.selectedUsers,
    });
  }

  selectUser(event, user) {
    event.stopPropagation();

    if (event.ctrlKey) {
      if (user.selected) {
        this.deselectUsers([user]);
      }
      else {
        user.selected = true;
        this.selectedUsers.push(user);
      }
    }
    else { // standard single user select
      // deselect the entire selected users array
      this.deselectUsers(this.selectedUsers);
      user.selected = true;
      this.selectedUsers = [user];
    }
  }

  downloadSelected(event) {
    event.stopPropagation();

  }

  previewSelected(event) {
    event.stopPropagation();

  }

  deleteSelected(event) {
    event.stopPropagation();
    let prompt_title = '';
    if(this.selectedUsers.length > 1) {
      prompt_title = 'Delete users?';
    }
    else {
      prompt_title = 'Delete user?';
    }
    let prompt_message = this.selectedUsers.map(function(user){
      return user.name;
    }).join(',<br>');
    let prompt = this.alertCtrl.create({
      title: prompt_title,
      message: prompt_message,
      buttons: [
        {
          text: 'Cancel',
          handler: data => {
            // Do nothing
          }
        },
        {
          text: 'Delete',
          handler: data => {
            // TODO: Make cinco calls to delete users
          }
        }
      ]
    });
    prompt.present();
  }

  /*
    Helper function to stop propagation on elements
  */
  stopEventPropagation(event) {
    event.stopPropagation();
  }

  deselectUsers(users) {
    if (!users) {
      return;
    }
    let i = users.length;
    while (i--) {
      let user = users[i];
      user.selected = false;
      let index = this.selectedUsers.indexOf(user);
      if(index !== -1) {
        this.selectedUsers.splice(index, 1);
      }
    }
  }

  modalClick(event) {
    // stray unhandled/unprevented click. deselect all users
    this.deselectUsers(this.selectedUsers);
  }

  /**
  * Callback executed when the visible button is pressed
  * @param  {Event}  event should be a mouse click event
  */
  uploadClicked(event: Event) {

    // trigger click event of hidden input
    let clickEvent: MouseEvent = new MouseEvent("click", {bubbles: true});
    this.renderer.invokeElementMethod(
      this.nativeInputBtn.nativeElement, "dispatchEvent", [clickEvent]
    );
  }

}
