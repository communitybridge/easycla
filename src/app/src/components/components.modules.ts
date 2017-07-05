import { IonicModule } from 'ionic-angular';
// import { ActionPopoverComponent } from './action-popover/action-popover';
import { UploadButtonComponent } from './upload-button/upload-button';
// import { LoadingSpinnerComponent } from './loading-spinner/loading-spinner';
import { NgModule } from '@angular/core';

@NgModule({
  declarations: [
    // ActionPopoverComponent,
    UploadButtonComponent,
    // LoadingSpinnerComponent
  ],
  imports: [
    IonicModule
  ],
  exports: [
    // ActionPopoverComponent,
    UploadButtonComponent,
    // LoadingSpinnerComponent,
  ]
})
export class ComponentsModule {}
