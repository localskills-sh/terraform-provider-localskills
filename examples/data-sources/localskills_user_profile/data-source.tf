data "localskills_user_profile" "me" {}

output "my_email" {
  value = data.localskills_user_profile.me.email
}

output "my_username" {
  value = data.localskills_user_profile.me.username
}
