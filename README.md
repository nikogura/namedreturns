# namedreturns

This linter enforces the use of named returns in Go functions. Named returns improve code readability and make function signatures more self-documenting.

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
        return user, err // returns user and err as defined in signature
    }
    return user, err  // clearly returns what was promised in the signature.
}

// Bad - unnamed returns
func processUser(id string) (*User, error) {
    user, err := fetchUser(id)
    if err != nil {
        return nil, err // less readable, requires looking at implementation
    }
    return user, nil
}

// Bad - Returning something other than what was promised
func processUser(id string) (user *User, err error) {
    myUser, err := fetchUser(id)
    if err != nil {
        return // returns something like what was promised.
    }

    return myUser, nil  // return values might be equivalent to what was promised, but in a long complicated function there could also be surprises.
}
```
## Go Version Compatibility

namedreturns supports analyzing codebases using **Go 1.21.0 and later**. The linter binary can be built with any Go version >= 1.21.0.

### Compatibility with Newer Go Versions

To analyze codebases using newer Go versions than the linter was built with:

```bash
# Simple rebuild with current Go version
make rebuild

# Or manually:
go build -o namedreturns .
```

**Why this works:** namedreturns uses only stable Go AST analysis APIs that are forward-compatible across Go versions.

### Version Strategy

- **Minimum Go Version**: 1.21.0 (set in go.mod)
- **Analysis Target**: Any Go 1.21.0+ codebase
- **Recommendation**: Rebuild with your current Go version for optimal compatibility

## Installation and Usage

> **Note**: This linter was proposed for inclusion in golangci-lint but was ultimately rejected. See [golangci-lint PR #6083](https://github.com/golangci/golangci-lint/pull/6083) for the discussion. The maintainers cited existing linters as duplicates, but this linter serves a different purpose:
>
> - **nonamedreturns** (firefart/nonamedreturns): Flags named returns as bad practice - the opposite philosophy
> - **gocritic (unnamedResult)**: Suggests adding names but doesn't enforce or validate them
> - **revive (bare-return)**: Warns against bare returns but doesn't ensure proper named return usage
>
> This linter uniquely enforces consistent use of named returns and catches shadowed variables and signature mismatches. As a result, this standalone version provides an easy way to use the linter outside of the golangci-lint ecosystem.

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


### Option 5: golangci-lint integration per golangci-lint.run docs (Doesn't work at the time of this writing.)
Add to your `.golangci.yml`:
```yaml
linters:
  custom:
    namedreturns:
      type: module
      path: github.com/nikogura/namedreturns
      description: enforces the use of named returns in Go functions
      original-url: github.com/nikogura/namedreturns
  enable:
    - namedreturns
```

Then run: `golangci-lint run`
### Option 6: Use 'custom' directive in .golangci-lint.yml (Doesn't work at the time of this writing.)
The following syntax is supported by `golangci-lint`:

```yaml
  custom:
    namedreturns:
      path: github.com/nikogura/namedreturns
      type: module
      description: enforces the use of named returns in Go functions
      original-url: github.com/nikogura/namedreturns
```

However, the schema parsers built in to many tools and IDE's are not entirely up to date as of the time of this writing. While `golangci-lint run` will run the linter, many syntax/schema testers will choke on the 'custom' section.

For example:
```yaml
    - name: Lint
      uses: golangci/golangci-lint-action@v8
      with:
        version: latest
        verify: false   # Need this to prevent the action from choking on the 'custom' section.
```

## Named Returns in Deferred Statements

Named errors used in defers are not reported. If you also want to report them set `report-error-in-defer` to true.

## Further Reading

Tutorial on how to write your own linter:
https://disaev.me/p/writing-useful-go-analysis-linter/
