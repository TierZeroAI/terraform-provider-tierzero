terraform {
  required_providers {
    tierzero = {
      source = "tierzero/tierzero"
    }
  }
}

provider "tierzero" {
  # API key can be set here or via TIERZERO_API_KEY environment variable
  # api_key = var.tierzero_api_key

  # Base URL defaults to https://api.tierzero.ai
  # base_url = "https://api.tierzero.ai"
}
