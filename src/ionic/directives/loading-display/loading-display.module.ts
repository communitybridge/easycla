import { NgModule } from '@angular/core';
import { LoadingDisplayDirective  } from './loading-display';
@NgModule({
  declarations: [
    LoadingDisplayDirective
  ],
  exports: [
    LoadingDisplayDirective
  ]
})
export class LoadingDisplayDirectiveModule { }
