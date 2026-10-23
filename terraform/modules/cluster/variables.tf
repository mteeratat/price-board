variable "name" {
  description = "kind cluster name"
  type        = string
}

variable "node_image" {
  description = "kindest/node image, pinned by digest"
  type        = string
}

variable "host_port" {
  description = "Port on the laptop"
  type        = number
  default     = 80
}

variable "node_port" {
  description = "Port inside the kind node (Traefik NodePort in step 4)"
  type        = number
  default     = 30080
}
