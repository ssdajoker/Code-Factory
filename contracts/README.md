# Contracts (Specifications)

This directory contains specification documents that describe what the software 
should do. These serve as the source of truth for the project.

## Current Specifications

- `factory_bootstrap_spec.md` - Complete onboarding and initialization flow
- `system_architecture.md` - Overall system design and component interaction
- `mode_specifications.md` - Detailed specs for INTAKE, REVIEW, CHANGE_ORDER, RESCUE modes

## Creating New Specifications

Use Factory's INTAKE mode to create new specifications:

```bash
factory intake
```

Or create manually using the template structure:

```markdown
# [Feature Name]

**Version:** 1.0.0
**Status:** Draft
**Created:** YYYY-MM-DD
**Author:** Your Name

## Overview
Brief description

## Requirements
### Functional Requirements
### Non-Functional Requirements

## Technical Design
### Architecture
### Components
### Data Models
### API Endpoints

## Implementation Plan
## Testing Strategy
## Security Considerations
## Dependencies
```

## Best Practices

1. Keep specs up-to-date with code changes
2. Use clear, actionable language
3. Include examples and diagrams
4. Version your specifications
5. Review specs with your team
