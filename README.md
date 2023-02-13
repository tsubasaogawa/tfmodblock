# tfmodblock (alpha)

## Overview

tfmodblock generates Terraform module block from variable blocks.

## Install

### Use install script

```bash
curl -H 'Accept: application/vnd.github.VERSION.raw' 'https://api.github.com/repos/tsubasaogawa/tfmodblock/contents/install.sh?ref=main' | bash
```

### Download an archive

Download an archive from [Releases](https://github.com/tsubasaogawa/tfmodblock/releases/latest) page.
Extract it and copy the binary to your PATH.

## Example

```hcl
$ cat example.tf
variable "foo" {
  type = string
}

variable "bar" {
  type = number
}

variable "baz" {
  type = map(number)
}

variable "lorem" {
  type = list(string)
}

variable "ipsum" {
  type = object({ a = string })
}
```

```hcl
$ tfmodblock .
module "tfmodblock" {
    source = "path/to/module"
    bar = 0
    baz = {}
    foo = ""
    ipsum = {}
    lorem = []
    
}
```

## Future works

- Create test code
- Auto indent
- Expand `object`

## Links

- <https://github.com/tsubasaogawa/tfmodblock-vscode-extension>
