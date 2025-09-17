# namedreturns

This linter enforces the use of named returns in Go functions. Named returns improve code readability and make function signatures more self-documenting.

> **Note**: This linter was proposed for inclusion in golangci-lint but was ultimately rejected. See [golangci-lint PR #6083](https://github.com/golangci/golangci-lint/pull/6083) for the discussion. The maintainers cited existing linters as duplicates, but this linter serves a different purpose:
>
> - **nonamedreturns** (firefart/nonamedreturns): Flags named returns as bad practice - the opposite philosophy
> - **gocritic (unnamedResult)**: Suggests adding names but doesn't enforce or validate them
> - **revive (bare-return)**: Warns against bare returns but doesn't ensure proper named return usage
>
> This linter uniquely enforces consistent use of named returns and catches shadowed variables and signature mismatches. As a result, this standalone version provides an easy way to use the linter outside of the golangci-lint ecosystem.

## Installation and Usage

Since it was rejected for inclusion in [https://github.com/golangci/golangci-lint](https://github.com/golangci/golangci-lint), we have to get creative.

### Option 1: Install via go install (Recommended)
```bash
go install github.com/nikogura/namedreturns@latest
namedreturns ./...
```

### Option 2: Run directly with go run
```bash
go run github.com/nikogura/namedreturns@latest ./...
```

### Option 3: Build and run locally
```bash
make build
./namedreturns ./...
```

### Option 4: Use Makefile targets
```bash
# Build and run on the project itself
make lint-self

# Or just build the binary
make build
```

## About

Tutorial on how to write your own linter:
https://disaev.me/p/writing-useful-go-analysis-linter/

Named errors used in defers are not reported. If you also want to report them set `report-error-in-defer` to true.

Inspired by, and derived from the excellent https://github.com/firefart/nonamedreturns .  I respect the author's opinion, but disagree with him on every point.

This linter will also catch "Shadowed Variables", and cases where the function signature claims return values are named, but the function returns something else.

I am forever inspired by the Go Proverb "Clear is better than clever."  Named returns, in my opinion, make your codebase clearer for the next person to come along.  Who knows?  That might even be you, 6 months later, when all the context has leaked out of your brain, when you're paged at 03:00 locally and people are screaming.

I'm always trying to do my "future self" a favor.

# Why are named returns beneficial?

1. **Improved Readability**
   Named returns make function signatures more self-documenting by clearly indicating what values the function returns. This is especially helpful for functions with multiple return values.

2. **Better Documentation**
   Named returns serve as inline documentation, making it clear what each return value represents without having to look at the function body or implementation details.

3. **Consistent API Design**
   Named returns encourage developers to think carefully about what their functions return and provide meaningful names for those return values.

4. **Easier Maintenance**
   When refactoring or modifying functions, named returns make it easier to understand the impact of changes on return values.

5. **Clear Intent**
   Named returns make the developer's intent explicit about what the function is designed to return, reducing ambiguity in the codebase.

6. **Better Error Handling**
   For error returns, named variables allow for more flexible error handling patterns, especially when combined with defer statements.

#### Example
```golang
// Good - named returns
func processUser(id string) (user *User, err error) {
    user, err = fetchUser(id)
    if err != nil {
        return // returns user and err as defined in signature
    }
    return
}

// Bad - unnamed returns
func processUser(id string) (*User, error) {
    user, err := fetchUser(id)
    if err != nil {
        return nil, err // less readable, requires looking at implementation
    }
    return user, nil
}
```
