variable "foo" {
  type    = string
  default = "a"
}

variable "bar" {
  type        = number
  description = "this is bar"
}

variable "baz" {
  type        = map(number)
  description = <<EOD
    this is heredoc
    description
  EOD
}

variable "lorem" {
  type    = list(string)
  default = ["lorem1", "lorem2"]
}

variable "ipsum" {
  type    = object({ a = string })
  default = { a = "ipsum1" }
}
