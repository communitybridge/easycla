import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'app-page-not-found',
  templateUrl: './page-not-found.component.html',
  styleUrls: ['./page-not-found.component.scss']
})
export class PageNotFoundComponent implements OnInit {
  message: string;
  constructor() { }

  ngOnInit(): void {
    this.message = 'The page you are looking for was not found.';
  }
}
