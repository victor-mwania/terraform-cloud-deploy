
provider "google" {
  project     = "PROJECT_ID"
  region      = "us-central1"
  credentials = file("key.json")

}

resource "google_cloud_run_service" "url_shortner" {
  name     = "url-shortner"
  location = "us-central1"

  template {
    spec {
      containers {
        image = "us-central1-docker.pkg.dev/PROJECT_ID/docker/url-shortner:latest"
        ports {
          container_port = 8090
        }

        liveness_probe {
          http_get {
            path = "/health"
          }
        }
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }

}

resource "google_cloud_run_service_iam_member" "run_all_users" {
  service  = "url-shortner"
  location = "us-central1"
  role     = "roles/run.invoker"
  member   = "allUsers"
}

data "google_cloud_run_service" "url_shortner_data" {
  name     = google_cloud_run_service.url_shortner.name
  location = google_cloud_run_service.url_shortner.location
}

output "service_url" {
  value = data.google_cloud_run_service.url_shortner_data.status[0].url
}
