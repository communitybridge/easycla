import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { WhitelistModal } from './whitelist-modal';

@NgModule({
  declarations: [
    WhitelistModal,
  ],
  imports: [
    IonicPageModule.forChild(WhitelistModal)
  ],
  entryComponents: [
    WhitelistModal,
  ]
})
export class WhitelistModalModule {}
