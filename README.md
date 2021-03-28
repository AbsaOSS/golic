# golic
declarative license tool, injecting licenses into source code
```
golic inject -c="2021 MyCompany ltd." --dry
```
![Screenshot 2021-03-08 at 11 42 52](https://user-images.githubusercontent.com/7195836/110310942-6d2f3680-8003-11eb-9540-b2e21b4f2b87.png)

## Quickstart 
Install and run **GOLIC**
```shell
# GO 1.16 
go install github.com/AbsaOSS/golic@v0.5.0
golic version
```
Golic has two configurations `.licignore` and `.golic.yaml`. The first is based on the .gitignore and determines which 
files will be selected for license injection. The second contains a configuration with license text and formatting rules.

### .licignore
.licignore determines which files will be selected for license injection. The syntax of the file is the same as for .gitignore
For simplicity, we have created inverse rules - we denied everything and allowed where to place license.

Create `.licignore` in project root
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
### .golic.yaml
golic.yaml contains a configuration with license text and formatting rules. Golic uses embeded [master configuration](https://raw.githubusercontent.com/AbsaOSS/golic/main/.golic.yaml) 
by default. The master configuration is compiled and goes with binary, so it can change from version to version.
If you need to change configuration, you can override it. For example, you want to change the text license, 
or set comments for a specific file type. All you have to do is create a `.golic.yaml` file in the project root. Golic will 
read it and overrides master configuration rules

Example below overrides master configuration by adding `apacheX` licenses and sets new rule for `*.go.txt` nad `.mzm`.
For more details see [master configuration](https://raw.githubusercontent.com/AbsaOSS/golic/main/.golic.yaml) example.
```yaml
# .golic.yaml 
golic:
  licenses:
    apacheX: |
      Copyright MyCompany
      Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Pellentesque pretium lectus id turpis. 
      Suspendisse sagittis ultrices augue. Integer pellentesque quam vel velit. In sem justo, commodo ut 
      suscipit at, pharetra vitae, orci
      
      for more details see https://github.com/mycompany/myproject/LICENSE
  rules:
    "*.go.txt":
      prefix: "/*"
      suffix: "*/"
    .mzm:
      prefix: ""    # no indent, no prefix or suffix, just place license text into top of the file 
```

### Running from commandline
If you already created `.licignore` and `.golic.yaml`, run command : 
```shell
golic inject -t apacheX
```
Consider to use `--dry` flag to preview which files will be affected, before golic modify files.
For more command line options (like placeholders, default values etc.), see [Usage](#usage) section.

## CI support
Usually you want to find out that something went wrong during CI / CD. For example, a file is missing a license. 
In terms of golic, we want the build pipe to end with an error if we find at least one file with a missing license.
The `-x` argument handles that.
```shell
  go install github.com/AbsaOSS/golic@v0.5.0
  golic inject --dry -x -t apache2
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

