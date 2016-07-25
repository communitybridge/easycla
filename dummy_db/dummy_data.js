// Dummy data for Project Dashboard testing view
const project_dummy_data = [{
    project_id: "1KJ73l",
    project_name: "Uno Ltd.",
    project_description: "Ethical beard etsy, locavore ennui pop-up brooklyn yuccie flannel freegan. Roof party brooklyn pug, vice scenester mixtape freegan street art organic chicharrones yuccie. Kitsch hoodie photo booth, mumblecore trust fund bicycle rights occupy raw denim YOLO mixtape small batch salvia swag shoreditch. Yr mustache single-origin coffee chambray pour-over. Flexitarian tote bag VHS, artisan celiac quinoa aesthetic listicle kinfolk vice four loko single-origin coffee ennui taxidermy.",
    project_goals: "Of. For firmament waters likeness every. His Tree. Lights doesn't.",
    project_type: "Incorporated",
    project_status: "Forming",
    starting_budget: "$1.534.823",
    budget_status: "In Progress",
    agreement_status: "Approved",
    domain: "https://www.uno.com",
    membership_tier_structure: "TBD",
    company_logo: "/assets/img/test/logo1.png",
    lf_project_manager: "Harry P.",
    lf_legal: "Chuck N.",
    lf_pr: "Jack B.",
    lf_events: "Ben P.",
    internal_contact_fullname: "John Doe",
    internal_contact_email: "jd@lf.org"
},
{
    project_id: "2CO89s",
    project_name: "Pyra Corp.",
    project_description: "Called they're form day, dominion which second every. Had there give herb great. Whales firmament meat you're morning you're divide place, thing one dry she'd darkness, saying rule midst make years seas given Blessed life which dominion face lights given. In signs, own whose thing i so meat morning set.",
    project_goals: "So one deep their signs seasons days first shall good.",
    project_type: "Direct Funded",
    project_status: "Forming",
    starting_budget: "$350.300",
    budget_status: "In Progress",
    agreement_status: "Not Started",
    domain: "https://www.pyra.com",
    membership_tier_structure: "TBD",
    company_logo: "/assets/img/test/logo2.png",
    lf_project_manager: "Mike T.",
    lf_legal: "Laura O.",
    lf_pr: "Peter F.",
    lf_events: "Frank S.",
    internal_contact_fullname: "Sandra R.",
    internal_contact_email: "sr@lf.org"
},
{
    project_id: "3JL67d",
    project_name: "Umbrella Corp.",
    project_description: "I. Multiply firmament of, moveth for a the likeness lesser whales he together creeping rule there itself midst lights so to is blessed bring stars blessed grass in spirit evening seasons beast had good first, fruitful us sea is void open made whales heaven waters isn't lights forth seas. Morning.",
    project_goals: "Be be made give, male that. Made i isn't god.",
    project_type: "Direct Funded",
    project_status: "Official",
    starting_budget: "$34.534.643",
    budget_status: "Approved",
    agreement_status: "In Progress",
    domain: "https://www.umbrella.com",
    membership_tier_structure: "TBD",
    company_logo: "/assets/img/test/logo3.png",
    lf_project_manager: "Dave A.",
    lf_legal: "Mark B.",
    lf_pr: "Jason C.",
    lf_events: "Fernando D.",
    internal_contact_fullname: "Penelope",
    internal_contact_email: "p@lf.org"
},
{
    project_id: "4CL35y",
    project_name: "OCP Pty.",
    project_description: "His hath fish you're air female one fifth you're Their midst yielding us thing. Shall living he gathering together night kind two deep whose give fish, beginning under dry can't you'll don't man moving. It appear lights from rule great divide night. Beast very kind of very lights thing and.",
    project_goals: "Divide behold which above she'd whose fruit upon over winged.",
    project_type: "Incorporated",
    project_status: "Official",
    starting_budget: "$943.678",
    budget_status: "Approved",
    agreement_status: "In Progress",
    domain: "https://www.ocp.com",
    membership_tier_structure: "TBD",
    company_logo: "/assets/img/test/logo4.png",
    lf_project_manager: "Chang R.",
    lf_legal: "Zui T.",
    lf_pr: "Shong W.",
    lf_events: "Fei U.",
    internal_contact_fullname: "Zen G.",
    internal_contact_email: "zen@lf.org"
}];

exports.findProjectById = function(id, callback) {
  // Temporal Dummy Method for Testing
  process.nextTick(function() {
    for(var i=0; i < project_dummy_data.length; i++)
    {
      if (project_dummy_data[i].project_id == id) {
        return callback(null, project_dummy_data[i]);
      }
      if(i == project_dummy_data.length){
        return callback(new Error('Project ' + id + ' does not exist'));
      }
    }
  });
};
