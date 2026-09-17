# Getting started with transit-cli

Transit is a CLI tool designed to help developers save time by running a set of usual terminal commands all at once.

## Installation
To install Transit using Homebrew, follow these steps:

```bash
brew tap Anslem1/homebrew-tap
```
*After tapping it, run:*
```bash
 brew install transit-cli
```

## Common Transit Commands

**Creating a Transit:** Initializes a new, empty transit.
```bash
# Global transit (stored in your user config)
transit create my-transit-name

# Project-local transit (saved in .transit.yaml in the current repo)
transit create my-transit-name --local
```

**Adding commands to a transit:**
```bash
transit add my-transit-name
# Or directly from the command line:
transit add my-transit-name "npm install" "npm run build"
```

**Dynamic Arguments / Placeholders:**
Commands can include placeholders like `$1`, `$2`, `"$@"`, or `{{1}}`:
```bash
transit add deploy "docker build -t myapp:$1 ." "docker push myapp:$1"
# Execute with arguments:
transit execute deploy v1.2.0
```

**Edit a command in transit:**
```bash
transit edit my-transit-name
```

**Search commands across transits:**
```bash
transit search "git checkout"
```

**Removing a command from a transit:**
```bash
transit remove my-transit-name
```

**Reordering commands in a transit:**
```bash
transit reorder my-transit-name
```

**List transits and commands:**
```bash
# View summary of all available transits (local and global)
transit list --all

# List commands inside a specific transit
transit list my-transit-name
```

**Delete a transit:**
```bash
transit delete my-transit-name
```

**Executing a transit:**
```bash
# Run sequentially (default, stops immediately on error)
transit execute my-transit-name

# Run all commands concurrently in parallel with color-tagged logs
transit execute my-transit-name --parallel # or -p

# Skip interactive prompts
transit execute my-transit-name --skip # or -s

# Continue running subsequent commands even if one fails
transit execute my-transit-name --continue-on-error # or -c

# Preview commands without executing
transit execute my-transit-name --dry-run
```
**You can also use 'transit [command]' without specifying a transit, and it will bring up an interactive prompt.**

## Shell Autocompletion

Transit supports shell autocompletion for `bash`, `zsh`, `fish`, and `powershell`. To enable it for `zsh`:

```bash
# Enable completion in current session:
source <(transit completion zsh)

# Or save permanently to your zsh completion dir:
transit completion zsh > "${fpath[1]}/_transit"
```
Once enabled, typing `transit execute <TAB>` or `transit add <TAB>` will automatically suggest all your available transits.

## Contributing 

Pull requests are welcome. For major changes, please open an issue first
to discuss what you would like to change.

Please make sure to update tests as appropriate.