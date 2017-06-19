import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { AccountSettingsPage } from './account-settings';

@NgModule({
  declarations: [
    AccountSettingsPage
  ],
  imports: [
    IonicPageModule.forChild(AccountSettingsPage)
  ],
  entryComponents: [
    AccountSettingsPage
  ]
})
export class AccountSettingsPageModule {}
