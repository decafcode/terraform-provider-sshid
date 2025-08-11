# AWS KMS

# ssh-keygen can issue certificates using this PKCS#11 provider:
# https://github.com/JackOfMostTrades/aws-kms-pkcs11

data "aws_kms_public_key" "example" {
  key_id = "alias/example"
}

output "example_aws" {
  value = provider::sshid::fromspki(data.aws_kms_public_key.example.public_key_pem)
}

# GCP KMS (untested)

# ssh-keygen can issue certificates using this PKCS#11 provider:
# https://github.com/GoogleCloudPlatform/kms-integrations

data "google_kms_key_ring" "example" {
  name     = "my-key-ring"
  location = "us-central1"
}

data "google_kms_crypto_key" "example" {
  name     = "my-crypto-key"
  key_ring = data.google_kms_key_ring.example.id
}

data "google_kms_crypto_key_latest_version" "example" {
  crypto_key = data.google_kms_crypto_key.example.id
}

output "example_gcp" {
  value = provider::sshid::fromspki(data.google_kms_crypto_key_latest_version.example.public_key)
}

# Azure Key Vault (untested)

# Not required: data.azurerm_key_vault_key natively provides an attribute that
# exports a public key in SSH format.

# ssh-keygen can issue certificates using this PKCS#11 provider:
# https://github.com/jepio/azure-keyvault-pkcs11
