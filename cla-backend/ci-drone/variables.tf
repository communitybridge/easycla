variable "access_key" {
  description = "The AWS access key value"
}

variable "secret_key" {
  description = "The AWS secret key value"
}

variable "domain" {
  description = "Domain name for the resource, shouldn't include www"
}

variable "github_client_id" {
  description = "The github client id, used for Auth0 "
}

variable "github_client_secret" {
  description = "The github client secret"
}

variable "key_file" {
  description = "The ssh public key used to authenticate with the server"
}

variable "email" {
  description = "Email address to receive security alerts from Let's Encrypt"
}

variable "safe_network_ip" {
  description = "Public ip address to allow ssh commuinication over"
}

variable "bucket" {
  description = "The terraform state bucket"
}
