import { Component, ElementRef, Input, ViewChild, Renderer, } from '@angular/core';

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
  * The callback executed when files are selected, set by parent
  */
  @Input()
  private btnCallback: Function;

  /**
   * The text used for the upload
   */
  @Input()
  private callbackContext: any;

  /**
   * Native upload button (hidden)
   */
  @ViewChild('input')
  private nativeInputBtn: ElementRef;

  constructor(private renderer: Renderer) {
    console.log ("upload button loaded");
  }

  /**
  * Callback executed when the visible button is pressed
  * @param  {Event}  event should be a mouse click event
  */
  public callback(event: Event): void {

    // trigger click event of hidden input
    let clickEvent: MouseEvent = new MouseEvent("click", {bubbles: true});
    this.renderer.invokeElementMethod(
        this.nativeInputBtn.nativeElement, "dispatchEvent", [clickEvent]);
  }

  /**
  * Callback which is executed after files from native popup are selected.
  * @param  {Event}    event change event containing selected files
  */
  public filesAdded(event: Event): void {
    let files: FileList = this.nativeInputBtn.nativeElement.files;
    this.btnCallback(files, this.callbackContext);
  }

}
