provider "multiverse" {}

resource "azureappconfig_config_setting" "test_config" {
  config_name = "app-config-terraform"
  key = "Thor"
  value = "44.133.122.111"
  label = "golang"
   
output "resources" {
  value = "${multiverse_custom_resource.spotinst_targetset_and_rules.id}"
}

output "test_targetset_id" {
  value = "${multiverse_custom_resource.spotinst_targetset_and_rules["testTargetSet"]}"
}

output "control_targetset_id" {
  value = "${multiverse_custom_resource.spotinst_targetset_and_rules["controlTargetSet"]}"
}