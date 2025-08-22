# Contributing to GoNest

Thank you for your interest in contributing to GoNest! This document provides guidelines and information for contributors.

## 🤝 How to Contribute

### Reporting Issues

Before creating an issue, please:

1. **Search existing issues** to see if your problem has already been reported
2. **Check the documentation** to see if your question is answered there
3. **Provide detailed information** including:
   - Go version
   - Operating system
   - Steps to reproduce
   - Expected vs actual behavior
   - Code examples

### Feature Requests

We welcome feature requests! Please:

1. **Describe the feature** clearly and concisely
2. **Explain the use case** and why it would be valuable
3. **Provide examples** of how you would use it
4. **Consider implementation** and discuss potential approaches

### Pull Requests

We love pull requests! Here's how to contribute:

#### Before You Start

1. **Fork the repository**
2. **Create a feature branch** from `main`
3. **Check existing issues** to see if your work is already in progress

#### Development Guidelines

1. **Follow Go conventions**:
   - Use `gofmt` for formatting
   - Follow [Effective Go](https://golang.org/doc/effective_go.html)
   - Use meaningful variable and function names

2. **Write tests**:
   - Add tests for new features
   - Ensure existing tests pass
   - Aim for good test coverage

3. **Update documentation**:
   - Update relevant documentation
   - Add examples for new features
   - Update API documentation if needed

4. **Commit messages**:
   - Use clear, descriptive commit messages
   - Follow conventional commit format: `type(scope): description`
   - Examples:
     - `feat(auth): add JWT refresh token support`
     - `fix(validation): resolve struct tag parsing issue`
     - `docs(api): update controller documentation`

#### Code Style

```go
// Good: Clear function name and documentation
// CreateUser creates a new user with the provided data
func (s *UserService) CreateUser(user *User) error {
    if user == nil {
        return gonest.BadRequestException("User data is required")
    }
    
    // Validate user data
    if err := s.validateUser(user); err != nil {
        return err
    }
    
    // Create user
    return s.repository.Create(user)
}

// Good: Clear error handling
func (c *UserController) GetUser(ctx echo.Context) error {
    id := ctx.Param("id")
    if id == "" {
        return gonest.BadRequestException("User ID is required")
    }
    
    user, err := c.userService.GetUser(id)
    if err != nil {
        return err
    }
    
    return ctx.JSON(http.StatusOK, user)
}
```

#### Testing

```go
func TestUserService_CreateUser(t *testing.T) {
    // Arrange
    service := NewUserService()
    user := &User{
        Name:  "John Doe",
        Email: "john@example.com",
    }
    
    // Act
    err := service.CreateUser(user)
    
    // Assert
    assert.NoError(t, err)
    assert.NotEmpty(t, user.ID)
}

func TestUserController_GetUser(t *testing.T) {
    // Arrange
    testApp := gonest.NewTestApp().
        WithModule(userModule).
        Build()
    defer testApp.Stop()
    
    // Act & Assert
    response := testApp.Request("GET", "/users/123").
        ExpectStatus(200).
        Get()
    
    assert.NotNil(t, response.JSON())
}
```

### Review Process

1. **Automated checks** must pass:
   - Tests
   - Linting
   - Code formatting

2. **Code review** by maintainers:
   - We'll review your code for quality and correctness
   - We may suggest improvements or ask questions
   - We'll help ensure your contribution fits well with the project

3. **Merge** once approved:
   - Your PR will be merged into the main branch
   - You'll be added to the contributors list

## 🏗️ Development Setup

### Prerequisites

- Go 1.21 or higher
- Git

### Local Development

1. **Clone the repository**:
   ```bash
   git clone https://github.com/ulims/GoNest.git
   cd gonest
   ```

2. **Install dependencies**:
   ```bash
   go mod tidy
   ```

3. **Run tests**:
   ```bash
   go test ./...
   ```

4. **Run examples**:
   ```bash
   go run examples/basic/main.go
   go run examples/advanced/main.go
   go run examples/mongodb/main.go
   ```

### Building the CLI

```bash
go build -o bin/gonest cmd/gonest/main.go
```

## 📋 Issue Templates

### Bug Report

```markdown
## Bug Description
Brief description of the bug

## Steps to Reproduce
1. Step 1
2. Step 2
3. Step 3

## Expected Behavior
What you expected to happen

## Actual Behavior
What actually happened

## Environment
- Go version: 
- OS: 
- GoNest version: 

## Additional Information
Any other context, logs, or screenshots
```

### Feature Request

```markdown
## Feature Description
Brief description of the feature

## Use Case
Why this feature would be valuable

## Proposed Implementation
How you think this could be implemented

## Examples
Code examples of how you would use this feature

## Additional Information
Any other context or considerations
```

## 🎯 Areas for Contribution

We're always looking for help in these areas:

### High Priority
- **Performance improvements**
- **Bug fixes**
- **Documentation improvements**
- **Test coverage**

### Medium Priority
- **New features**
- **Examples and tutorials**
- **CLI improvements**
- **Integration examples**

### Low Priority
- **Code style improvements**
- **Minor optimizations**
- **Additional examples**

## 📞 Getting Help

- **Documentation**: Check the [main documentation](docs/DOCUMENTATION.md)
- **Issues**: Search existing [issues](https://github.com/ulims/GoNest/issues)
- **Discussions**: Join [discussions](https://github.com/ulims/GoNest/discussions)
- **Email**: Contact maintainers directly

## 🙏 Recognition

Contributors will be:

- **Listed in the contributors section** of the README
- **Mentioned in release notes** for significant contributions
- **Invited to join the maintainer team** for consistent contributors

## 📄 License

By contributing to GoNest, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing to GoNest! 🚀
