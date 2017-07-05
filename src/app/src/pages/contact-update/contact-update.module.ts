import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { ContactUpdate } from './contact-update';
import { LoadingSpinnerComponentModule } from '../../components/loading-spinner/loading-spinner.module';
import { LoadingDisplayDirectiveModule } from '../../directives/loading-display/loading-display.module';
import { UploadButtonComponentModule } from '../../components/upload-button/upload-button.module';

@NgModule({
  declarations: [
    ContactUpdate,
  ],
  imports: [
    LoadingSpinnerComponentModule,
    LoadingDisplayDirectiveModule,
    UploadButtonComponentModule,
    IonicPageModule.forChild(ContactUpdate),
  ],
  entryComponents: [
    ContactUpdate,
  ]
})
export class ContactUpdateModule {}
