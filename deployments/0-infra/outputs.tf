locals {
  outputs_location = try(pathexpand(var.outputs_location), "")
  tfvars = {
    project_id = var.project_id
    region     = var.region
    csql_db = {
      id   = try(google_sql_database.main_db.id, null)
      name = var.dbname
      tier = var.dbtier
    }
    csql_instance = {
      "connection_name" = try(google_sql_database_instance.main.connection_name, null)
      "uri"             = try(google_sql_database_instance.main.self_link, null)
      "public_ip"       = try(google_sql_database_instance.main.public_ip_address, null)
      "private_ip"      = try(google_sql_database_instance.main.private_ip_address, null)
    }
    secret_ids       = local.secret_ids
    bucket_name      = google_storage_bucket.exported_pdfs.name
    backend_env_vars = local.backend_env_variables
    mcp_env_vars     = local.mcp_env_variables
    artifactregistry = {
      "repository" = google_artifact_registry_repository.docker.name
      "location"   = google_artifact_registry_repository.docker.location
    }
    iam = {
      "cloud_build" = {
        "email" = try(google_service_account.cloud_build_sa.email, null)
        "id"    = try(google_service_account.cloud_build_sa.id, null)
      }
      "cloud_run" = {
        "email" = try(google_service_account.cloud_run_sa.email, null)
        "id"    = try(google_service_account.cloud_run_sa.id, null)
      }
    }
  }
}


output "tfvars" {
  description = "Terraform variable files for the following stages."
  value       = local.tfvars
}


resource "local_file" "tfvars" {
  for_each        = var.outputs_location == null ? {} : { 1 = 1 }
  file_permission = "0644"
  filename        = "${local.outputs_location}/infra.auto.tfvars.json"
  content         = jsonencode(local.tfvars)
}