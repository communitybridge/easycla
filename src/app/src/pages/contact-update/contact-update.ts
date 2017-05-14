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
  contact: any; // Should be used when Updating
  contactId: string;
  roleId: string;
  memberContactRoles: any;
  orgContactRoles: any;


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
    this.org = this.navParams.get('org');
    let originalContact = this.navParams.get('contact');

    if (originalContact.id) {
      this.contactId = originalContact.id;
    }

    if (originalContact.role) {
      this.roleId = originalContact.role;
    }

    // Deep copy originalContact to contact
    this.contact = Object.assign({}, originalContact);
    console.log('contact contact:');
    console.log(this.contact);
  }

  ngOnInit() {
    this.getMemberContactRoles();
    this.getOrgContactRoles();
  }

  getDefaults() {
    // Instantiate member data
    this.org = {
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

  getOrgContactRoles() {
    // TODO: replace with call to cinco
    this.orgContactRoles = [
      {
        key: 'NONE',
        pretty_value: '',
      },
      {
        key: 'IT',
        pretty_value: 'IT',
      },
      {
        key: 'DIRECTOR',
        pretty_value: 'Director',
      },
      {
        key: 'EXECUTIVE',
        pretty_value: 'Executive',
      },
      {
        key: 'FINANCE',
        pretty_value: 'Finance',
      },
      {
        key: 'MANAGER_NO_SUBS',
        pretty_value: 'Manager (without subordinates)',
      },
      {
        key: 'MANAGER_WITH_SUBS',
        pretty_value: 'Manager (with subordinates)',
      },
      {
        key: 'OPERATIONS',
        pretty_value: 'Operations',
      },
      {
        key: 'OWNER_PARTNER',
        pretty_value: 'Owner/Partner',
      }
    ];
  }

  getMemberContactRoles() {
    // TODO: replace with call to cinco
    this.memberContactRoles = [
      {
        key: 'BILLING_CONTACT',
        pretty_value: 'Billing Contact',
      },
      {
        key: 'LEGAL_CONTACT',
        pretty_value: 'Legal Contact',
      },
      {
        key: 'PRESS_CONTACT',
        pretty_value: 'Press Contact',
      },
      {
        key: 'MARKETING_CONTACT',
        pretty_value: 'Marketing Contact',
      },
      {
        key: 'TECHNICAL_CONTACT',
        pretty_value: 'Technical Contact',
      },
      {
        key: 'REP_VOTING_CONTACT',
        pretty_value: 'Representative/Voting Contact',
      }
    ];
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
              this.contact.primaryContact = 'no';
            }
          },
          {
            text: 'Assign',
            handler: data => {
              this.contact.primaryContact = 'yes';
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
              this.contact.boardMember = 'no';
            }
          },
          {
            text: 'Assign',
            handler: data => {
              this.contact.boardMember = 'yes';
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
    console.log(this.memberId);
    console.log(this.contact);
    if (this.contact.primaryContact === 'yes') {
      this.contact.primaryContact = true;
    }
    else {
      this.contact.primaryContact = false;
    }

    if (this.contact.boardMember === 'yes') {
      this.contact.boardMember = true;
    }
    else {
      this.contact.boardMember = false;
    }

    if (this.contact.id) {
      // TODO: fix up roleId logic so it is based off the contact when first recieved and determines if it has been attached to a member.
      // if (this.roleId) {
      if (false) {
        this.cincoService.updateMemberContact(this.projectId, this.memberId, this.contact.id, this.contact).subscribe(response => {
          if(response) {
            console.log('updateMemberContact response:');
            console.log(response);
          }
        });
      }
      else {
        this.cincoService.addMemberContact(this.projectId, this.memberId, this.contact.id, this.contact).subscribe(response => {
          if (response) {
            console.log('addMemberContact response:');
            console.log(response);
          }
        });
      }
    }
    else {
      // Add new contact to organization
      this.cincoService.createOrganizationContact(this.org.id, this.contact).subscribe(response => {
        if (response) {
          console.log('createOrganizationContact response:');
          console.log(response);
          this.contact.id = response;
          // add to member
          this.cincoService.addMemberContact(this.projectId, this.memberId, this.contact.id, this.contact).subscribe(response => {
            if (response) {
              console.log('addMemberContact response:');
              console.log(response);
            }
          });
        }
      });  
    }

  }

  removeContact() {
    console.log('remove contact');
    console.log(this.contact);
    this.cincoService.removeMemberContact(this.projectId, this.memberId, this.contact.id).subscribe(response => {
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
