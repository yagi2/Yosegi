# Security Policy

## Overview

Yosegi is a git worktree management tool that executes external commands. This document outlines the security measures implemented to prevent command injection and other security vulnerabilities.

## Security Measures

### 1. Input Validation

All user inputs are validated before being passed to external commands:

#### Branch Name Validation
- Must match the pattern `^[a-zA-Z0-9._/-]+$`
- Cannot start with a dash (`-`)
- Cannot start or end with a dot (`.`)
- Cannot contain consecutive dots (`..`)
- Cannot contain dangerous shell characters: `;`, `&`, `|`, `$`, `` ` ``, `(`, `)`, `<`, `>`, `"`, `'`, `\`, newlines

#### Path Validation
- Cannot be empty
- Cannot contain dangerous shell characters (same as above)
- Cannot contain deep directory traversal patterns (`../../../`)
- Cannot attempt to access restricted system directories:
  - `/etc/`
  - `/usr/bin/`
  - `/bin/`
  - `/sbin/`
  - `/root/`
  - `/home/root/`

### 2. Command Execution Security

#### Git Commands
All git commands are executed using Go's `exec.Command()` with:
- Hardcoded command names (`git`)
- Validated arguments
- Explicit working directory setting (`cmd.Dir`)
- No shell interpretation

#### Terminal Commands
Terminal control commands (`stty`) use:
- Hardcoded command names and arguments
- No user input in command construction

### 3. Prevented Attack Vectors

The following attack vectors are prevented by the validation system:

#### Command Injection
```bash
# These would be blocked:
yosegi new "/tmp/test;rm -rf /"
yosegi new "/tmp/test|cat /etc/passwd"
yosegi new "/tmp/test`whoami`"
```

#### Directory Traversal
```bash
# These would be blocked:
yosegi new "/tmp/../../../etc/passwd"
yosegi new "/etc/passwd"
```

#### Argument Injection
```bash
# These would be blocked:
yosegi new --path="/tmp" --branch="--exec=malicious"
yosegi new --path="/tmp" --branch="-dangerous"
```

## Testing

Security measures are thoroughly tested with comprehensive test suites:

- `TestValidateBranchName`: Tests branch name validation
- `TestValidatePath`: Tests path validation  
- `TestManagerAddWithSecurityValidation`: Tests worktree addition security
- `TestManagerRemoveWithSecurityValidation`: Tests worktree removal security
- `TestManagerDeleteBranchWithSecurityValidation`: Tests branch deletion security

Run security tests with:
```bash
go test -v ./internal/git/ -run "SecurityValidation|Validate"
```

## Reporting Security Issues

If you discover a security vulnerability, please:

1. **Do not** create a public GitHub issue
2. Email the maintainers privately
3. Include detailed information about the vulnerability
4. Allow time for the issue to be fixed before public disclosure

## Security Best Practices for Users

When using Yosegi:

1. **Validate inputs**: Always verify branch names and paths before use
2. **Use absolute paths**: When possible, use absolute paths to avoid confusion
3. **Review commands**: Use the `--dry-run` flag (if available) to preview actions
4. **Keep updated**: Always use the latest version with security patches

## Implementation Details

### Validation Functions

The core security is implemented in `/internal/git/worktree.go`:

- `validateBranchName(branch string) error`: Validates git branch names
- `validatePath(path string) error`: Validates file system paths

These functions are called by all public methods that accept user input:
- `Add(path, branch string, createBranch bool) error`
- `Remove(path string, force bool) error`
- `DeleteBranch(branch string, force bool) error`
- `HasUnpushedCommits(branch string) (bool, int, error)`

### Security Boundaries

Yosegi operates within these security boundaries:

1. **File system access**: Limited to the git repository and its worktrees
2. **Command execution**: Only executes `git` and `stty` commands
3. **Network access**: No direct network operations (git may access remotes)
4. **User permissions**: Runs with the same permissions as the invoking user

## Changelog

### Version 0.1.0
- Initial security implementation
- Added input validation for all user inputs
- Added comprehensive security test suite
- Prevented command injection vulnerabilities
- Added path traversal protection