import { Component, ChangeDetectorRef, ViewChild } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { NavController, NavParams, ViewController, AlertController, IonicPage, Content } from 'ionic-angular';
import { PhoneNumberValidator } from  '../../validators/phonenumber';
import { EmailValidator } from  '../../validators/email';
import { CincoService } from '../../app/services/cinco.service';

@IonicPage({
  segment: 'console-user-update'
})
@Component({
  selector: 'console-user-update-modal',
  templateUrl: 'console-user-update-modal.html',
  providers: [CincoService]
})
export class ConsoleUserUpdateModal {
  user: any;
  userRoles: any;
  keysGetter;
  _form: FormGroup;
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
    console.log("user:");
    console.log(this.user);

    this._form = formBuilder.group({
      userId:[this.user.userId, Validators.required],
      email:[this.user.email, Validators.compose([Validators.required, EmailValidator.isValid])],
      // 
      // boardMember:[this.contact.boardMember, Validators.required],
      // bio:[this.contact.contact.bio, Validators.required],
    });

  }

  ngOnInit() {
    this.getUserRoles();
  }

  getDefaults() {
    // Instantiate user data
    this.user = {
      userId: '',
      email: '',
      roles: [],
    };

    this.userRoles = {};
  }

  getUserRoles() {
    // this.cincoService.getUserRoles().subscribe(response => {
    //   if(response) {
    //     this.userRoles = response;
    //   }
    // });
  }

  // ContactUpdate modal dismiss
  dismiss() {
    this.viewCtrl.dismiss();
  }

  primarySelectChanged(value) {
    // // normalize the value from string to bool
    // if (value == 'true') {
    //   let prompt = this.alertCtrl.create({
    //     title: 'Assign as Primary?',
    //     message: 'This will replace the exising primary contact on this project.',
    //     buttons: [
    //       {
    //         text: 'Cancel',
    //         handler: data => {
    //           this.contactUpdateForm.patchValue({primaryContact:'false'});
    //         }
    //       },
    //       {
    //         text: 'Assign',
    //         handler: data => {
    //           this.contactUpdateForm.patchValue({primaryContact:'true'});
    //         }
    //       }
    //     ]
    //   });
    //   prompt.present();
    // }
  }

  boardSelectChanged(value) {
    // // normalize the value from string to bool
    // if (value == 'true') {
    //   let prompt = this.alertCtrl.create({
    //     title: 'Assign to Board?',
    //     message: 'This will replace the exising member company Board member on this project.',
    //     buttons: [
    //       {
    //         text: 'Cancel',
    //         handler: data => {
    //           this.contactUpdateForm.patchValue({boardMember:'false'});
    //         }
    //       },
    //       {
    //         text: 'Assign',
    //         handler: data => {
    //           this.contactUpdateForm.patchValue({boardMember:'true'});
    //         }
    //       }
    //     ]
    //   });
    //   prompt.present();
    // }
  }

  removeContactPrompt() {
    // let prompt = this.alertCtrl.create({
    //   title: 'Remove contact?',
    //   message: 'This will remove this contact from the project and remove them from any email alias or list.',
    //   buttons: [
    //     {
    //       text: 'Cancel',
    //       handler: data => {
    //         // Do nothing
    //       }
    //     },
    //     {
    //       text: 'Remove',
    //       handler: data => {
    //         this.removeContact();
    //       }
    //     }
    //   ]
    // });
    // prompt.present();
  }


  showPrompt() {
    // let prompt = this.alertCtrl.create({
    //   title: 'Add Email Group',
    //   message: '',
    //   inputs: [
    //     {
    //       name: 'email',
    //       placeholder: 'Title'
    //     },
    //   ],
    //   buttons: [
    //     {
    //       text: 'Cancel',
    //       handler: data => {
    //         // Don't add contact
    //       }
    //     },
    //     {
    //       text: 'Save',
    //       handler: data => {
    //         this.addEmailGroup(data.email);
    //       }
    //     }
    //   ]
    // });
    // prompt.present();
  }


  saveContact() {
    this.submitAttempt = true;
    this.currentlySubmitting = true;
    // if (!this.contactUpdateForm.valid){
    //   this.content.scrollToTop();
    //   this.currentlySubmitting = false;
    //   // prevent submit
    //   return;
    // }
    // let primaryContact = this.contactUpdateForm.value.primaryContact;
    // primaryContact = (primaryContact === true || primaryContact === 'true')
    //       ? true
    //       : false;
    // let boardMember = this.contactUpdateForm.value.boardMember;
    // boardMember = (boardMember === true || boardMember === 'true')
    //       ? true
    //       : false;
    // var memberContact = {
    //   id: this.contact.id,
    //   memberId: this.contact.memberId,
    //   type: this.contactUpdateForm.value.role,
    //   primaryContact: primaryContact,
    //   boardMember: boardMember,
    //   contact: {
    //     id: this.contact.contact.id,
    //     email: this.contactUpdateForm.value.email,
    //     givenName: this.contactUpdateForm.value.givenName,
    //     familyName: this.contactUpdateForm.value.familyName,
    //     title: this.contactUpdateForm.value.title,
    //     phone: this.contactUpdateForm.value.phone,
    //     type: this.contactUpdateForm.value.type,
    //     bio: this.contactUpdateForm.value.bio,
    //   },
    // };
    // if (this.contactId) {
    //   if (this.roleId) {
    //     this.cincoService.updateOrganizationContact(this.org.id, this.contactId, memberContact.contact).subscribe(response => {
    //       if (response) {
    //         // update org contact with response from update
    //         // should be the same as what was sent, but we will just be sure
    //         memberContact.contact = response;
    //         // add as a member contact
    //         this.cincoService.updateMemberContact(this.projectId, this.memberId, this.contactId, this.roleId, memberContact).subscribe(response => {
    //           if(response) {
    //             this.dismiss();
    //           }
    //         });
    //       }
    //     });
    //   }
    //   else {
    //     this.cincoService.updateOrganizationContact(this.org.id, this.contactId, memberContact.contact).subscribe(response => {
    //       if (response) {
    //         // update org contact with response from update
    //         // should be the same as what was sent, but we will just be sure
    //         memberContact.contact = response;
    //         // add as a member contact
    //         this.cincoService.addMemberContact(this.projectId, this.memberId, this.contactId, memberContact).subscribe(response => {
    //           if (response) {
    //             this.dismiss();
    //           }
    //         });
    //       }
    //     });
    //   }
    // }
    // else {
    //   // Add new contact to organization
    //   this.cincoService.createOrganizationContact(this.org.id, memberContact.contact).subscribe(response => {
    //     if (response) {
    //       this.contactId = response;
    //       memberContact.contact.id = this.contactId;
    //       // add to member
    //       this.cincoService.addMemberContact(this.projectId, this.memberId, this.contactId, memberContact).subscribe(response => {
    //         if (response) {
    //           this.dismiss();
    //         }
    //       });
    //     }
    //   });
    // }

  }

  removeContact() {
    this.currentlySubmitting = true;
    // this.cincoService.removeMemberContact(this.projectId, this.memberId, this.contactId, this.roleId).subscribe(response => {
    //   if(response) {
    //     this.dismiss();
    //   }
    // });
  }

}
