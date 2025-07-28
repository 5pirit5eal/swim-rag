variable "project_id" {
  description = "The GCP Project ID"
  type        = string
}

variable "region" {
  description = "The GCP Region"
  type        = string
}

variable "apis" {
  description = "The GCP APIs to enable"
  type        = list(string)
}

variable "outputs_location" {
  description = "The GCP output locations"
  type        = string
  default     = null
}

variable "dbname" {
  description = "The GCP Cloud SQL Database Name"
  type        = string
}

variable "dbuser" {
  description = "The GCP Cloud SQL Database User"
  type        = string
}

variable "dbtier" {
  description = "The GCP Cloud SQL Database Tier"
  type        = string
}

variable "github_token" {
  description = "The GitHub token"
  type        = string
}

variable "github_app_installation_id" {
  description = "The GitHub app installation ID"
  type        = number
}

variable "github_uri" {
  description = "The GitHub URI"
  type        = string
  default     = null
}

variable "model" {
  description = "The model name"
  type        = string
}
variable "embedding_name" {
  description = "The embedding name"
  type        = string
}
variable "embedding_model" {
  description = "The embedding model"
  type        = string
}
variable "embedding_size" {
  description = "The embedding size"
  type        = number
}
variable "port" {
  description = "The port number"
  type        = number
}
variable "log_level" {
  description = "The log level"
  type        = string
}

variable "service_cpu" {
  description = "The CPU for the service"
  type        = number
}

variable "service_memory" {
  description = "The memory for the service"
  type        = string
}

variable "service_timeout" {
  description = "The timeout for the service"
  type        = number
}

variable "backend_url" {
  description = "The backend URL for the MCP service"
  type        = string
  default     = "http://localhost:8080"
}