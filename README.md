# Validation Library

This is a simple validation library in Go that uses `go-playground/validator/v10` for validating structs. It supports defining custom validation rules programmatically and using struct tags for field names, making it suitable for REST API validation.

## Installation

```sh
go get github.com/umbeluzi/validation
```

## Usage

### Basic Usage

To use the validation library, first, create a struct that you want to validate and then define validation rules for its fields.

```go
package main

import (
    "fmt"
    "github.com/umbeluzi/validation"
)

type User struct {
    Name  string `json:"name"`
    Email string `json:"email"`
    Age   string `json:"age"`
}

func main() {
    v := validation.New()

    user := User{
        Name:  "JohnDoe",
        Email: "john.doe@example.com",
        Age:   "30",
    }

    rules := []validation.Rule{
        {
            Field:    "Name",
            Tag:      "required,alpha",
            Message:  "Name must contain only letters",
        },
        {
            Field:    "Email",
            Tag:      "required,email",
            Message:  "Email must be a valid email address",
        },
        {
            Field:    "Age",
            Tag:      "required,numeric",
            Message:  "Age must be a number",
        },
    }

    if errs := v.ValidateStruct(user, rules); errs != nil {
        fmt.Println("Validation failed:", errs)
    } else {
        fmt.Println("Validation succeeded")
    }
}
```

### Custom Validators

You can register custom validators using the `Register` method. The custom validation function should implement the `ValidatorFunc` interface.

```go
package main

import (
    "fmt"
    "github.com/umbeluzi/validation"
)

type User struct {
    Name  string `json:"name"`
    Email string `json:"email"`
    Age   string `json:"age"`
}

func main() {
    v := validation.New()

    // Register a custom validation function
    v.Register("is-even", func(fl validation.FieldLevel) bool {
        value := fl.Field().(string)
        if len(value) == 0 {
            return false
        }
        // Check if the string length is even
        return len(value)%2 == 0
    })

    user := User{
        Name:  "JohnDoe",
        Email: "john.doe@example.com",
        Age:   "30",
    }

    rules := []validation.Rule{
        {
            Field:    "Name",
            Tag:      "required,alpha",
            Message:  "Name must contain only letters",
        },
        {
            Field:    "Email",
            Tag:      "required,email",
            Message:  "Email must be a valid email address",
        },
        {
            Field:    "Age",
            Tag:      "required,numeric,is-even",
            Message:  "Age must be an even number",
        },
    }

    if errs := v.ValidateStruct(user, rules); errs != nil {
        fmt.Println("Validation failed:", errs)
    } else {
        fmt.Println("Validation succeeded")
    }
}
```

### Custom Error Messages

You can customize error messages for specific validation tags using the `SetMessage` method.

```go
package main

import (
    "fmt"
    "github.com/umbeluzi/validation"
)

type User struct {
    Name  string `json:"name"`
    Email string `json:"email"`
    Age   string `json:"age"`
}

func main() {
    v := validation.New()

    user := User{
        Name:  "JohnDoe",
        Email: "john.doe@example.com",
        Age:   "30",
    }

    rules := []validation.Rule{
        {
            Field:    "Name",
            Tag:      "required,alpha",
            Message:  "Name must contain only letters",
        },
        {
            Field:    "Email",
            Tag:      "required,email",
            Message:  "Email must be a valid email address",
        },
        {
            Field:    "Age",
            Tag:      "required,numeric",
            Message:  "Age must be a number",
        },
    }

    if errs := v.ValidateStruct(user, rules); errs != nil {
        // Customize error messages
        errs.SetMessage("Age", "Age must be a valid integer")
        fmt.Println("Validation failed:", errs)
    } else {
        fmt.Println("Validation succeeded")
    }
}
```

## Testing

Run the tests using:

```sh
go test ./...
```

**Example Test**

```go
package validation

import (
    "reflect"
    "testing"
)

type User struct {
    Name  string `json:"name"`
    Email string `json:"email"`
    Age   string `json:"age"`
}

func TestValidateStruct(t *testing.T) {
    v := New()

    validUser := User{
        Name:  "JohnDoe",
        Email: "john.doe@example.com",
        Age:   "30",
    }

    rules := []Rule{
        {
            Field:    "Name",
            Tag:      "required,alpha",
            Message:  "Name must contain only letters",
        },
        {
            Field:    "Email",
            Tag:      "required,email",
            Message:  "Email must be a valid email address",
        },
        {
            Field:    "Age",
            Tag:      "required,numeric",
            Message:  "Age must be a number",
        },
    }

    if errs := v.ValidateStruct(validUser, rules); errs != nil {
        t.Errorf("expected no error, got %v", errs)
    }

    invalidUser := User{
        Name:  "John123",
        Email: "invalid-email",
        Age:   "thirty",
    }

    expectedErrors := Errors{
        {Field: "name", Message: "Name must contain only letters"},
        {Field: "email", Message: "Email must be a valid email address"},
        {Field: "age", Message: "Age must be a number"},
    }

    errs := v.ValidateStruct(invalidUser, rules)
    if !reflect.DeepEqual(errs, expectedErrors) {
        t.Errorf("expected %v, got %v", expectedErrors, errs)
    }
}
