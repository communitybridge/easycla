const AWS = require('aws-sdk');
const SSM = require('aws-sdk/clients/ssm');
const program = require('commander');

program
  .version('1.0.0')
  .option('-d, --debug', 'output extra debugging')
  .requiredOption('-r, --region <region>', 'specifies the AWS region')
  .requiredOption('-p, --parameter <name>', 'the parameter name')
  .parse(process.argv);
if (program.debug) console.log(program.opts());


// Configure AWS
AWS.config.update({ region: program.region });

console.log(`Querying SSM Parameter: ${program.parameter} in region ${program.region}...`);
const query = {
  "Name": program.parameter,
  "WithDecryption": false,
};

const ssm = new SSM();
let param = ssm.getParameter(query, (err, data) => {
  if (err == null) {
    console.log(`${data.Parameter.Name} = ${data.Parameter.Value}`)
    console.log('Details: %o', data);
  } else {
    console.log(`Error fetching parameter ${program.parameter}. Error code: ${err.code}, Message: ${err.message}`);
    if (program.debug) {
      console.log('error = %o', err);
    }
  }
});
