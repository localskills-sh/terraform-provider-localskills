data "localskills_user_tokens" "all" {}

output "token_names" {
  value = [for t in data.localskills_user_tokens.all.tokens : t.name]
}
