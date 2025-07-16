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
- `perf`: Performance improvements

### 5. Push and Create Pull Request
```bash
git push origin feature/your-feature-name
```
Then follow these steps to create a Pull Request through GitHub:
- Push your branch: The above command uploads your changes to your forked repository.
- Open a Pull Request: Go to the GitHub repository, click "Pull requests," then "New Pull Request." Select your branch (`feature/your-feature-name`) as the compare branch.
- Fill out the form: Use a clear title (e.g., `docs: improve PR instructions`) and description. Reference related issues (e.g., `Fixes #123`) if applicable.
- Add visuals (optional): Include screenshots or demos if your changes affect the UI or require visual explanation.
- Submit: Click "Create pull request" to submit for review.

## 📝 Coding Standards

### Go Code Style
- Follow the [official Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use [golangci-lint](https://github.com/golangci/golangci-lint) for comprehensive linting
- Use `gofmt` and `goimports` for formatting
- Write meaningful variable and function names
- Add comments for exported functions and complex logic

### Solidity Code Style
- Follow our [Solidity Style Guide](https://github.com/mitosis-org/shared-rules/blob/main/solidity/reference/coinbase-style-guide.mdc)
- Use ERC-7201 for storage patterns
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

- **Be respectful and inclusive**: Treat all community members with kindness and ensure an welcoming environment for everyone.
- **Help newcomers get started**: Offer guidance to new contributors to help them navigate the project.
- **Share knowledge and best practices**: Contribute tips or examples to improve the community's collective expertise.
- **Follow our [Code of Conduct](CODE_OF_CONDUCT.md)**: Adhere to the guidelines outlined to maintain a positive community culture.

## 🆘 Getting Help

- **GitHub Discussions**: Technical questions and ideas
- **Discord**: Real-time community chat
- **Issues**: Bug reports and feature requests

## 📚 Additional Resources

- [Mitosis Developer Docs](https://docs.mitosis.org/developers/)
- [Cosmos SDK Documentation](https://docs.cosmos.network/)

---

Thank you for contributing to Mitosis Chain! Your efforts help build the future of modular blockchain infrastructure. 🚀
