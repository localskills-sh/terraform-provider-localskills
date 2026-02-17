# Browse the most downloaded public skills
data "localskills_explore" "popular" {
  sort = "downloads"
  type = "rule"
}

# Search for skills by keyword
data "localskills_explore" "search" {
  query = "terraform"
  sort  = "newest"
}

output "popular_skills" {
  value = [for s in data.localskills_explore.popular.skills : s.name]
}
