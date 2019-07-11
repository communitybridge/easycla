// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

const  routes  =  require('next-routes')
module.exports  =  routes()
.add('home', '/')
.add('notfound', '/*')
