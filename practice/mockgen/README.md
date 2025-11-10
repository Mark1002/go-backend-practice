# Go Mock Practice with gomock

This folder contains practical examples demonstrating the core concepts and usage patterns of **gomock** - a powerful mocking framework for Go unit testing.

## 📋 Table of Contents

- [Overview](#overview)
- [Core Concepts](#core-concepts)
- [Project Structure](#project-structure)
- [Installation & Setup](#installation--setup)
- [Understanding the Code](#understanding-the-code)
- [Mock Generation](#mock-generation)
- [Testing Patterns](#testing-patterns)
- [Advanced Features](#advanced-features)
- [Best Practices](#best-practices)
- [Common Pitfalls](#common-pitfalls)
- [References](#references)

## 🎯 Overview

**gomock** is a mocking framework that helps you create mock objects for testing in Go. It allows you to:

- **Isolate** the code under test from external dependencies
- **Control** the behavior of dependencies during testing
- **Verify** that interactions with dependencies happen as expected
- **Simulate** various scenarios including error conditions

### Why Use Mocking?

Mocking is essential when testing components that depend on:
- Databases
- External APIs
- File systems
- Network services
- Complex business logic components

Instead of using real dependencies (which can be slow, unreliable, or unavailable), mocks provide controlled, predictable behavior.

## 🧠 Core Concepts

### 1. Interfaces and Dependency Injection

Go's interface system is fundamental to effective mocking. In our example:

```go
type IUserRepo interface {
    GetUserByID(id int) (*User, error)
    Insert(user User) error
    Update(id int, user User) error
    Delete(id int) error
}
```

The `UserService` depends on this interface, not a concrete implementation:

```go
type UserService struct {
    repo IUserRepo  // Dependency injection through interface
}
```

### 2. Mock Object Architecture

gomock generates mock objects with three key components:

#### Controller (`*gomock.Controller`)
- **Purpose**: Manages the lifecycle of all mock objects
- **Responsibilities**: 
  - Tracks expected method calls
  - Validates that expectations are met
  - Reports test failures
  - Coordinates with Go's testing framework

#### Mock Object (`MockIUserRepo`)
- **Purpose**: Implements the interface being mocked
- **Structure**:
  ```go
  type MockIUserRepo struct {
      ctrl     *gomock.Controller  // Reference to controller
      recorder *MockIUserRepoMockRecorder  // Sets up expectations
      isgomock struct{}  // Marker for gomock objects
  }
  ```

#### Recorder (`MockIUserRepoMockRecorder`)
- **Purpose**: Bridge between test code and controller
- **Function**: When you call `mock.EXPECT().Method()`, you're actually calling the recorder, which registers expectations with the controller

### 3. Test Lifecycle

```go
func TestExample(t *testing.T) {
    // 1. Create controller (manages mock lifecycle)
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()  // Optional in Go 1.14+
    
    // 2. Create mock object
    mockRepo := NewMockIUserRepo(ctrl)
    
    // 3. Set expectations (what should happen)
    mockRepo.EXPECT().GetUserByID(1).Return(&User{ID: 1}, nil)
    
    // 4. Execute code under test
    service := UserService{repo: mockRepo}
    user, err := service.GetUserByID(1)
    
    // 5. Verify results
    assert.NoError(t, err)
    assert.Equal(t, 1, user.ID)
    
    // 6. Controller automatically verifies all expectations were met
}
```

## 📁 Project Structure

```
mockgen/
├── domain.go          # Domain entities (User struct)
├── repo.go           # Repository interface definition
├── service.go        # Business logic (UserService)
├── mock_repo.go      # Generated mock (DO NOT EDIT)
├── service_test.go   # Test examples using mocks
└── README.md         # This file
```

## ⚙️ Installation & Setup

### 1. Install gomock

```bash
go get -u go.uber.org/mock/gomock
go install go.uber.org/mock/mockgen@latest
```

### 2. Verify Installation

```bash
mockgen -version
```

### 3. Add to Project

```go
// In your go.mod
require (
    go.uber.org/mock v0.4.0
    github.com/stretchr/testify v1.8.4
)
```

## 🔍 Understanding the Code

### Domain Layer (`domain.go`)
```go
type User struct {
    ID   int
    Name string
}
```
Simple domain entity representing a user.

### Repository Interface (`repo.go`)
```go
type IUserRepo interface {
    GetUserByID(id int) (*User, error)
    Insert(user User) error
    Update(id int, user User) error
    Delete(id int) error
}
```
Defines the contract for data access operations.

### Business Logic (`service.go`)
```go
type UserService struct {
    repo IUserRepo  // Dependency injection
}

func (u *UserService) Upsert(user User) error {
    if user.ID <= 0 {
        return invalidUserIDError
    }
    
    // Check if user exists
    existingUser, err := u.repo.GetUserByID(user.ID)
    if err != nil {
        return err
    }
    
    // Insert new user or update existing
    if existingUser == nil {
        return u.repo.Insert(user)
    }
    return u.repo.Update(user.ID, user)
}
```

The `Upsert` method demonstrates complex business logic that requires multiple repository calls.

## 🏭 Mock Generation

### Generate Mock

```bash
# From project root
mockgen -source=practice/mockgen/repo.go \
        -destination=practice/mockgen/mock_repo.go \
        -package=mockgen
```

### Understanding Generated Code

The generated mock implements the interface:

```go
// Generated mock structure
type MockIUserRepo struct {
    ctrl     *gomock.Controller
    recorder *MockIUserRepoMockRecorder
    isgomock struct{}
}

// Method implementation
func (m *MockIUserRepo) GetUserByID(id int) (*User, error) {
    m.ctrl.T.Helper()  // Skip this frame in error traces
    ret := m.ctrl.Call(m, "GetUserByID", id)  // Register the call
    ret0, _ := ret[0].(*User)
    ret1, _ := ret[1].(error)
    return ret0, ret1
}

// Expectation recorder
func (mr *MockIUserRepoMockRecorder) GetUserByID(id interface{}) *gomock.Call {
    mr.mock.ctrl.T.Helper()
    return mr.mock.ctrl.RecordCallWithMethodType(
        mr.mock, "GetUserByID", 
        reflect.TypeOf((*MockIUserRepo)(nil).GetUserByID), id)
}
```

## 🧪 Testing Patterns

### 1. Basic Expectation Setting

```go
// Expect specific input and return specific output
mockRepo.EXPECT().GetUserByID(1).Return(&User{ID: 1, Name: "John"}, nil)

// Expect method to be called once (default)
mockRepo.EXPECT().Insert(gomock.Any()).Times(1)

// Expect method to be called multiple times
mockRepo.EXPECT().Delete(gomock.Any()).Times(3)

// Allow method to be called any number of times
mockRepo.EXPECT().GetUserByID(gomock.Any()).AnyTimes()
```

### 2. Parameter Matching

```go
import "go.uber.org/mock/gomock"

// Exact match
mockRepo.EXPECT().GetUserByID(1)

// Any value
mockRepo.EXPECT().GetUserByID(gomock.Any())

// Conditional matching
mockRepo.EXPECT().Insert(gomock.Cond(func(u User) bool {
    return u.ID > 0 && u.Name != ""
}))

// Equality check
mockRepo.EXPECT().Update(gomock.Eq(1), gomock.Any())
```

### 3. Error Simulation

```go
// Simulate database error
mockRepo.EXPECT().
    GetUserByID(999).
    Return(nil, errors.New("user not found"))

// Simulate network timeout
mockRepo.EXPECT().
    Insert(gomock.Any()).
    Return(context.DeadlineExceeded)
```

### 4. Complex Scenarios with Multiple Calls

```go
func TestUpsertNewUser(t *testing.T) {
    ctrl := gomock.NewController(t)
    mockRepo := NewMockIUserRepo(ctrl)
    
    user := User{ID: 1, Name: "John"}
    
    // Set up expectation chain for upsert of new user
    mockRepo.EXPECT().
        GetUserByID(1).
        Return(nil, nil).  // User doesn't exist
        Times(1)
    
    mockRepo.EXPECT().
        Insert(user).
        Return(nil).  // Successful insert
        Times(1)
    
    // Execute test
    service := UserService{repo: mockRepo}
    err := service.Upsert(user)
    
    // Verify
    assert.NoError(t, err)
}
```

### 5. Table-Driven Tests with Mocks

```go
func TestUpsertUser(t *testing.T) {
    tests := []struct {
        name                 string
        user                 User
        specifyFunctionCalls func(mock *MockIUserRepo)
        expectedError        error
    }{
        {
            name: "Should insert new user",
            user: User{ID: 1, Name: "John"},
            specifyFunctionCalls: func(mock *MockIUserRepo) {
                mock.EXPECT().GetUserByID(1).Return(nil, nil).Times(1)
                mock.EXPECT().Insert(User{ID: 1, Name: "John"}).Return(nil).Times(1)
            },
            expectedError: nil,
        },
        {
            name: "Should update existing user",
            user: User{ID: 1, Name: "John Updated"},
            specifyFunctionCalls: func(mock *MockIUserRepo) {
                mock.EXPECT().GetUserByID(1).Return(&User{ID: 1, Name: "John"}, nil).Times(1)
                mock.EXPECT().Update(1, User{ID: 1, Name: "John Updated"}).Return(nil).Times(1)
            },
            expectedError: nil,
        },
        {
            name:          "Should reject invalid user ID",
            user:          User{ID: -1, Name: "Invalid"},
            expectedError: invalidUserIDError,
        },
    }

    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            ctrl := gomock.NewController(t)
            mockRepo := NewMockIUserRepo(ctrl)
            
            if test.specifyFunctionCalls != nil {
                test.specifyFunctionCalls(mockRepo)
            }
            
            service := UserService{repo: mockRepo}
            err := service.Upsert(test.user)
            
            assert.Equal(t, test.expectedError, err)
        })
    }
}
```

## 🚀 Advanced Features

### 1. Callback Functions with `Do()`

```go
// Execute custom logic when method is called
mockRepo.EXPECT().
    Delete(1).
    Do(func(id int) {
        t.Logf("Deleting user with ID: %d", id)
    }).
    Return(nil)
```

### 2. Call Ordering with `InOrder()`

```go
// Ensure methods are called in specific order
gomock.InOrder(
    mockRepo.EXPECT().GetUserByID(1).Return(&User{ID: 1}, nil),
    mockRepo.EXPECT().Delete(1).Return(nil),
)
```

### 3. Dynamic Return Values with `DoAndReturn()`

```go
// Calculate return value dynamically
mockRepo.EXPECT().
    GetUserByID(gomock.Any()).
    DoAndReturn(func(id int) (*User, error) {
        if id <= 0 {
            return nil, errors.New("invalid ID")
        }
        return &User{ID: id, Name: fmt.Sprintf("User %d", id)}, nil
    }).
    AnyTimes()
```

### 4. Testing Concurrent Code

```go
func TestConcurrentAccess(t *testing.T) {
    ctrl := gomock.NewController(t)
    mockRepo := NewMockIUserRepo(ctrl)
    
    // Setup expectations for concurrent calls
    mockRepo.EXPECT().
        GetUserByID(gomock.Any()).
        Return(&User{}, nil).
        Times(10)  // Expect 10 concurrent calls
    
    service := UserService{repo: mockRepo}
    
    var wg sync.WaitGroup
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            service.GetUserByID(id)
        }(i)
    }
    
    wg.Wait()
}
```

## ✅ Best Practices

### 1. **One Controller Per Test**
```go
func TestSomething(t *testing.T) {
    ctrl := gomock.NewController(t)  // New controller for each test
    // ... test logic
}
```

### 2. **Use Dependency Injection**
```go
// Good: Accept interface
func NewUserService(repo IUserRepo) *UserService {
    return &UserService{repo: repo}
}

// Avoid: Hard-coded dependencies
func NewUserService() *UserService {
    return &UserService{repo: &MySQLUserRepo{}}  // Hard to test
}
```

### 3. **Test Behavior, Not Implementation**
```go
// Good: Test the outcome
func TestDeleteUser(t *testing.T) {
    // Setup mock to return success
    mockRepo.EXPECT().Delete(1).Return(nil)
    
    err := service.DeleteUserByID(1)
    assert.NoError(t, err)  // Verify the result
}

// Avoid: Over-specifying internal calls
func TestDeleteUser(t *testing.T) {
    // Too specific about internal implementation
    mockRepo.EXPECT().GetUserByID(1).Return(&User{}, nil)
    mockRepo.EXPECT().ValidateUser(gomock.Any()).Return(true)
    mockRepo.EXPECT().LogDeletion(1)
    mockRepo.EXPECT().Delete(1).Return(nil)
    // ... too many internal details
}
```

### 4. **Use Meaningful Test Names**
```go
func TestUpsert_WhenUserExists_ShouldUpdate(t *testing.T) { }
func TestUpsert_WhenUserDoesNotExist_ShouldInsert(t *testing.T) { }
func TestUpsert_WhenInvalidID_ShouldReturnError(t *testing.T) { }
```

### 5. **Group Related Expectations**
```go
// Good: Group related setup
func setupNewUserScenario(mock *MockIUserRepo, user User) {
    mock.EXPECT().GetUserByID(user.ID).Return(nil, nil)
    mock.EXPECT().Insert(user).Return(nil)
}

// Usage
setupNewUserScenario(mockRepo, user)
```

## ⚠️ Common Pitfalls

### 1. **Forgetting to Set Expectations**
```go
// Wrong: No expectations set
func TestGetUser(t *testing.T) {
    ctrl := gomock.NewController(t)
    mockRepo := NewMockIUserRepo(ctrl)
    
    service := UserService{repo: mockRepo}
    user, err := service.GetUserByID(1)  // Will fail: unexpected call
}

// Correct: Set expectations
func TestGetUser(t *testing.T) {
    ctrl := gomock.NewController(t)
    mockRepo := NewMockIUserRepo(ctrl)
    
    mockRepo.EXPECT().GetUserByID(1).Return(&User{ID: 1}, nil)  // Set expectation
    
    service := UserService{repo: mockRepo}
    user, err := service.GetUserByID(1)
}
```

### 2. **Expectations Not Called**
```go
// Wrong: Expectation set but method never called
func TestSomething(t *testing.T) {
    ctrl := gomock.NewController(t)
    mockRepo := NewMockIUserRepo(ctrl)
    
    mockRepo.EXPECT().GetUserByID(1).Return(&User{}, nil)  // Expectation set
    
    // But GetUserByID is never actually called!
    // Test will fail when ctrl.Finish() is called
}
```

### 3. **Incorrect Parameter Matching**
```go
// Wrong: Parameter mismatch
mockRepo.EXPECT().GetUserByID(1).Return(&User{}, nil)
service.GetUserByID(2)  // Called with 2, expected 1 - will fail

// Correct: Match actual parameters
mockRepo.EXPECT().GetUserByID(2).Return(&User{}, nil)
service.GetUserByID(2)
```

### 4. **Over-Mocking**
```go
// Avoid: Mocking everything unnecessarily
func TestSimpleValidation(t *testing.T) {
    // Don't mock simple validation logic that doesn't need external dependencies
    user := User{ID: -1}
    service := UserService{}  // No repo needed for validation
    
    err := service.validateUser(user)  // Just test the validation logic
    assert.Error(t, err)
}
```

### 5. **Lifecycle Misalignment with Test Suites**

When using test suites (e.g., `testify/suite`), ensure controller lifecycle aligns with test lifecycle:

```go
// Wrong: Controller created once for entire suite
type UserSuite struct {
    suite.Suite
    ctrl *gomock.Controller
    mockRepo *MockIUserRepo
}

func (s *UserSuite) SetupSuite() {
    s.ctrl = gomock.NewController(s.T())  // Lives for entire suite
    s.mockRepo = NewMockIUserRepo(s.ctrl)
}

// Correct: Controller per test
func (s *UserSuite) SetupTest() {
    s.ctrl = gomock.NewController(s.T())  // New controller per test
    s.mockRepo = NewMockIUserRepo(s.ctrl)
}

func (s *UserSuite) TearDownTest() {
    s.ctrl.Finish()  // Clean up after each test
}
```

## 🔗 References

- [gomock Official Documentation](https://pkg.go.dev/go.uber.org/mock/gomock)
- [Mastering Mocking in Go: Comprehensive Guide](https://towardsdev.com/mastering-mocking-in-go-comprehensive-guide-to-gomock-with-practical-examples-e12c1773842f)
- [Understanding gomock Architecture and Lifecycle](https://blog.kenwsc.com/gomock-architecture-and-lifecycle/)
- [Go Testing Package](https://pkg.go.dev/testing)
- [testify/assert Documentation](https://pkg.go.dev/github.com/stretchr/testify/assert)

---

## 🎯 Key Takeaways

1. **gomock excels at isolating business logic** from external dependencies
2. **Interfaces are crucial** for effective mocking in Go
3. **Controller manages the lifecycle** of all mock objects and expectations
4. **Recorder bridges** test expectations with the controller
5. **Test behavior, not implementation details** to maintain maintainable tests
6. **One controller per test** ensures proper isolation
7. **Use dependency injection** to make code testable

Understanding these concepts will help you write more effective, maintainable unit tests in Go! 🚀