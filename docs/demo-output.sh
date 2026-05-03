#!/bin/sh

cat <<'EOF'
Planning infrastructure changes...

  # aws_iam_policy.deploy will be changed
  ~ resource aws_iam_policy.deploy {
      policy = jsonencode({
        Statement = [
          {
            Action = *
            Effect = Allow
            Resource = *
          }
        ]
      })
    }

  # aws_db_instance.main will be updated in-place
  ~ resource aws_db_instance.main {
      instance_class = db.t3.medium -> db.t4g.medium
      backup_retention_period = 7 -> 1
    }

Plan: 0 to add, 2 to change, 0 to destroy.
EOF
