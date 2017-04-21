import { Component, ChangeDetectorRef } from '@angular/core';
import { NavController, NavParams, ViewController, AlertController } from 'ionic-angular';

@Component({
  selector: 'contact-update',
  templateUrl: 'contact-update.html'
})
export class ContactUpdate {
  project: any;
  member: any;
  contact: any;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    public viewCtrl: ViewController,
    public alertCtrl: AlertController,
    private changeDetectorRef: ChangeDetectorRef,
  ) {
    this.project=this.navParams.get('project');
    this.member=this.navParams.get('member');
    this.contact=this.navParams.get('contact');
  }

  // ContactUpdate modal dismiss
  dismiss() {
    this.viewCtrl.dismiss();
  }

  showPrompt() {
    let prompt = this.alertCtrl.create({
      title: 'Add Email Group',
      message: "",
      inputs: [
        {
          name: 'email',
          placeholder: 'Title'
        },
      ],
      buttons: [
        {
          text: 'Cancel',
          handler: data => {
            // Don't add contact
          }
        },
        {
          text: 'Save',
          handler: data => {
            this.addEmailGroup(data.email);
          }
        }
      ]
    });
    prompt.present();
  }

  addEmailGroup(groupName) {
    if(!this.contact.email_groups) {
      this.contact.email_groups = [];
    }
    this.contact.email_groups.push(
      {
        name: groupName
      }
    );
  }

  removeEmailGroup(index) {
    this.contact.email_groups.splice(index, 1);
  }

  saveContact() {
    console.log("save contact");
    console.log(this.contact);
  }

  removeContact(contact) {
    console.log("remove contact");
    console.log(this.contact);
  }

  filesNotify(files) {
    this.contact.photos = files;
  }

}
