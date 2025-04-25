resource "couchbase-capella_allowlist" "allowlist" {
  organization_id = "<organization_id>"
  project_id      = "<project_id>"
  cluster_id      = "<cluster_id>"
  cidr            = "0.0.0.0/0"
  comment         = "allows all ips"
  expires_at      = "2023-12-30T23:59:59.465Z"
}