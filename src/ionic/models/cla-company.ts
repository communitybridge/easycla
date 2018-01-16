export class ClaCompanyModel {

  // This model is based on CLA Company class
  company_external_id: string;
  company_id: string;
  company_manager_id: string;
  company_name: string;
  company_whitelist: Array<string>;
  company_whitelist_patterns: Array<string>;
  date_created: string;
  date_modified: string;
  version: string;

  constructor() {
    
  }

}
