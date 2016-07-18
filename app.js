var express = require('express');
var passport = require('passport');
var GoogleStrategy = require('passport-google-oauth20').Strategy;
var path = require('path');
var bodyParser = require('body-parser');

var app = express();

// App config
app.set('view engine', 'ejs');
app.set('views', path.join(__dirname, 'views'));

// Middleware
app.use(express.static(path.join(__dirname, 'public')));

// Modules required by Angular 2
app.use('/node_modules/zone.js/dist/', express.static(path.join(__dirname, 'node_modules/zone.js/dist/')));
app.use('/node_modules/reflect-metadata/', express.static(path.join(__dirname, 'node_modules/reflect-metadata/')));
app.use('/node_modules/systemjs/dist/', express.static(path.join(__dirname, 'node_modules/systemjs/dist/')));
app.use('/node_modules/core-js/client/', express.static(path.join(__dirname, 'node_modules/core-js/client/')));
app.use('/node_modules/@angular/', express.static(path.join(__dirname, 'node_modules/@angular/')));
app.use('/node_modules/angular2-in-memory-web-api/', express.static(path.join(__dirname, 'node_modules/angular2-in-memory-web-api/')));
app.use('/node_modules/rxjs/', express.static(path.join(__dirname, 'node_modules/rxjs/')));

//app.use(require('morgan')('combined')); // HTTP request logger middleware
app.use(require('cookie-parser')());
app.use(require('body-parser').urlencoded({ extended: true }));
app.use(require('express-session')({ secret: process.env['SESSION_SECRET'] != null ? process.env['SESSION_SECRET'] : 'lhb.sdu3erw lwfe rlfwe oThge3 823dwj34 @#kbdwe3 ghdklnj32lj l2303', resave: false, saveUninitialized: false }));

app.use(passport.initialize());
app.use(passport.session());

app.use(function (req, res, next) {
  res.locals.req = req;
  next();
});

// Routes
var routes = require('./routes')
app.use(routes);

const port = process.env['UI_PORT'] != null ? process.env['UI_PORT'] : 8081
app.listen(port, function() {});

passport.use(new GoogleStrategy({
    clientID: process.env.GOOGLE_CLIENT_ID,
    clientSecret: process.env.GOOGLE_CLIENT_SECRET,
    callbackURL: "http://lf-integration-console-sandbox.us-west-2.elasticbeanstalk.com/auth/google/callback"
    // callbackURL: "http://localhost:8081/auth/google/callback"
  },
  function(accessToken, refreshToken, profile, callback) {
    //User.findOrCreate({ googleId: profile.id }, function (err, user) {
      //return callback(err, user);
    //});
    return callback(null, profile);
  }
));

passport.serializeUser(function(user, callback) {
  callback(null, user);
});

passport.deserializeUser(function(obj, callback) {
  callback(null, obj);
});
