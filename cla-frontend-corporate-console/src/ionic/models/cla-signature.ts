export class ClaSignatureModel {

  // This model is based on CLA Signature class
  date_created: string;
  date_modified: string;
  signature_approved: boolean;
  signature_callback_url: string;
  signature_document_major_version: string;
  signature_document_minor_version: string;
  signature_external_id: string;
  signature_id: string;
  signature_project_id: string;
  signature_reference_id: string;
  signature_reference_type: string;
  signature_return_url: string;
  signature_sign_url: string;
  signature_signed: boolean;
  signature_type: string;
  signature_user_ccla_company_id: string;
  version: string;
  domain_whitelist: Array<string>;
  email_whitelist: Array<string>;

  constructor() {

  }

}
