import { Component, Renderer, } from '@angular/core';
import { NavController, NavParams, ViewController, AlertController, ToastController, IonicPage  } from 'ionic-angular';
import { CincoService } from '../../app/services/cinco.service'

@IonicPage({
  segment: 'project-user-management'
})
@Component({
  selector: 'project-user-management',
  templateUrl: 'project-user-management.html',
  providers: [CincoService]
})
export class ProjectUserManagementModal {
  projectId: string; // Always Needed
  users: any;
  selectedUsers: any;
  enteredUserId: string;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    public viewCtrl: ViewController,
    private cincoService: CincoService,
    public toastCtrl: ToastController,
    private renderer: Renderer,
    public alertCtrl: AlertController,
  ) {
    this.projectId = this.navParams.get('projectId');
    this.getDefaults();
  }

  ngOnInit() {
    this.getProjectConfig();
  }

  getDefaults() {
    this.users = [];
    this.selectedUsers = [];
  }

  // ContactUpdate modal dismiss
  dismiss() {
    this.viewCtrl.dismiss();
  }

  getProjectConfig() {
    this.cincoService.getProjectConfig(this.projectId).subscribe(response => {
      if(response) {
        this.users = response.programManagers;
      }
    });
  }

  selectUser(event, user) {
    event.stopPropagation();

      // deselect the entire selected users array
      this.deselectUsers(this.selectedUsers);
      user.selected = true;
      this.selectedUsers = [user];
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

  removeUser(index) {
    this.users.splice(index, 1);
    this.cincoService.updateProjectManagers(this.projectId, this.users).subscribe(response => {
      if(response) {
        // console.log("response");
        // console.log(response);
        // this.users = response;
      }
    });
  }

  modalClick(event) {
    // stray unhandled/unprevented click. deselect all users
    this.deselectUsers(this.selectedUsers);
  }

  addprojectManager(userId) {
    this.enteredUserId = '';
    this.users.push(userId);
    // console.log(this.users);
    this.cincoService.updateProjectManagers(this.projectId, this.users).subscribe(response => {
      if(response) {
        // console.log("response");
        // console.log(response);
        // this.users = response;
      }
    });
  }

}
