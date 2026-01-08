# Project Contracts

This directory contains all project specifications and architecture documentation for Code-Factory.

## Structure

### `/specs/`
Feature specifications - detailed requirements and implementation contracts for individual features.

*Coming soon: Feature specifications will be placed here as the project develops.*

### `/architecture/`
System architecture documents - high-level design decisions and technical architecture.

**Current documents:**
- `system_architecture.md` - Complete system architecture and design patterns

### `/decisions/`
Architecture Decision Records (ADRs) - track important architectural decisions over time.

*Coming soon: ADRs will be added as architectural decisions are made.*

## Foundation Specifications

The following foundational specifications define the core system:

### 1. **factory_bootstrap_spec.md**
Comprehensive specification for the one-click onboarding experience (`factory init`).

**Contents:**
- One-click onboarding flow
- GitHub OAuth and App installation
- LLM configuration (Ollama vs BYOK)
- Secret management strategy
- Fallback scenarios and error handling
- UX design and user feedback
- Technical implementation details

### 2. **system_architecture.md**
Complete system architecture including component design, data flow, and technology stack.

**Contents:**
- High-level architecture overview
- Component architecture (TUI, LLM, GitHub, Modes)
- Data flow diagrams
- Storage model (git-native)
- Integration points
- Security architecture
- Performance considerations

### 3. **mode_specifications.md**
Detailed specifications for all four operational modes.

**Contents:**
- **INTAKE Mode** - Capture requirements as specs
- **REVIEW Mode** - Analyze code against specs
- **CHANGE_ORDER Mode** - Implement changes from specs
- **RESCUE Mode** - Debug and fix issues
- Cross-mode features and shared capabilities

## Using Specifications

These specifications serve as the **source of truth** for implementation. When developing:

1. **Read the relevant spec first** - Understand requirements before coding
2. **Follow the spec exactly** - Deviations should be documented and approved
3. **Update specs when requirements change** - Keep specs in sync with reality
4. **Reference specs in PRs** - Link code changes to their specifications
5. **Review code against specs** - Use Code-Factory's REVIEW mode

## Spec Format

All specifications follow this structure:

```markdown
---
id: unique-identifier
title: Feature Title
status: draft | approved | implemented | deprecated
created: YYYY-MM-DD
updated: YYYY-MM-DD
author: username
tags: [tag1, tag2]
priority: low | medium | high | critical
---

# Feature Title

## Overview
Brief description...

## Requirements
### Functional
### Non-Functional

## Technical Specification
...

## Testing Strategy
...
```

## Contributing

When adding new specifications:

1. Use the template format above
2. Place in appropriate directory (specs/, architecture/, decisions/)
3. Link to related specifications
4. Update this README's structure section
5. Commit with clear message: `docs: Add specification for X`

## Questions?

For questions about specifications, see:
- [ARCHITECTURE.md](../docs/ARCHITECTURE.md) - Technical architecture details
- [README.md](../README.md) - Project overview
- GitHub Discussions - Community discussions

---

**Maintained by:** Code-Factory Core Team  
**Last updated:** 2026-01-07
