import { Component, ChangeDetectorRef } from '@angular/core';
import { NavController, NavParams, ViewController, AlertController } from 'ionic-angular';
import { CincoService } from '../../app/services/cinco.service'

@Component({
  selector: 'contact-update',
  templateUrl: 'contact-update.html',
  providers: [CincoService]
})
export class ContactUpdate {
  projectId: string; // Always Needed
  memberId: string; // Always Needed
  member: any; // Use if passed, otherwise generate from memberId
  contactId: string; // Needed for Updating.
  contact: any; // Should be used when Updating

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    public viewCtrl: ViewController,
    public alertCtrl: AlertController,
    private changeDetectorRef: ChangeDetectorRef,
    private cincoService: CincoService
  ) {
    this.getDefaults();
    this.projectId = this.navParams.get('projectId');
    this.memberId = this.navParams.get('memberId');
    this.member = this.navParams.get('member');
    this.contactId = this.navParams.get('contactId');
    let originalContact = this.navParams.get('contact');
    // Deep copy originalContact to contact
    this.contact = Object.assign({}, originalContact);
    console.log('contact contact:');
    console.log(this.contact);
  }

  ngOnInit() {

  }

  getDefaults() {
    // Instantiate member data
    this.member = {
      name: '',
    }
    // Instantiate contact data
    this.contact = {
      email: '',
      company: '',
      name: '',
      phone: '',
      title: '',
      timezone: '',
      role: '',
      primary: 'no',
      board: 'no',
      emailGroups: [],
      photos: [],
    }
  }

  // ContactUpdate modal dismiss
  dismiss() {
    this.viewCtrl.dismiss();
  }

  primarySelectChanged(event) {
    if (event == 'yes') {
      let prompt = this.alertCtrl.create({
        title: 'Assign as Primary?',
        message: 'This will replace the exising primary contact on this project.',
        buttons: [
          {
            text: 'Cancel',
            handler: data => {
              this.contact.primary = 'no';
            }
          },
          {
            text: 'Assign',
            handler: data => {
              this.contact.primary = 'yes';
            }
          }
        ]
      });
      prompt.present();
    }
  }

  boardSelectChanged(event) {
    if (event == 'yes') {
      let prompt = this.alertCtrl.create({
        title: 'Assign to Board?',
        message: 'This will replace the exising member company Board member on this project.',
        buttons: [
          {
            text: 'Cancel',
            handler: data => {
              this.contact.board = 'no';
            }
          },
          {
            text: 'Assign',
            handler: data => {
              this.contact.board = 'yes';
            }
          }
        ]
      });
      prompt.present();
    }
  }

  removeContactPrompt() {
    let prompt = this.alertCtrl.create({
      title: 'Remove contact?',
      message: 'This will remove this contact from the project and remove them from any email alias or list.',
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
            this.removeContact();
          }
        }
      ]
    });
    prompt.present();
  }


  showPrompt() {
    let prompt = this.alertCtrl.create({
      title: 'Add Email Group',
      message: '',
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
    if(!this.contact.emailGroups) {
      this.contact.emailGroups = [];
    }
    this.contact.emailGroups.push(
      {
        name: groupName
      }
    );
  }

  removeEmailGroup(index) {
    this.contact.emailGroups.splice(index, 1);
  }

  saveContact() {
    console.log('save contact');
    console.log(this.projectId);
    console.log(this.member);
    console.log(this.contact);
    if(this.contact.id) {
      this.cincoService.updateMemberContact(this.projectId, this.member.id, this.contact.id, this.contact).subscribe(response => {
        if(response) {
          console.log('updateMemberContact response:');
          console.log(response);
        }
      });
    }
    else {
      this.cincoService.addMemberContact(this.projectId, this.member.id, this.contact).subscribe(response => {
        if(response) {
          console.log('addMemberContact response:');
          console.log(response);
        }
      });
    }

  }

  removeContact() {
    console.log('remove contact');
    console.log(this.contact);
    this.cincoService.removeMemberContact(this.projectId, this.member.id, this.contact.id).subscribe(response => {
      if(response) {
        console.log('removeMemberContact response:');
        console.log(response);
      }
    });
  }

  filesNotify(files) {
    this.contact.photos = files;
  }

}
