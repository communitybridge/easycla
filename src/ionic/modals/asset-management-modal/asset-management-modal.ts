import { Component, Renderer, ElementRef, ViewChild, } from '@angular/core';
import { NavController, NavParams, ViewController, AlertController, ToastController, IonicPage, ActionSheetController  } from 'ionic-angular';
import { CincoService } from '../../services/cinco.service'

@IonicPage({
  segment: 'asset-management-modal/:projectId'
})
@Component({
  selector: 'asset-management-modal',
  templateUrl: 'asset-management-modal.html',
  providers: [CincoService]
})
export class AssetManagementModal {

  projectId: string; // Always Needed
  logoClassifier: string; //can be anything, but it's meant to be something like "main" or "black-and-white" or "thumbnail".
  documentClassifier: string;
  image: any;
  document: any;

  files: any;
  folders: any;
  selectedFiles: any;
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
    public actionSheetCtrl: ActionSheetController
  ) {
    this.projectId = navParams.get('projectId');
    this.selectedFiles = [];
    this.uploadTypes = 'jpg,jpeg,png,gif,tif,psd,ai,docx,pptx,pdf';
    this.uploadSizeMax = 50000000; // 50MB
    this.getDefaults();
  }

  ngOnInit() {

  }

  ionViewDidEnter(){

    this.cincoService.getProjectLogos(this.projectId).subscribe(response => {

      // CINCO sample response
      // logoClassifier: "main"
      // key:"logos/project/a090n000000uQhdAAE/main.png"
      // publicUrl:"http://docker.for.mac.localhost:50563/public-media.platform.linuxfoundation.org/logos/project/a090n000000uQhdAAE/main.png

      console.log("getProjectLogos");
      console.log(response);
    });

    this.cincoService.getProjectDocuments(this.projectId).subscribe(response => {
      console.log("getProjectDocuments");
      console.log(response);
    });

  }

  getDefaults() {
    this.files = [
      // {
      //   id: 'A000000001',
      //   name: 'Zephyr_Bylaws.pdf',
      //   type: 'file',
      //   lastUpdated: '3/3/2017',
      //   notes: ''
      // },
      // {
      //   id: 'A000000002',
      //   name: 'Zephyr_LF_membership_agreement.pdf',
      //   type: 'file',
      //   lastUpdated: '3/3/2017',
      //   notes: 'Linux Foundation membership agreement'
      // },
    ];
    this.folders = [
      {
        name: 'Logos',
        type: 'folder',
      },
    ];
  }

  // ContactUpdateModal modal dismiss
  dismiss() {
    this.viewCtrl.dismiss(this.newLogoRef);
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
    // trigger click event of hidden input
    // let clickEvent: MouseEvent = new MouseEvent("click", {bubbles: true});
    // this.renderer.invokeElementMethod(
    //   this.nativeInputBtn.nativeElement, "dispatchEvent", [clickEvent]
    // );
  }

  presentActionSheet() {
    let actionSheet = this.actionSheetCtrl.create({
      title: 'I want to upload a',
      buttons: [
        {
        text: 'Logo',
        role: 'logo',
        handler: () => {
          console.log('Logo clicked');
          this.uploadMode = "logo";
          // trigger click event of hidden input
          let clickEvent: MouseEvent = new MouseEvent("click", {bubbles: true});
          this.renderer.invokeElementMethod(
            this.nativeInputBtn.nativeElement, "dispatchEvent", [clickEvent]
          );
        }
      },{
        text: 'Document',
        handler: () => {
          console.log('Upload clicked');
          this.uploadMode = "document";
          // trigger click event of hidden input
          let clickEvent: MouseEvent = new MouseEvent("click", {bubbles: true});
          this.renderer.invokeElementMethod(
            this.nativeInputBtn.nativeElement, "dispatchEvent", [clickEvent]
          );
        }
      },{
        text: 'Cancel',
        role: 'cancel',
        handler: () => {
          this.uploadMode = "";
          console.log('Cancel clicked');
        }
      }
    ]
  });
  actionSheet.present();
}

  /**
  * Callback which is executed after files from native popup are selected.
  * @param  {Event}    event change event containing selected files
  */
  filesAdded(event: Event, uploadMode: string): void {
    console.log("uploadMode: ", uploadMode);
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
            this.logoClassifier = "main";
            const imgBase64 = (<any>event).target.result;
            this.image = {
              imageBytes: imgBase64,
              contentType: contentType
            }
            console.log("test!");
            this.cincoService.obtainLogoS3URL(this.projectId, this.logoClassifier, this.image).subscribe(response => {
              console.log("obtainLogoS3URL response");
              console.log(response);
              if(response.url) {
                let S3URL = response.url;
                console.log("S3URL");
                console.log(S3URL);

                this.cincoService.uploadToS3(S3URL, file, this.image.contentType).subscribe(response => {
                   // This is to refresh an image that have same URL
                   this.newLogoRef = response.url.split("?", 1)[0] + "?" + new Date().getTime();
                   this.dismiss();
                 });

              }
            });
          }
          else if (uploadMode == 'document') {
            this.documentClassifier = 'minutes';

            const docBytes = (<any>event).target.result;

            console.log("this.documentClassifier: ", this.documentClassifier);
            contentType = file.type;
            console.log("file");
            console.log(file);
            console.log("file.name");
            console.log(file.name);
            console.log("contentType");
            console.log(contentType);
            console.log("docBytes");
            console.log(docBytes);
            this.cincoService.obtainDocumentS3URL(this.projectId, this.documentClassifier, docBytes, file.name, contentType).subscribe(response => {
              console.log("obtainDocumentS3URL");
              console.log(response);
              if(response.url) {
                let S3URL = response.url;
                console.log("S3URL");
                console.log(S3URL);
                console.log("file");
                console.log(file);
                console.log("contentType");
                console.log(contentType);
                this.dismiss();
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
