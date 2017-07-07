import { Input, Component } from '@angular/core';

@Component({
  selector: 'loading-spinner',
  templateUrl: 'loading-spinner.html'
})
export class LoadingSpinnerComponent {

  /**
   * The text used for the upload label
   */
  @Input('loading')
  private loading: boolean;

  constructor() {
  }

}
