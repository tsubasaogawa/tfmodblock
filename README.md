# tfmodblock (alpha)

## Overview

tfmodblock generates Terraform module block from variable blocks.

## Install

Download a binary from Releases page.
Copy the binary to your PATH.

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
- VSCode Extension
