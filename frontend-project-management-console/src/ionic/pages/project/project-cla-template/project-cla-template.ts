import {Component, ViewChild} from "@angular/core";
import {
  NavController,
  NavParams,
  IonicPage, Nav, Events, AlertController
} from "ionic-angular";
import {ClaService} from "../../../services/cla.service";
import {Restricted} from "../../../decorators/restricted";
import { DomSanitizer } from '@angular/platform-browser';
@Restricted({
  roles: ["isAuthenticated", "isPmcUser"]
})
@IonicPage({
  segment: "project/:projectId/cla/template/:projectTemplateId"
})
@Component({
  selector: "project-cla-template",
  templateUrl: "project-cla-template.html"
})
export class ProjectClaTemplatePage {
  loading: any;
  sfdcProjectId: string;
  projectTemplateId: string;
  templates: any[] = [];
  selectedTemplate = {};
  templateValues = {};
  pdfPath = {
    CorporatePDFURL: '/assets/sample-cla-pdf.pdf',
    IndividualPDFURL: '/assets/sample-cla-pdf.pdf'
  };
  currentPDF = 'cclaPdfUrl';
  step = 'selection';

  @ViewChild(Nav) nav: Nav;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    public claService: ClaService,
    public sanitizer: DomSanitizer
  ) {
    this.sfdcProjectId = navParams.get("projectId");
    this.projectTemplateId = navParams.get("projectTemplateId");
    this.getDefaults();
  }

  getDefaults() {
    this.getTemplates();
  }

  getTemplates() {
    this.claService.getTemplates().subscribe(templates => this.templates = templates);
  }

  ngOnInit() {
  }

  getPdfPath() {
    return this.sanitizer.bypassSecurityTrustResourceUrl(this.pdfPath[this.currentPDF]);
  }

  showPDF (type) {
    this.currentPDF = type;
  }

  selectTemplate(template) {
    this.selectedTemplate = template;
    this.step = 'values';
  }

  reviewSelectedTemplate() {
    this.claService.postClaGroupTemplate(this.projectTemplateId, this.templateValues)
      .subscribe(response => {
        this.pdfPath = response;
        this.goToStep('review');
      })
  }

  goToStep(step) {
    this.step = step;
  }

  backToProject() {
    this.navCtrl.pop();
  }

  viewDocusign() {}

  submit() {}
}
