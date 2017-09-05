{{ $service_name := env "CONSUL_SERVICE" }}
upstream cla {
  least_conn;
  {{range service $service_name }}server {{.NodeAddress}}:{{.Port}} max_fails=3 fail_timeout=60 weight=1;
  {{end}}
}

server {
  listen 80;
  server_name ${APP_DOMAINS};

  # If request came through the ELB as http, we forward to https
  if ($http_x_forwarded_proto != "https") {
    return 301 https://$host$request_uri;
  }

  location / {
    proxy_pass       http://cla;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header Host $http_host;
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
