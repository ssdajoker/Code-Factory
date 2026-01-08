# Reports

This directory contains generated reports from Factory's analysis modes.

## Report Types

- **Review Reports** - Code compliance analysis against specifications
- **Change Orders** - Tracking specification drift over time
- **Rescue Reports** - Reverse-engineered specifications from existing code

## Generating Reports

Reports are automatically generated when you run Factory modes:

```bash
factory review          # Generates review report
factory change-order    # Generates change order report
factory rescue          # Generates rescue report
```

## Report Format

All reports are in Markdown format for easy reading and version control.

Example report structure:

```markdown
# [Report Type] Report

**Date:** YYYY-MM-DD
**Specification:** spec_name.md
**Overall Compliance:** XX%

## Summary
Brief overview of findings

## Detailed Findings
### Critical Issues
### Warnings
### Compliant Items

## Recommendations
Actionable next steps
```

## Best Practices

1. Review reports regularly
2. Address critical issues first
3. Track trends over time
4. Share reports with your team
5. Use reports to improve specs and code
