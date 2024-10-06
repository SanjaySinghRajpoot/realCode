variable "region" {
    default = "ap-south-1"
}

variable "key_name" {
  description = "key_name"
  type        = string
}

variable "public_key" {
  description = "public_key"
  type        = string
}

variable "private_key" {
  description = "private_key"
  type        = string
}