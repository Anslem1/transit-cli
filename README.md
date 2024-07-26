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
transit create my-transit-name
```
**Adding command to a transit:** adds a new command to a transit.
```bash
transit add my-transit-name
```
**Edit a command in transit:** edits an existing command in a transit.
```bash
transit edit my-transit-name
```
**Search a command in transit:** This looks through all your transits and search for a the specified command.
```bash
transit search my-transit-name
```
**Eemoving command in a transit:** removes commands within an existing transi.
```bash
transit remove my-transit-name
```
**Reordering commands in a transit:** Changes the order at which commands are placed within a transit.
```bash
transit reorder my-transit-name
```
**List a transit:** Lists commands within a transit.
```bash
transit list my-transit-name
```

**Delete a transit:**
```bash
transit delete my-transit-name
```

**Execute a transit:** Executes commands within a transit.
```bash
transit execute my-transit-name
```
**You can also use 'transit [command]' without specifying the exact transit, and it will bring up a prompt to select the transit interactively.**

## Contributing 

Pull requests are welcome. For major changes, please open an issue first
to discuss what you would like to change.

Please make sure to update tests as appropriate.