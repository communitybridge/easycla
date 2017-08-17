if (process.env['NEWRELIC_LICENSE']) require('newrelic');
var express = require('express');
var config = require('config');
var path = require('path');
var session = require('express-session');
var RedisStore = require('connect-redis')(session);

var app = express();

// App config
app.set('views', path.join(__dirname, 'views'));

// Middleware
app.use(require('morgan')('combined')); // HTTP request logger middleware
app.use(require('cookie-parser')());
app.use(session({
  store: new RedisStore({
    host: config.get('console.redisHost'),
    port: config.get('console.redisPort')
  }),
  secret: config.get('session.secret'),
  resave: false,
  saveUninitialized: false
}));

app.disable('x-powered-by');

// Serve Angular / Ionic App
app.use('/', express.static(path.join(__dirname, 'app/www')));

app.get('*', function(req, res) {
  res.redirect('/');
});

// AWS  nginx proxy server uses 8081 by default
const appPort = 8081;
app.listen(appPort, function() {});
console.log("Node App listening port: ", appPort);

// Docker resolves the port mapping to "console.endpoint.port:8081"
const pmcURL = config.get('console.endpoint');
console.log('Docker project-management-console-url: ' + pmcURL);


displayBanner();

function displayBanner() {
  let fs = require('fs')
  let max = 4;
  let min = 1;
  let version = Math.floor (Math.random() * (max - min + 1) ) + min;
  let banner = './config/banners/' + version + '.txt';
  fs.readFile(banner, 'utf8', function(err, data) {
    if (err) throw err;
    console.log(data)
  });
};

/******************************
* ██████╗ ███╗   ███╗ ██████╗ *
* ██╔══██╗████╗ ████║██╔════╝ *
* ██████╔╝██╔████╔██║██║      *
* ██╔═══╝ ██║╚██╔╝██║██║      *
* ██║     ██║ ╚═╝ ██║╚██████╗ *
* ╚═╝     ╚═╝     ╚═╝ ╚═════╝ *
*******************************/
