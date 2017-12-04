import { Component, Injector } from '@angular/core';
//
// // @Restricted()
// export function Restricted(data: any = {}) {
//   console.log(data);
//   return (target:any, key:string, descriptor:any) => {
//     console.log('target');
//     console.log(target);
//     console.log('key');
//     console.log(key);
//     console.log('descriptor');
//     console.log(descriptor);
//     // var original = descriptor.value;
//     // var localMessage = message.replace('{name}', name);
//     // descriptor.value = function() {
//     //   console.warn(`Function ${name} is deprecated: ${localMessage}`);
//     //   return instance
//     // };
//     // return descriptor;
//   };
// }

// @deprecate('Please use other methode')
export function deprecate(message: string = '{name}') {
  return (instance:any, name:string, descriptor:any) => {
    // var original = descriptor.value;
    console.log('descriptor');
    console.log(descriptor);
    // let descriptor = {};
    var localMessage = message.replace('{name}', name);
    descriptor.value = function() {
      console.log(`Function ${name} is deprecated: ${localMessage}`);
      return instance;
    };
    return descriptor;
  };
}
