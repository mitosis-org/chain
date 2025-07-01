# Contributing to Mitosis Chain

Welcome to Mitosis Chain! We're excited to have you contribute to the next-generation modular blockchain infrastructure. This guide will help you get started with contributing to our project.

## 🎯 Quick Start

### Prerequisites

- Go 1.21 or higher
- Node.js 18+ (for frontend tools)
- Git
- Docker (for running devnet)

### Development Setup

First, fork the repository to your account or an organization on GitHub.

```bash
# Fork and Clone
git clone https://github.com/YOUR_USERNAME/chain.git
cd chain
git remote add upstream https://github.com/mitosis-org/chain.git

# Verify installation
make build
make test
```

## 🔄 Development Workflow

### 1. Create a Feature Branch
```bash
git checkout -b feature/your-feature-name
```

### 2. Make Changes
- Write your code following our coding standards
- Add tests for new functionality
- Update documentation if needed

### 3. Test Your Changes
```bash
# Run unit tests
make test

# Run linting
make lint
```

### 4. Commit Your Changes
We follow [Conventional Commits](https://www.conventionalcommits.org/) specification:

```bash
git commit -m "feat: add new validator management feature"
git commit -m "fix: resolve consensus layer sync issue"
git commit -m "docs: update API documentation"
git commit -m "test: add unit tests for evmvalidator module"
```

**Commit Types:**
- `feat`: New features
- `fix`: Bug fixes
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Build process or auxiliary tool changes

### 5. Push and Create Pull Request
```bash
git push origin feature/your-feature-name
```

Then create a Pull Request through GitHub with:
- Clear title and description
- Reference any related issues
- Include screenshots/demos if applicable

## 📝 Coding Standards

### Go Code Style
- Follow the [official Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use [golangci-lint](https://github.com/golangci/golangci-lint) for comprehensive linting
- Use `gofmt` and `goimports` for formatting
- Write meaningful variable and function names
- Add comments for exported functions and complex logic

### Solidity Code Style
- Follow our [Solidity Style Guide](.cursor/rules/shared-rules/solidity/reference/coinbase-style-guide.mdc)
- Use ERC7201 for storage patterns
- Implement comprehensive tests for all contracts

### Project Structure
```
chain/
├── app/           # Cosmos SDK app configuration
├── cmd/           # CLI applications
├── x/             # Custom Cosmos SDK modules
├── types/         # Common types
├── proto/         # Protocol buffer definitions
├── infra/         # Infrastructure and deployment scripts for testing environments
└── scripts/       # Build and utility scripts
```

## 🤝 Community Guidelines

- Be respectful and inclusive
- Help newcomers get started
- Share knowledge and best practices
- Follow our [Code of Conduct](CODE_OF_CONDUCT.md)

## 🆘 Getting Help

- **GitHub Discussions**: Technical questions and ideas
- **Discord**: Real-time community chat
- **Issues**: Bug reports and feature requests

## 📚 Additional Resources

- [Mitosis Documentation](https://docs.mitosis.org/)
- [Cosmos SDK Documentation](https://docs.cosmos.network/)

---

Thank you for contributing to Mitosis Chain! Your efforts help build the future of modular blockchain infrastructure. 🚀