import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { MemberPage } from './member';

@NgModule({
  declarations: [
    MemberPage
  ],
  imports: [
    IonicPageModule.forChild(MemberPage)
  ],
  entryComponents: [
    MemberPage
  ]
})
export class MemberPageModule {}
