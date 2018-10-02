import { Component } from "@angular/core";
import {
  NavController,
  ModalController,
  NavParams,
  IonicPage
} from "ionic-angular";
import { ClaService } from "../../services/cla.service";
import { ClaCompanyModel } from "../../models/cla-company";
import { ClaUserModel } from "../../models/cla-user";
import { ClaSignatureModel } from "../../models/cla-signature";
import { SortService } from "../../services/sort.service";
import { RolesService } from "../../services/roles.service";
import { Restricted } from "../../decorators/restricted";

@Restricted({
  roles: ["isAuthenticated"]
})
@IonicPage({
  segment: "company/:companyId"
})
@Component({
  selector: "company-page",
  templateUrl: "company-page.html"
})
export class CompanyPage {
  companyId: string;
  company: ClaCompanyModel;
  manager: ClaUserModel;
  companySignatures: ClaSignatureModel[];
  projects: any;
  loading: any;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    private claService: ClaService,
    public modalCtrl: ModalController,
    private rolesService: RolesService // for @Restricted
  ) {
    this.companyId = navParams.get("companyId");
    this.getDefaults();
  }

  getDefaults() {
    this.loading = {
      companySignatures: true
    };
    this.company = new ClaCompanyModel();
    this.projects = {};
  }

  ngOnInit() {
    this.getCompany();
    this.getCompanySignatures();
  }

  getCompany() {
    this.claService.getCompany(this.companyId).subscribe(response => {
      this.company = response;
      this.getUser(this.company.company_manager_id);
    });
  }

  getUser(userId) {
    this.claService.getUser(userId).subscribe(response => {
      this.manager = response;
    });
  }

  getCompanySignatures() {
    this.claService.getCompanySignatures(this.companyId).subscribe(response => {
      this.companySignatures = response.filter(signature =>
        signature.signature_signed === true
      );
      for (let signature of this.companySignatures) {
          this.getProject(signature.signature_project_id);
      }
    });
  }

  getProject(projectId) {
    this.claService.getProject(projectId).subscribe(response => {
      this.projects[projectId] = response;
    });
  }

  openProjectPage(projectId) {
    this.navCtrl.push("ProjectPage", {
      companyId: this.companyId,
      projectId: projectId
    });
  }

  openCompanyModal() {
    let modal = this.modalCtrl.create("AddCompanyModal", {
      company: this.company
    });
    modal.onDidDismiss(data => {
      // A refresh of data anytime the modal is dismissed
      this.getCompany();
    });
    modal.present();
  }

  openWhitelistEmailModal() {
    let modal = this.modalCtrl.create("WhitelistModal", {
      type: "email",
      company: this.company
    });
    modal.onDidDismiss(data => {
      // A refresh of data anytime the modal is dismissed
      this.getCompany();
    });
    modal.present();
  }

  openWhitelistDomainModal() {
    let modal = this.modalCtrl.create("WhitelistModal", {
      type: "domain",
      company: this.company
    });
    modal.onDidDismiss(data => {
      // A refresh of data anytime the modal is dismissed
      this.getCompany();
    });
    modal.present();
  }

  openProjectsCclaSelectModal() {
    let modal = this.modalCtrl.create("ProjectsCclaSelectModal", {
      company: this.company
    });
    modal.onDidDismiss(data => {
      // A refresh of data anytime the modal is dismissed
      this.getCompany();
    });
    modal.present();
  }
}
