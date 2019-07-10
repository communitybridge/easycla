const express = require('express');
const helmet = require('helmet')
const dev = process.env.NODE_ENV !== 'production';
const next = require('next');

const app = next({ dev, dir: './src' });
const  routes  =  require("./src/routes")
const handle = routes.getRequestHandler(app);

app
  .prepare()
  .then(() => {
    express()
    .use(handle)
    .use(helmet())
    .listen(process.env.PORT || 3000, () =>  process.stdout.write(`Server started on set port or 3000\n`))
  })
