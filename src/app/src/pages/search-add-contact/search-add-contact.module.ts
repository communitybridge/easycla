import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { SearchAddContact } from './search-add-contact';

@NgModule({
  declarations: [
    SearchAddContact,
  ],
  imports: [
    IonicPageModule.forChild(SearchAddContact)
  ],
  entryComponents: [
    SearchAddContact,
  ]
})
export class SearchAddContactModule {}
