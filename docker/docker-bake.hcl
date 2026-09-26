target "docker-metadata-action" {}

target "release" {
  inherits = ["docker-metadata-action"]
}
