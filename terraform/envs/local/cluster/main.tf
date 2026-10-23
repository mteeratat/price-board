provider "kind" {}

module "cluster" {
  source = "../../../modules/cluster"

  name       = "price-board"
  node_image = "kindest/node:v1.36.4@sha256:099e049362a1526b2db71494e1947aae99bd16290d7c895f2b7ea312e3cbfaed"
}
