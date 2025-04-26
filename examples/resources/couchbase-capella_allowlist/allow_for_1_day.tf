resource "couchbase-capella_allowlist" "allowlist" {
  organization_id = "<organization_id>"
  project_id      = "<project_id>"
  cluster_id      = "<cluster_id>"
  cidr            = "95.75.125.11/32"
  comment         = "allows only my ip"
  expires_at      = time_offset.one_day.base_rfc3339
}

resource "time_offset" "one_day" {
  offset_days = 1
}