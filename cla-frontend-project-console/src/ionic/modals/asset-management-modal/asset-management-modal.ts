// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component, Renderer, ElementRef, ViewChild, } from '@angular/core';
import { NavController, NavParams, ViewController, AlertController, ToastController, IonicPage, ActionSheetController  } from 'ionic-angular';
import { CincoService } from '../../services/cinco.service'
import { S3Service } from '../../services/s3.service'

@IonicPage({
  segment: 'asset-management-modal/:projectId'
})
@Component({
  selector: 'asset-management-modal',
  templateUrl: 'asset-management-modal.html',
})
export class AssetManagementModal {

  projectId: string; // Always Needed
  projectName: string;
  logoClassifier: string; //can be anything, but it's meant to be something like "main" or "black-and-white" or "thumbnail".
  documentClassifier: string;
  image: any;
  document: any;

  projectLogos: any;
  projectDocuments: any;

  files: any;
  folders: any;
  selectedFiles: any;
  selectedDirectory: any;
  loading: any;

  newLogoRef: string = "";

  uploadMode: string = "";
  /**
   * Comma separated array of allowed file extensions
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
    public actionSheetCtrl: ActionSheetController,
    public s3Service: S3Service,
  ) {
    this.projectId = navParams.get('projectId');
    this.projectName = navParams.get('projectName');
    this.selectedFiles = [];
    this.uploadTypes = 'jpg,jpeg,png,gif,tif,psd,ai,docx,pptx,pdf';
    this.uploadSizeMax = 50000000; // 50MB
  }

  ngOnInit() {

  }

  async ionViewDidEnter(){
    this.getAllProjectLogosDocuments();
    this.getDefaults();
    this.selectDirectory("root");
  }

  getAllProjectLogosDocuments(){
    this.cincoService.getProjectLogos(this.projectId).subscribe(response => {
      // CINCO sample response
      // classifier: "main"
      // key:"logos/project/a090n000000uQhdAAE/main.png"
      // publicUrl:"http://docker.for.mac.localhost:50563/public-media.platform.linuxfoundation.org/logos/project/a090n000000uQhdAAE/main.png
      if(response) {
        this.projectLogos = response;
        // Temporary Fix until CINCO returns filename in GET logos response
        for(let eachLogo of this.projectLogos) {
          eachLogo.key = eachLogo.key.split("/");
          eachLogo.key = eachLogo.key[3];
        }
      }
    });

    this.cincoService.getProjectDocuments(this.projectId).subscribe(response => {
      // CINCO sample response
      // classifier:"minutes"
      // expiresOn:"2017-10-10T16:02:02.654Z"
      // name:"samplecontract.pdf"
      // url:"http://docker.for.mac.localhost:51409/private-media.platform.linuxfoundation.org/documents/project/a093F0000001YReQAM/minutes/samplecontract.pdf?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Date=20171010T150202Z&X-Amz-SignedHeaders=host&X-Amz-Expires=3600&X-Amz-Credential=NDB7TUF7W7FAQ7Z0U64V%2F20171010%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Signature=9ef3d41d712766c493ba07b535cec189579e87b7966dc88f8ffeefea0212c383"
      if(response){
        this.projectDocuments = response;
      }
    });
  }

  getDefaults() {
    this.files = [
      // {
      //   id: 'A000000001',
      //   name: 'Zephyr_Bylaws.pdf',
      //   type: 'file',
      //   lastUpdated: '3/3/2017',
      //   classifier: 'bylaws'
      // },
      // {
      //   id: 'A000000002',
      //   name: 'Zephyr_LF_membership_agreement.pdf',
      //   type: 'file',
      //   lastUpdated: '3/3/2017',
      //   classifier: 'bylaws'
      // },
    ];
    this.folders = [
      {
        displayName: 'Logos',
        type: 'folder',
        value: 'logos'
      },
      {
        displayName: 'Documents',
        type: 'folder',
        value: 'documents'
      }
    ];
  }

  // ContactUpdateModal modal dismiss
  dismiss() {
    this.viewCtrl.dismiss(this.newLogoRef);
  }

  selectDirectory(directory){
    this.selectedDirectory = directory;
  }

  selectFile(event, file) {
    event.stopPropagation();

    if (event.ctrlKey) {
      if (file.selected) {
        this.deselectFiles([file]);
      }
      else {
        file.selected = true;
        this.selectedFiles.push(file);
      }
    }
    else { // standard single file select
      // deselect the entire selected files array
      this.deselectFiles(this.selectedFiles);
      file.selected = true;
      this.selectedFiles = [file];
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
    if(this.selectedFiles.length > 1) {
      prompt_title = 'Delete files?';
    }
    else {
      prompt_title = 'Delete file?';
    }
    let prompt_message = this.selectedFiles.map(function(file){
      return file.name;
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
            // TODO: Make cinco calls to delete files
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

  deselectFiles(files) {
    if (!files) {
      return;
    }
    let i = files.length;
    while (i--) {
      let file = files[i];
      file.selected = false;
      let index = this.selectedFiles.indexOf(file);
      if(index !== -1) {
        this.selectedFiles.splice(index, 1);
      }
    }
  }

  modalClick(event) {
    // stray unhandled/unprevented click. deselect all files
    this.deselectFiles(this.selectedFiles);
  }

  /**
   * Callback executed when the visible button is pressed
   * @param  {Event}  event should be a mouse click event
   */
  uploadClicked(event: Event) {
    this.presentActionSheet();
  }

  presentActionSheet() {
    let actionSheet = this.actionSheetCtrl.create({
      title: 'I want to upload a',
      buttons: [
        {
          text: 'Logo',
          role: 'logo',
          handler: () => {
            this.uploadMode = "logo";
            this.showLogoClassifier();
          }
        },{
          text: 'Document',
          role: 'document',
          handler: () => {
            this.uploadMode = "document";
            this.showDocumentClassifier();
          }
        },{
          text: 'Cancel',
          role: 'cancel',
          handler: () => {
            this.uploadMode = "";
          }
        }
      ]
    });
    actionSheet.present();
  }

  showLogoClassifier() {
    let alert = this.alertCtrl.create();
    alert.setTitle('Choose Logo Category');

    alert.addInput({
      type: 'radio',
      label: 'Main Logo',
      value: 'main',
      checked: true
    });

    alert.addInput({
      type: 'radio',
      label: 'Black and White Logo',
      value: 'black-and-white',
      checked: false
    });

    alert.addInput({
      type: 'radio',
      label: 'Thumbnail Logo',
      value: 'thumbnail',
      checked: false
    });

    alert.addButton('Cancel');
    alert.addButton({
      text: 'OK',
      handler: data => {
        this.logoClassifier = data;
        // trigger click event of hidden input
        let clickEvent: MouseEvent = new MouseEvent("click", {bubbles: true});
        this.renderer.invokeElementMethod(
          this.nativeInputBtn.nativeElement, "dispatchEvent", [clickEvent]
        );

      }
    });
    alert.present();
  }


  showDocumentClassifier() {
    let alert = this.alertCtrl.create();
    alert.setTitle('Choose Document Category');

    alert.addInput({
      type: 'radio',
      label: 'Minutes',
      value: 'minutes',
      checked: true
    });

    alert.addInput({
      type: 'radio',
      label: 'Bylaws',
      value: 'bylaws',
      checked: false
    });

    alert.addButton('Cancel');
    alert.addButton({
      text: 'OK',
      handler: data => {
        this.documentClassifier = data;
        // trigger click event of hidden input
        let clickEvent: MouseEvent = new MouseEvent("click", {bubbles: true});
        this.renderer.invokeElementMethod(
          this.nativeInputBtn.nativeElement, "dispatchEvent", [clickEvent]
        );

      }
    });
    alert.present();
  }

  /**
   * Callback which is executed after files from native popup are selected.
   * @param  {Event}    event change event containing selected files
   */
  filesAdded(event: Event, uploadMode: string): void {

    let addedFiles: FileList = this.nativeInputBtn.nativeElement.files;

    for(let i=0; i< addedFiles.length; i++) {
      let file = addedFiles.item(i);
      let contentType;


      let valid = this.validateFile(file);
      if(valid) {

        contentType = file.type;

        let reader = new FileReader();
        reader.onload = event => {
          if(uploadMode == 'logo') {
            const imgBase64 = (<any>event).target.result;
            this.image = {
              imageBytes: imgBase64,
              contentType: contentType
            }
            this.cincoService.obtainLogoS3URL(this.projectId, this.logoClassifier, this.image).subscribe(response => {
              if(response.putUrl.url) {
                let S3URL = response.putUrl.url;
                this.s3Service.uploadToS3(S3URL, file, this.image.contentType).subscribe(response => {
                  // This is to refresh an Main Logo image that have same URL as its previous upload. E.g. main.png.
                  if(this.logoClassifier == 'main') {
                    this.newLogoRef = response.url.split("?", 1)[0] + "?" + new Date().getTime();
                    this.dismiss();
                  }
                  else {
                    this.getAllProjectLogosDocuments();
                    this.selectDirectory('logos');
                  }
                });

              }
            });
          }
          else if (uploadMode == 'document') {
            const docBytes = (<any>event).target.result;
            this.cincoService.obtainDocumentS3URL(this.projectId, this.documentClassifier, docBytes, file.name, contentType).subscribe(response => {
              if(response.putUrl.url) {
                let S3URL = response.putUrl.url;
                this.s3Service.uploadToS3(S3URL, file, file.type).subscribe(response => {
                  this.getAllProjectLogosDocuments();
                  this.selectDirectory('documents');
                });
              }
            });
          }
          else {
            this.dismiss();
          }
        }

        reader.readAsDataURL(file);

        // merge files from the input with the files
        /*
          TODO: send a call to cinco with the new file data
          from the response, add the file to the files array

          if(!this.files) {
            this.files = [];
          }
          this.files.push(fileResponse);

         */
      }
    }
  }

  validateFile(file) {
    var extensionValid = false;
    if(typeof this.uploadTypes == 'undefined') {
      extensionValid = true;
    }
    else {
      // Validate extension by checking extension in filename against uploadTypes
      var validTypes = this.uploadTypes.split(',');

      for (var i = 0; i < validTypes.length; i++) {
        var currentType = validTypes[i];
        if (file.name.substr(file.name.length - currentType.length, currentType.length).toLowerCase() == currentType.toLowerCase()) {
          extensionValid = true;
        }
      }
    }

    if (!extensionValid) {
      this.uploadError("Sorry, " + file.name + " is invalid, allowed extensions are: " + validTypes.join(", "));
      return false;
    }

    var sizeValid = false;
    if (typeof this.uploadSizeMax == 'undefined') {
      sizeValid = true;
    }
    else {
      if (file.size < this.uploadSizeMax) {
        sizeValid = true;
      }
    }

    if (!sizeValid) {
      let maxSize = this.uploadSizeMax + 'bytes';
      this.uploadError("Sorry, " + file.name + " is too big, max size is: " + maxSize);
      return false;
    }

    // All individual checks/returns should have happened by now
    return true;

  }

  uploadError(message) {
    let toast = this.toastCtrl.create({
      message: message,
      showCloseButton: true,
      closeButtonText: 'Dismiss',
      duration: 3000
    });
    toast.present();
  }

}

