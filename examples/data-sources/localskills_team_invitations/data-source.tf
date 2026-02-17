# List pending invitations for a team
data "localskills_team_invitations" "pending" {
  tenant_id = localskills_team.engineering.id
}

output "pending_emails" {
  value = [for inv in data.localskills_team_invitations.pending.invitations : inv.email]
}
