import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { AuthorityYesnoPage } from './authority-yesno-page';
import { LayoutModule } from "../../layout/layout.module";

@NgModule({
  declarations: [
    AuthorityYesnoPage,
  ],
  imports: [
    IonicPageModule.forChild(AuthorityYesnoPage),
    LayoutModule
  ],
  entryComponents: [
    AuthorityYesnoPage
  ]
})
export class AuthorityYesnoPageModule {}
