# Directory-Based Config Plan

## Overview

Change config from single file (`config.sh`) to directory-based (all `*.sh` files in directory, recursively).

## Current Behavior

- Config path: `$XDG_CONFIG_HOME/ifg/config.sh` or `~/.ifg/config.sh`
- Single file parsing
- Backward compatible with existing `config.sh`

## New Behavior

- Config directory: `$XDG_CONFIG_HOME/ifg/` or `~/.ifg/`
- Read all `*.sh` files recursively
- Sort alphabetically by relative path
- Parse each file, merge entries
- If directory doesn't exist, create it with default `config.sh`

## File Ordering

Alphabetical by relative path within config directory:

```
01-git.sh
02-docker.sh
personal/notes.sh
work/01-ssh.sh
work/02-deploy.sh
```

Users can prefix with numbers for custom ordering.

## Implementation

### 1. `internal/config/config.go`

**Rename functions:**
- `GetConfigPath()` → `GetConfigDir()` - returns directory path

**Update `LoadConfig()`:**
- Accept directory path instead of file path
- Use `filepath.WalkDir` to collect all `*.sh` files
- Sort files alphabetically
- Parse each file, append entries
- Return error if directory doesn't exist

**Update `CreateDefaultConfig()`:**
- Create directory if it doesn't exist
- Write default `config.sh` inside directory

**Add `parseFile()` helper:**
- Extract current parsing logic into separate function
- Takes file path, returns `[]Entry`

### 2. `main.go`

**Update config loading:**
- Call `config.GetConfigDir()`
- Check if directory exists
- If not, call `config.CreateDefaultConfig(configDir)`
- Call `config.LoadConfig(configDir)`

**Update `--help` output:**
- Show directory path instead of file path
- Update wording if needed

### 3. `README.md`

**Update Configuration section:**
- Document directory-based config
- Mention recursive subdirectory support
- Show example structure:
  ```
  ~/.ifg/
  ├── 01-git.sh
  ├── 02-docker.sh
  ├── personal/
  │   └── scripts.sh
  └── work/
      ├── 01-ssh.sh
      └── 02-deploy.sh
  ```

## Test Cases

### `internal/config/config_test.go`

1. **TestGetConfigDir**
   - Returns `$XDG_CONFIG_HOME/ifg/` when `XDG_CONFIG_HOME` is set
   - Returns `~/.ifg/` when `XDG_CONFIG_HOME` is not set

2. **TestLoadConfigEmptyDirectory**
   - Empty directory returns empty entries

3. **TestLoadConfigSingleFile**
   - Directory with one `*.sh` file returns correct entries

4. **TestLoadConfigMultipleFiles**
   - Directory with multiple `*.sh` files
   - Entries are merged in alphabetical order

5. **TestLoadConfigRecursive**
   - Directory with subdirectories containing `*.sh` files
   - All files are included, sorted alphabetically by path

6. **TestLoadConfigIgnoresNonShFiles**
   - Non-`.sh` files are ignored

7. **TestLoadConfigNonexistentDirectory**
   - Returns error for nonexistent directory

8. **TestCreateDefaultConfig**
   - Creates directory if it doesn't exist
   - Creates default `config.sh` inside directory

9. **TestParseFile**
   - Correctly parses entries from a single file

## Files to Modify

| File | Changes |
|------|---------|
| `internal/config/config.go` | Rename functions, update logic for directory |
| `main.go` | Update config loading, update help output |
| `README.md` | Update documentation |
| `internal/config/config_test.go` | Add/update tests |

## Backward Compatibility

Existing users with `~/.ifg/config.sh`:
- `~/.ifg/` is already a directory
- `config.sh` is still a valid `.sh` file
- Will be read as part of directory scan
- No breaking change

## Edge Cases

1. **Empty config directory** - Returns empty entries, no error
2. **Symlinks** - `filepath.WalkDir` follows symlinks by default, may need to handle cycles
3. **Permission errors** - Return error from `WalkDir` function
4. **Invalid UTF-8 filenames** - Go handles UTF-8 natively, should be fine

## Alternatives Considered

1. **Non-recursive only** - Simpler but less flexible. Decided against since recursive is easy.
2. **Explicit include order via config** - More complex, users can use filename prefixes instead.