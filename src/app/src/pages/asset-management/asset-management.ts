import { Component, Renderer, ElementRef, ViewChild, } from '@angular/core';
import { NavController, NavParams, ViewController, AlertController, ToastController, IonicPage  } from 'ionic-angular';
import { CincoService } from '../../app/services/cinco.service'

@IonicPage({
  segment: 'asset-management'
})
@Component({
  selector: 'asset-management',
  templateUrl: 'asset-management.html',
  providers: [CincoService]
})
export class AssetManagementModal {
  projectId: string; // Always Needed
  files: any;
  folders: any;
  selectedFiles: any;

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
  ) {
    this.selectedFiles = [];
    this.uploadTypes = 'jpg,jpeg,png,gif,tif,psd,ai,docx,pptx,pdf';
    this.uploadSizeMax = 50000000; // 50MB
    this.getDefaults();
  }

  ngOnInit() {

  }

  getDefaults() {
    this.files = [
      {
        id: 'A000000001',
        name: 'Zephyr_Bylaws.pdf',
        type: 'file',
        lastUpdated: '3/3/2017',
        notes: ''
      },
      {
        id: 'A000000002',
        name: 'Zephyr_LF_membership_agreement.pdf',
        type: 'file',
        lastUpdated: '3/3/2017',
        notes: 'Linux Foundation membership agreement'
      },
      {
        id: 'A000000003',
        name: 'Zephyr_project_membership_agreement.pdf',
        type: 'file',
        lastUpdated: '3/3/2017',
        notes: 'Project membership agreement, updated on March 2nd.'
      },
      {
        id: 'A000000004',
        name: 'Zephyr_sow.pdf',
        type: 'file',
        lastUpdated: '3/3/2017',
        notes: ''
      },
      {
        id: 'A000000005',
        name: 'Technical_steering_ctr.pdf',
        type: 'file',
        lastUpdated: '3/3/2017',
        notes: 'Technical steering charter, last updated March 1st'
      },
      {
        id: 'A000000006',
        name: 'Zephyr_Bylaws.pdf',
        type: 'file',
        lastUpdated: '3/3/2017',
        notes: ''
      },
    ];
    this.folders = [
      {
        name: 'Logos',
        type: 'folder',
      },
    ];
  }

  // ContactUpdate modal dismiss
  dismiss() {
    this.viewCtrl.dismiss();
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

    // trigger click event of hidden input
    let clickEvent: MouseEvent = new MouseEvent("click", {bubbles: true});
    this.renderer.invokeElementMethod(
      this.nativeInputBtn.nativeElement, "dispatchEvent", [clickEvent]
    );
  }

  /**
  * Callback which is executed after files from native popup are selected.
  * @param  {Event}    event change event containing selected files
  */
  filesAdded(event: Event): void {
    let addedFiles: FileList = this.nativeInputBtn.nativeElement.files;

    for(let i=0; i< addedFiles.length; i++) {
      let file = addedFiles.item(i);
      let valid = this.validateFile(file);
      if(valid) {
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
