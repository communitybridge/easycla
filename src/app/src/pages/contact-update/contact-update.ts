import { Component, ChangeDetectorRef } from '@angular/core';
import { NavController, NavParams, ViewController, AlertController } from 'ionic-angular';
import { CincoService } from '../../app/services/cinco.service'

@Component({
  selector: 'contact-update',
  templateUrl: 'contact-update.html',
  providers: [CincoService]
})
export class ContactUpdate {
  projectId: string;
  memberId: string;
  project: any;
  member: any;
  originalContact: any;
  contact: any;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    public viewCtrl: ViewController,
    public alertCtrl: AlertController,
    private changeDetectorRef: ChangeDetectorRef,
    private cincoService: CincoService
  ) {
    this.projectId = this.navParams.get('projectId');
    this.member = this.navParams.get('member');
    console.log("contact member:");
    console.log(this.member);
    this.originalContact = this.navParams.get('contact');
    this.contact = Object.assign({}, this.originalContact);
    console.log("contact contact:");
    console.log(this.contact);
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
    console.log(this.projectId);
    console.log(this.member);
    console.log(this.contact);
    this.cincoService.updateMemberContact(this.projectId, this.member.id, this.contact.id, this.contact).subscribe(response => {
      if(response) {
        console.log("updateMemberContact response:");
        console.log(response);
      }
    });
  }

  removeContact(contact) {
    console.log("remove contact");
    console.log(this.contact);
  }

  filesNotify(files) {
    this.contact.photos = files;
  }

}
