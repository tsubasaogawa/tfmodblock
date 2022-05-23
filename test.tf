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
