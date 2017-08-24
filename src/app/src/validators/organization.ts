import { FormControl } from '@angular/forms';
declare var require: any;
var http = require('@angular/http');

export class OrganizationValidator {

  static isValid(control: FormControl): any {

    let url = "https://github.com/";
    http.get(url + control.value).subscribe(response => {
      console.log(response);
    });

    return null;
  }

}
