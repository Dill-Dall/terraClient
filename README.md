# terraClient

This client is a wrapper for the terragrunt wrapper. It consist of the following functions:

## importfromto

Method for importing state in one state directory to another.

### usage

```
terraClient stateImport --help
Does import of item from state list in src state into destination state.
        Fetches id for import from state output.
        Delete original state at end.
        Error: States can be dependant on eachother. Import does not remove in original state on fail.

Usage:
  terraClient stateImport [flags]

Flags:
  -d, --dest string     folder with terragruntstate to import to
  -h, --help            help for stateImport
  -s, --src string      folder with terragruntstate to import from
  -u, --update-config   also update config file (default: false)   [ not supported ]

terracli  importfromto  -s "/home/thomas/training/go/terragrunt/two" -d "/home/thomas/training/go/terragrunt/one"
```

### Why it exists

When upgrading in our setup and there are breaking changes. Every configuration needs to be upgraded at the same time. Making migrations a real hassle.
With this client the state can be split up and upgrades can be done one configuration at the time.
