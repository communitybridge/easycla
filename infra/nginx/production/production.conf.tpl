server {
  listen 80;
  server_name ${APP_DOMAINS};
  root /srv/app/src/www;
  gzip on;
  index index.html;

  # If request came through the ELB as http, we forward to https
  if ($http_x_forwarded_proto != "https") {
    return 301 https://$host$request_uri;
  }

  location / {
    try_files $uri $uri/ /index.html;
  }
}

server {
  listen 80;
  server_name ${EC2_LOCAL_IP} 127.0.0.1 localhost;

  # Server status for the NewRelic nginx monitoring
  location = /status {
    access_log off;
    stub_status on;
    allow 127.0.0.1;
    deny all;
  }

  # Server status for the AWS ELB
  location = /elb-status {
    access_log off;
    return 200;
  }
}
