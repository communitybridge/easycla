import { Component, ChangeDetectorRef } from '@angular/core';

import { NavController, NavParams, ViewController } from 'ionic-angular';

import { AlertController } from 'ionic-angular';


@Component({
  selector: 'contact-update',
  templateUrl: 'contact-update.html'
})
export class ContactUpdate {
  project: any;
  member: any;
  contact: any;
  imagePreview: any;
  contactUpdate: any = this;

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
            this.AddEmailGroup(data.email);
          }
        }
      ]
    });
    prompt.present();
  }

  AddEmailGroup(groupName) {
    if(!this.contact.email_groups) {
      this.contact.email_groups = [];
    }
    this.contact.email_groups.push(
      {
        name: groupName
      }
    );
  }

  RemoveEmailGroup(index) {
    this.contact.email_groups.splice(index, 1);
  }

  SaveContact(contact) {

  }

  RemoveContact(contact) {

  }

  upload(files, context) {
    context.readFiles(files);
  }

  filesNotify(files) {
    console.log("files notify");
    console.log(files);
    console.log("this");
    console.log(this.member);
    this.readFiles(files);
  }

  readFiles(files, index=0){
    // Create the file reader
    let reader = new FileReader();
    // If there is a file
    if(index in files){
      // Start reading this file
      this.readFile(files[index], reader, (result) =>{
        this.imagePreview = result;
      });
    }else{
      // When all files are done This forces a change detection
      this.changeDetectorRef.detectChanges();
    }
  }

  readFile(file, reader, callback){
    reader.onload = () => {
      callback(reader.result);
    }

    reader.readAsDataURL(file);
  }

}
