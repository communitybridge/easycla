import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { AuthorityYesnoPage } from './authority-yesno-page';

@NgModule({
  declarations: [
    AuthorityYesnoPage,
  ],
  imports: [
    IonicPageModule.forChild(AuthorityYesnoPage),
  ],
  entryComponents: [
    AuthorityYesnoPage
  ]
})
export class AuthorityYesnoPageModule {}
