import { Input, Component } from '@angular/core';

@Component({
  selector: 'sorting-display',
  templateUrl: 'sorting-display.html'
})
export class SortingDisplayComponent {

  /**
   * The text used for the upload label
   */
  @Input('sorting')
  private sorting: string;

  constructor() {
    this.sorting;
  }

}
