import Ajv from 'ajv';

const ajv = new Ajv();
 //To validate API response using schema
export function validateApiResponse (schemaPath,response) {
  cy.fixture(schemaPath).then(
    (schema) => {
  const validate = ajv.compile(schema);
  const isValid = validate(response.body);

  // Assert that the response matches the schema

  if (isValid) {
    cy.log('API response schema is valid');
    expect(isValid, 'API response schema is valid').to.be.true;
  } else {    
    Cypress.on('test:after:run', (test, runnable) => {
        const testName = `${runnable.parent.title} - ${test.title}`
        cy.log(`API response schema is not valid for Test Case : ${testName}`)
        console.log(`API response schema is not valid for Test Case : ${testName}`)
        cy.log('Schema Error : ', validate.errors);
        console.error('Schema Error : ', validate.errors);    
    })
  }

});
};

//To validate & assert 200 response of api
export function validate_200_Status(response){
  expect(response.status).to.eq(200);
  expect(response.body).to.not.be.null;
  const jsonResponse = JSON.stringify(response.body, null, 2);
  cy.log(jsonResponse);
};

let bearerToken = "";
export function getTokenKey(){
  
  cy.request({
    method: 'POST',
    url: Cypress.env("AUTH0_TOKEN_API"),
   
    body: {
      "grant_type": "http://auth0.com/oauth/grant-type/password-realm",
      "realm": "Username-Password-Authentication",
      "username":Cypress.env("AUTH0_USER_NAME"),
      "password":Cypress.env("AUTH0_PASSWORD"),
      "client_id":Cypress.env("AUTH0_CLIENT_ID"),
      "audience": "https://api-gw.dev.platform.linuxfoundation.org/",
      "scope": "access:api openid profile email"
    }
  }).then(response => {        
    expect(response.status).to.eq(200);       
    bearerToken = response.body.access_token;    
    cy.window().then((win) => {
      win.localStorage.setItem('bearerToken', response.body.access_token);
    });
  });
};
