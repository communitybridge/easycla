import { Component, ChangeDetectorRef } from '@angular/core';
import { NavController, NavParams, ViewController, AlertController, IonicPage } from 'ionic-angular';
import { CincoService } from '../../app/services/cinco.service';

@IonicPage({
  segment: 'contact-update'
})
@Component({
  selector: 'contact-update',
  templateUrl: 'contact-update.html',
  providers: [CincoService]
})
export class ContactUpdate {
  projectId: string; // Always Needed
  memberId: string; // Always Needed
  org: any;
  contact: any;
  contactId: string;
  roleId: string;
  memberContactRoles: any;
  orgContactRoles: any;
  keysGetter;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    public viewCtrl: ViewController,
    public alertCtrl: AlertController,
    private changeDetectorRef: ChangeDetectorRef,
    private cincoService: CincoService
  ) {
    this.getDefaults();
    this.keysGetter = Object.keys;
    this.projectId = this.navParams.get('projectId');
    this.memberId = this.navParams.get('memberId');
    this.org = this.navParams.get('org');
    let originalContact = this.navParams.get('contact');

    if (originalContact.id) {
      this.roleId = originalContact.id;
    }

    if (originalContact.contact.id) {
      this.contactId = originalContact.contact.id;
    }

    // Deep copy originalContact to contact
    this.contact = Object.assign({}, originalContact);
  }

  ngOnInit() {
    this.getMemberContactRoles();
    this.getOrgContactRoles();
  }

  getDefaults() {
    // Instantiate member data
    this.org = {
      name: '',
    };
    // Instantiate contact data
    this.contact = {
      title:"",
      type:"",
      primaryContact:false,
      boardMember:false,
      contact:{
        email:"",
        givenName:"",
        familyName:"",
        phone:"",
        type:"",
        bio:"",
      },
    };

    this.memberContactRoles = {};
    this.orgContactRoles = {};

  }

  getOrgContactRoles() {
    this.cincoService.getOrganizationContactTypes().subscribe(response => {
      if(response) {
        this.orgContactRoles = response;
      }
    });
  }

  getMemberContactRoles() {
    this.cincoService.getMemberContactRoles().subscribe(response => {
      if(response) {
        this.memberContactRoles = response;
      }
    });
  }

  // ContactUpdate modal dismiss
  dismiss() {
    this.viewCtrl.dismiss();
  }

  primarySelectChanged(event) {
    if (event == 'true') {
      let prompt = this.alertCtrl.create({
        title: 'Assign as Primary?',
        message: 'This will replace the exising primary contact on this project.',
        buttons: [
          {
            text: 'Cancel',
            handler: data => {
              this.contact.primaryContact = 'false';
            }
          },
          {
            text: 'Assign',
            handler: data => {
              this.contact.primaryContact = 'true';
            }
          }
        ]
      });
      prompt.present();
    }
  }

  boardSelectChanged(event) {
    if (event == 'true') {
      let prompt = this.alertCtrl.create({
        title: 'Assign to Board?',
        message: 'This will replace the exising member company Board member on this project.',
        buttons: [
          {
            text: 'Cancel',
            handler: data => {
              this.contact.boardMember = 'false';
            }
          },
          {
            text: 'Assign',
            handler: data => {
              this.contact.boardMember = 'true';
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
    if (this.contactId) {
      if (this.roleId) {
        this.cincoService.updateOrganizationContact(this.org.id, this.contactId, this.contact.contact).subscribe(response => {
          if (response) {
            // update org contact with response from update
            // should be the same as what was sent, but we will just be sure
            this.contact.contact = response;
            // add as a member contact
            this.cincoService.updateMemberContact(this.projectId, this.memberId, this.contactId, this.roleId, this.contact).subscribe(response => {
              if(response) {
                this.dismiss();
              }
            });
          }
        });
      }
      else {
        this.cincoService.updateOrganizationContact(this.org.id, this.contactId, this.contact.contact).subscribe(response => {
          if (response) {
            // update org contact with response from update
            // should be the same as what was sent, but we will just be sure
            this.contact.contact = response;
            // add as a member contact
            this.cincoService.addMemberContact(this.projectId, this.memberId, this.contactId, this.contact).subscribe(response => {
              if (response) {
                this.dismiss();
              }
            });
          }
        });
      }
    }
    else {
      // Add new contact to organization
      this.cincoService.createOrganizationContact(this.org.id, this.contact.contact).subscribe(response => {
        if (response) {
          this.contactId = response;
          this.contact.contact.id = this.contactId;
          // add to member
          this.cincoService.addMemberContact(this.projectId, this.memberId, this.contactId, this.contact).subscribe(response => {
            if (response) {
              this.dismiss();
            }
          });
        }
      });
    }

  }

  removeContact() {
    this.cincoService.removeMemberContact(this.projectId, this.memberId, this.contactId, this.roleId).subscribe(response => {
      if(response) {
        this.dismiss();
      }
    });
  }

  filesNotify(files) {
    this.contact.contact.photos = files;
  }

}
