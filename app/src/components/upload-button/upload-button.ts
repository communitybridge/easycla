import { Component, ElementRef, Input, Output, ViewChild, Renderer, EventEmitter } from '@angular/core';

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
  uploadPreview: string;
  uploadTitle: string;

  private displayableTypes: Array<string>;
  /**
   * The icon used for the button
   */
  @Input()
  private btnIcon: String;

  /**
   * The text used for the upload
   */
  @Input()
  private btnText: String;


  /**
   * Native upload button (hidden)
   */
  @ViewChild('input')
  private nativeInputBtn: ElementRef;

  @Output() notify: EventEmitter<FileList> = new EventEmitter<FileList>();

  constructor(private renderer: Renderer) {
    console.log ("upload button loaded");
    this.displayableTypes = [
      'image/png',
      'image/jpeg',
      'image/gif',
    ];
  }

  /**
  * Callback executed when the visible button is pressed
  * @param  {Event}  event should be a mouse click event
  */
  public callback(event: Event): void {

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
  public filesAdded(event: Event): void {
    let files: FileList = this.nativeInputBtn.nativeElement.files;
    this.displayFile(files[0]);
    console.log("filesAdded:");
    console.log(files);
    this.notify.emit(files);
  }

  displayFile(file: File) {
    this.uploadTitle = file.name;
    let reader = new FileReader();
    if (this.displayableTypes.indexOf(file.type)!=-1) {
      this.readFile(file, reader, (result) =>{
        this.uploadPreview = result;
      });
    }
    else {
      this.uploadPreview = '';
    }
  }

  readFile(file, reader, callback){
    reader.onload = () => {
      callback(reader.result);
    }
    reader.readAsDataURL(file);
  }

}
