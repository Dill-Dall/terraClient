include "root" {
  path = find_in_parent_folders()
}

terraform {
  source = "/home/thomas/training/go/terragrunt/modules//s3"
}

inputs = {
  configurations = merge(
		read_terragrunt_config("configurations/testing-them-buckets.hcl").inputs,
	)
}