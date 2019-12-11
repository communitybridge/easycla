// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component, ElementRef, Input, Output, ViewChild, Renderer, EventEmitter } from '@angular/core';
import { PopoverController, ToastController } from 'ionic-angular';

/*
  Generated class for the UploadButton component.

  See https://angular.io/docs/ts/latest/api/core/index/ComponentMetadata-class.html
  for more info on Angular 2 Components.
*/
@Component({
  selector: 'upload-button',
  templateUrl: 'upload-button.html'
})
export class UploadButtonComponent {
  /**
   * The fileList maintained as files are added and removed
   */
  fileList: Array<any>;

  /**
   * The text used for the upload label
   */
  @Input()
  private files: Array<any>;

  /**
   * The text used for the upload label
   */
  @Input()
  private uploadText: String;

  /**
   * Comma separated array of allowed file extensions
   */
  @Input()
  private uploadTypes: String;

  /**
   * If multiple files can be uploaded
   * @TODO: implement this feature so it affects the input multiple attribute
   *        and also build in validation for number of files allowed (1 or many)
   */
  @Input()
  private multiple: boolean;

  /**
   * Native upload button (hidden)
   */
  @ViewChild('input')
  private nativeInputBtn: ElementRef;

  /**
   * Event emitted after FileList is modified
   */
  @Output() notify: EventEmitter<Array<any>> = new EventEmitter<Array<any>>();

  constructor(private renderer: Renderer, private popoverCtrl: PopoverController, public toastCtrl: ToastController) {}

  ngOnInit() {
    this.fileList = this.files;

    if (typeof this.multiple == 'undefined') {
      this.multiple = false;
    }

    if (!this.uploadText) {
      this.uploadText = 'Upload';
    }
  }

  /**
   * Callback executed when the visible button is pressed
   * @param  {Event}  event should be a mouse click event
   */
  public callback(event: Event): void {
    // trigger click event of hidden input
    let clickEvent: MouseEvent = new MouseEvent('click', { bubbles: true });
    this.renderer.invokeElementMethod(this.nativeInputBtn.nativeElement, 'dispatchEvent', [clickEvent]);
  }

  /**
   * Callback which is executed after files from native popup are selected.
   * @param  {Event}    event change event containing selected files
   */
  filesAdded(event: Event): void {
    let addedFiles: FileList = this.nativeInputBtn.nativeElement.files;

    for (let i = 0; i < addedFiles.length; i++) {
      let file = addedFiles.item(i);
      let valid = this.validateFile(file);
      if (valid) {
        // merge files from the input with the fileList
        if (!this.fileList) {
          this.fileList = [];
        }
        this.fileList.push(file);
      }
    }
    this.notify.emit(this.fileList);
  }

  validateFile(file) {
    if (typeof this.uploadTypes == 'undefined') {
      return true;
    }
    // Validate extension by checking extension in filename against uploadTypes
    var validTypes = this.uploadTypes.split(',');
    var extensionValid = false;
    for (var i = 0; i < validTypes.length; i++) {
      var currentType = validTypes[i];
      if (
        file.name.substr(file.name.length - currentType.length, currentType.length).toLowerCase() ==
        currentType.toLowerCase()
      ) {
        extensionValid = true;
        return extensionValid;
      }
    }
    if (!extensionValid) {
      this.uploadError('Sorry, ' + file.name + ' is invalid, allowed extensions are: ' + validTypes.join(', '));
      return false;
    }
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

  presentPopover(ev, index) {
    let popoverData = {
      items: [
        {
          label: 'Download file',
          callback: 'fileDownload',
          callbackData: {
            index: index
          }
        },
        {
          label: 'Delete file',
          callback: 'fileDelete',
          callbackData: {
            index: index
          }
        }
      ]
    };

    let popover = this.popoverCtrl.create('ActionPopoverComponent', popoverData);

    popover.present({
      ev: ev
    });

    popover.onDidDismiss(popoverData => {
      if (popoverData) {
        this.popoverResponse(popoverData);
      }
    });
  }

  /**
   * Called if popover dismissed with data. Passes data to a callback function
   * @param  {object} popoverData should contain .callback and .callbackData
   */
  popoverResponse(popoverData) {
    let callback = popoverData.callback;
    if (this[callback]) {
      this[callback](popoverData.callbackData);
    }
  }

  fileDelete(data) {
    this.fileList.splice(data.index, 1);
    this.notify.emit(this.fileList);
  }

  fileDownload(data) {
    // File download:
    // data.index
    // this.fileList
  }
}
