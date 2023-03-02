# tfmodblock (alpha)

## Overview

tfmodblock generates Terraform module block HCL using variables from tf files.

## Features

- Auto generate module block
- Specify `source` by relative path
- Insert a value using `default` attribute
- Insert a description

## Install

### a. Use install script

```bash
curl -H 'Accept: application/vnd.github.VERSION.raw' 'https://api.github.com/repos/tsubasaogawa/tfmodblock/contents/install.sh?ref=main' | bash
```

### b. Download an archive manually

Download an archive from [Releases](https://github.com/tsubasaogawa/tfmodblock/releases/latest) page.
Extract it and copy the binary to your PATH.

## Example

```hcl
$ cat example.tf
variable "foo" {
  type = string
}

variable "bar" {
  type        = number
  description = "this is bar"
}

variable "baz" {
  type = map(number)
}

variable "lorem" {
  type    = list(string)
  default = ["lorem1", "lorem2"]
}

variable "ipsum" {
  type = object({ a = string })
}
```

```hcl
$ tfmodblock .
module "tfmodblock" {
    source = "."

    // this is bar
    bar = 0
    baz = {}
    foo = ""
    ipsum = {}
    lorem = [lorem1 lorem2]
}
```

## Help

Run tfmodblock with `--help` option.

## Future works

- Create test code
- Auto indent
- Expand `object`

## Links

- <https://github.com/tsubasaogawa/tfmodblock-vscode-extension>
