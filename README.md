# golic
license management tool, injecting licenses into source code
```
golic inject -c="2021 MyCompany ltd." --dry
```
![Screenshot 2021-03-08 at 11 42 52](https://user-images.githubusercontent.com/7195836/110310942-6d2f3680-8003-11eb-9540-b2e21b4f2b87.png)


## Running from commandline

create `.licignore` in project root
```shell
# Ignore everything
*

# But not these files...
!Dockerfile*
!Makefile
!*.go

# ...even if they are in subdirectories
!*/
````
Install and run **GOLIC**
```shell
# GO 1.16 
go install github.com/AbsaOSS/golic@v0.4.8
golic inject -c="2021 MyCompany ltd."
```

## Usage
```
Available Commands:
  help        Help about any command
  inject      Injects license
  version     Print the version number of Golic

Flags:
  -h, --help      help for this command
  -v, --verbose   verbose output

Usage inject:
   inject [flags]

Flags:
  -p, --config-path string   path to the local configuration overriding config-url (default ".golic.yaml")
  -u, --config-url string    default config URL (default "https://raw.githubusercontent.com/AbsaOSS/golic/main/.golic.yaml")
  -c, --copyright string     company initials entered into license (default "2021 MyCompany")
  -d, --dry                  dry run
  -h, --help                 help for inject
  -l, --licignore string     .licignore path (default ".licignore")
  -x, --modified-exit        If enabled, exits with status 1 when any file is modified. The settings is used by CI
  -t, --template string      license key (default "apache2")


Global Flags:
  -v, --verbose   verbose output
```

## Configuration
Golic uses embeded [master configuration](https://raw.githubusercontent.com/AbsaOSS/golic/main/.golic.yaml) by default.
The master configuration is compiled within binary, and if you need to change it, create PR.
However, it is much better to create a local configuration that overrides the master configuration settings. All 
you have to do is create a `.golic.yaml` file in the project root, or use the` -p` flag.
Example below overrides master configuration by custom licenses
```yaml
# .golic.yaml 
golic:
  licenses:
    apache2: |
      Copyright {{copyright}}

      This is my custom license text
    apacheX: |
      Copyright {{copyright}}
      for more details see https://github.com/mycompany/myproject/LICENSE
```

