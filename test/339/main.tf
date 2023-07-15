terraform {
  required_providers {
    nexus = {
      source  = "nduyphuong/nexus"
      version = "1.23.0"
    }
  }
}

provider "nexus" {
  insecure = true
  password = "123123123"
  url      = "http://127.0.0.1:8081"
  username = "admin"
}

resource "nexus_repository_bower_hosted" "internal" {
  name   = "internal"
  online = true

  storage {
    blob_store_name                = "default"
    strict_content_type_validation = true
    write_policy                   = "ALLOW"
  }
}

resource "nexus_repository_bower_proxy" "bower_io" {
  name   = "bower-io"
  online = true

  rewrite_package_urls = true

  storage {
    blob_store_name                = "default"
    strict_content_type_validation = true
  }

  proxy {
    remote_url       = "https://registry.bower.io"
    content_max_age  = 1440
    metadata_max_age = 1440
  }

  negative_cache_enabled = true
  negative_cache_ttl     = 1440

  http_client {
    blocked    = false
    auto_block = true
  }
}

resource "nexus_repository_bower_group" "group" {
  name   = "bower-group"
  online = true

  group {
    member_names = [
      nexus_repository_bower_hosted.internal.name,
      nexus_repository_bower_proxy.bower_io.name,
    ]
  }

  storage {
    blob_store_name                = "default"
    strict_content_type_validation = true
  }
}
