import { IonicModule } from 'ionic-angular';
import { UploadButtonComponent } from './upload-button/upload-button';
import { NgModule } from '@angular/core';

@NgModule({
  declarations: [
    UploadButtonComponent
  ],
  imports: [
    IonicModule
  ],
  exports: [
    UploadButtonComponent
  ]
})
export class ComponentsModule {}
