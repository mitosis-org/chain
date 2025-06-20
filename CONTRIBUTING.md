# Contributing to Mitosis Chain

Welcome to Mitosis Chain! We're excited to have you contribute to the next-generation modular blockchain infrastructure. This guide will help you get started with contributing to our project.

## 🎯 Quick Start

### Prerequisites

- Go 1.21 or higher
- Node.js 18+ (for frontend tools)
- Git
- Docker (for running devnet)

### Development Setup

1. **Fork and Clone**
   ```bash
   git clone https://github.com/YOUR_USERNAME/chain.git
   cd chain
   git remote add upstream https://github.com/mitosis-org/chain.git
   ```

2. **Setup Development Environment**
   ```bash
   # Fetch submodules
   git submodule update --init --recursive
   
   # Setup localnet (for testing)
   make setup-geth
   make setup-mitosisd
   ```

3. **Verify Installation**
   ```bash
   make test
   make build
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

# Run integration tests
make test-integration

# Run linting
make lint

# Test localnet setup
make localnet-test
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
├── infra/         # Infrastructure and deployment
└── scripts/       # Build and utility scripts
```

## 🧪 Testing Guidelines

### Unit Tests
- Write tests for all new functions
- Aim for >80% code coverage
- Use table-driven tests where appropriate
- Mock external dependencies

### Integration Tests
- Test module interactions
- Verify end-to-end workflows
- Include both happy path and error cases

### Example Test Structure
```go
func TestValidatorCreation(t *testing.T) {
    tests := []struct {
        name    string
        input   ValidatorInput
        want    Validator
        wantErr bool
    }{
        {
            name: "valid validator creation",
            input: ValidatorInput{...},
            want: Validator{...},
            wantErr: false,
        },
        // Add more test cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

## 🐛 Bug Reports

When filing a bug report, please include:

1. **Environment Information**
   - Operating system and version
   - Go version
   - Mitosis Chain version

2. **Steps to Reproduce**
   - Clear, numbered steps
   - Expected vs actual behavior
   - Relevant logs or error messages

3. **Additional Context**
   - Screenshots if applicable
   - Configuration files
   - Related issues or PRs

## 💡 Feature Requests

For new features:

1. **Check existing issues** to avoid duplicates
2. **Describe the use case** and problem you're solving
3. **Propose a solution** with technical details
4. **Consider breaking changes** and backward compatibility

## 📖 Documentation

- Update relevant documentation with your changes
- Follow markdown best practices
- Include code examples where helpful
- Update the changelog for significant changes

### Documentation Types
- **README.md**: Project overview and quick start
- **API Documentation**: Generated from code comments
- **User Guides**: Step-by-step instructions
- **Developer Guides**: Technical implementation details

## 🔍 Code Review Process

### For Contributors
- Respond to review feedback promptly
- Make requested changes in new commits
- Keep discussions focused and constructive

### For Reviewers
- Be constructive and respectful
- Focus on code quality, security, and maintainability
- Approve when satisfied with the changes

## 🏆 Recognition

Contributors are recognized through:
- **Contributor list** in our README
- **Release notes** acknowledgments
- **Community highlights** in our social channels

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

- [Cosmos SDK Documentation](https://docs.cosmos.network/)
- [Go Documentation](https://golang.org/doc/)
- [Mitosis Chain Documentation](https://docs.mitosis.org/)

---

Thank you for contributing to Mitosis Chain! Your efforts help build the future of modular blockchain infrastructure. 🚀 