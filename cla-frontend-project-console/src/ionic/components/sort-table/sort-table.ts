import { Component, Input } from '@angular/core';

@Component({
  selector: 'easy-sort-table',
  templateUrl: 'sort-table.html'
})
export class SortTableComponent {

  @Input() column = [];
  @Input() columnData = [];
  @Input() data: any;
  @Input() childTitle: string;
  @Input() childTable = false;
  @Input() columnDataChildKey;
  @Input() childColumn = [];
  toggle = false;
  currentColumn: string;
  showIndex: any;
  


  constructor() {}


  ngOnInit(){
    this.sort(this.column[0].dataKey);
  }
  
  compareValues(key, order = 'asc') {
    return function innerSort(a, b) {
      if (!a.hasOwnProperty(key) || !b.hasOwnProperty(key)) {
        return 0;
      }
  
      const varA = (typeof a[key] === 'string')
        ? a[key].toLowerCase().trim() : a[key];
      const varB = (typeof b[key] === 'string')
        ? b[key].toLowerCase().trim()  : b[key];
  
      let comparison = 0;
      if (varA > varB) {
        comparison = 1;
      } else if (varA < varB) {
        comparison = -1;
      }
      return (
        (order === 'desc') ? (comparison * -1) : comparison
      );
    };
  }

  sort(col) {
    this.currentColumn = col.head;
    this.toggle = !this.toggle;
    this.currentColumn = col.head;
    const sortIn = this.toggle ? 'asc' : 'desc'
    this.columnData.sort(this.compareValues(col.dataKey, sortIn))
  }

  show(i) {
    return this.showIndex = this.showIndex === i ? '' : i
  }

  
}
