import { NgModule } from '@angular/core';

import { IonicPageModule } from 'ionic-angular';

import { AllMembersPage } from './all-members';

@NgModule({
  declarations: [
    AllMembersPage
  ],
  imports: [
    IonicPageModule.forChild(AllMembersPage)
  ],
})
export class AllMembersPageModule {}
