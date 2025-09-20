# Contributing to CMDB Lite

We welcome contributions to CMDB Lite! This document provides guidelines and instructions for contributing to the project.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Process](#development-process)
- [Pull Request Process](#pull-request-process)
- [Coding Standards](#coding-standards)
- [Documentation](#documentation)
- [Reporting Issues](#reporting-issues)

## Code of Conduct

Please read and follow our [Code of Conduct](CODE_OF_CONDUCT.md) to ensure a welcoming and inclusive environment for all contributors.

## Getting Started

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/yourusername/cmdb-lite.git
   cd cmdb-lite
   ```
3. Add the original repository as a remote:
   ```bash
   git remote add upstream https://github.com/your-org/cmdb-lite.git
   ```
4. Create a new branch for your changes:
   ```bash
   git checkout -b feature/your-feature-name
   ```

## Development Process

### Setting Up Your Development Environment

1. Follow the instructions in the [README.md](README.md) to set up your development environment
2. Ensure you have all the necessary prerequisites installed
3. Set up your environment variables by copying the example files:
   ```bash
   cp .env.example .env
   cp backend/.env.example backend/.env
   cp frontend/.env.example frontend/.env
   cp database/.env.example database/.env
   ```

### Making Changes

1. Create a new branch for each feature or bug fix:
   ```bash
   git checkout -b feature/your-feature-name
   # or
   git checkout -b fix/your-fix-name
   ```
2. Make your changes following the coding standards outlined below
3. Test your changes thoroughly
4. Update documentation as needed
5. Commit your changes with a clear and descriptive commit message:
   ```bash
   git commit -m "Add feature: your feature description"
   # or
   git commit -m "Fix issue: your fix description"
   ```

### Keeping Your Fork Up-to-Date

Before submitting a pull request, make sure your branch is up-to-date with the upstream repository:

```bash
git fetch upstream
git rebase upstream/main
```

## Pull Request Process

1. Ensure your code follows the coding standards
2. Write tests for new functionality and ensure existing tests pass
3. Update documentation as needed
4. Push your changes to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```
5. Create a pull request against the main branch of the original repository
6. Provide a clear description of your changes in the pull request
7. Link any relevant issues in the pull request description
8. Wait for your pull request to be reviewed and address any feedback

### Pull Request Template

When creating a pull request, please use the following template:

```markdown
## Description
[Provide a clear description of the changes made]

## Type of Change
- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Documentation update

## Testing
[Describe how you tested your changes]

## Checklist
- [ ] My code follows the coding standards of this project
- [ ] I have performed a self-review of my own code
- [ ] I have commented my code, particularly in hard-to-understand areas
- [ ] I have made corresponding changes to the documentation
- [ ] My changes generate no new warnings
- [ ] I have added tests that prove my fix is effective or that my feature works
- [ ] New and existing unit tests pass locally with my changes
- [ ] Any dependent changes have been merged and published in downstream modules
```

## Coding Standards

### Go (Backend)

- Follow the [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- Use `gofmt` to format your code:
  ```bash
  gofmt -w .
  ```
- Use `golint` to check for style issues:
  ```bash
  golint ./...
  ```
- Write clear, concise, and idiomatic Go code
- Use meaningful variable and function names
- Add comments for exported functions, types, and constants
- Handle errors appropriately

### JavaScript/Vue (Frontend)

- Follow the [Vue Style Guide](https://v2.vuejs.org/v2/style-guide/)
- Use ESLint to check for style issues:
  ```bash
  npm run lint
  ```
- Use Prettier to format your code:
  ```bash
  npm run format
  ```
- Use meaningful variable and function names
- Prefer arrow functions for callbacks
- Use template literals for string interpolation
- Use const and let instead of var

### SQL (Database)

- Use uppercase for SQL keywords (SELECT, INSERT, UPDATE, etc.)
- Use lowercase for table and column names
- Use snake_case for table and column names
- Add comments to explain complex queries
- Use transactions for multiple related operations

### General Guidelines

- Write clear, concise, and self-documenting code
- Keep functions and methods small and focused on a single responsibility
- Use meaningful names for variables, functions, and classes
- Add comments to explain complex logic or algorithms
- Avoid code duplication by extracting common functionality into reusable functions or components

## Documentation

- Update documentation for any new features or changes to existing functionality
- Ensure all API endpoints are documented
- Add comments to explain complex code
- Keep the README.md up-to-date with any changes to the project structure or setup process

### Documentation Structure

- User documentation goes in `docs/user/`
- Developer documentation goes in `docs/developer/`
- Operator documentation goes in `docs/operator/`

## Reporting Issues

If you find a bug or have a feature request, please create an issue in the GitHub repository.

### Issue Template

When creating an issue, please use the appropriate template:

#### Bug Report

```markdown
## Bug Description
[Provide a clear description of the bug]

## Expected Behavior
[Describe what you expected to happen]

## Actual Behavior
[Describe what actually happened]

## Steps to Reproduce
[Provide steps to reproduce the bug]

## Environment
- OS: [e.g., Ubuntu 20.04]
- Browser: [e.g., Chrome 90.0]
- CMDB Lite Version: [e.g., v1.0.0]

## Additional Context
[Add any other context about the problem here]
```

#### Feature Request

```markdown
## Feature Description
[Provide a clear description of the feature]

## Problem Statement
[Describe the problem this feature would solve]

## Proposed Solution
[Describe your proposed solution]

## Alternatives Considered
[Describe any alternative solutions you've considered]

## Additional Context
[Add any other context or screenshots about the feature request here]
```

## Questions?

If you have any questions about contributing to CMDB Lite, please feel free to:

1. Create an issue with the "question" label
2. Join our community chat (if available)
3. Contact the maintainers directly

Thank you for your interest in contributing to CMDB Lite!