---
# ghv-clfp
title: Fix template-syntax-valid rule to recognize {{define}} as control structure
status: completed
type: bug
priority: normal
created_at: 2026-01-18T01:52:40Z
updated_at: 2026-01-18T01:53:45Z
---

The template-syntax-valid rule doesn't recognize `{{define}}` as a valid control structure. Need to add it to the list of recognized template directives.