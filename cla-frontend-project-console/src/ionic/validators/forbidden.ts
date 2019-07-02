// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { AbstractControl, ValidatorFn } from '@angular/forms';

/** A hero's name can't match the given regular expression */
export function forbiddenValidator(nameRe: RegExp): ValidatorFn {
  return (control: AbstractControl): {[key: string]: any} => {
    const forbidden = nameRe.test(control.value);
    return forbidden ? {'forbidden': {value: control.value}} : null;
  };
}
