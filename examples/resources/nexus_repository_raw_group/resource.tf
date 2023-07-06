resource "nexus_repository_raw_hosted" "internal" {
  name   = "raw-internal"
  online = true

  storage {
    blob_store_name                = "default"
    strict_content_type_validation = false
    write_policy                   = "ALLOW"
  }

}


resource "nexus_repository_raw_group" "group" {
  name   = "raw-group"
  online = true

  group {
    member_names {
      name  = nexus_repository_raw_hosted.internal.name
      order = 1
    }
  }

  storage {
    blob_store_name                = "default"
    strict_content_type_validation = true
  }
}
