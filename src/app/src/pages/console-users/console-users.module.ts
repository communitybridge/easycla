import { NgModule } from '@angular/core';

import { IonicPageModule } from 'ionic-angular';

import { ConsoleUsersPage } from './console-users';

@NgModule({
  declarations: [
    ConsoleUsersPage
  ],
  imports: [
    IonicPageModule.forChild(ConsoleUsersPage)
  ],
})
export class ConsoleUsersPageModule {}
