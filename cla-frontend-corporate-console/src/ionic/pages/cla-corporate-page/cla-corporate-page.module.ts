import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { ClaCorporatePage } from './cla-corporate-page';

@NgModule({
  declarations: [
    ClaCorporatePage,
  ],
  imports: [
    IonicPageModule.forChild(ClaCorporatePage),
  ],
  entryComponents: [
    ClaCorporatePage
  ]
})
export class ClaCorporatePageModule {}
