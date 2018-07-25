import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { ClaMessageSentPage } from './cla-message-sent';

@NgModule({
  declarations: [
    ClaMessageSentPage,
  ],
  imports: [
    IonicPageModule.forChild(ClaMessageSentPage),
  ],
  entryComponents: [
    ClaMessageSentPage
  ]
})
export class ClaMessageSentPageModule {}
