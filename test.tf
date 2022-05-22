variable "foo" {
  type = string
}

variable "bar" {
  type = number
}

variable "baz" {
  type = any
}

variable "lorem" {
  type = list(string)
}

variable "ipsum" {
  type = object({ a = string })
}

variable "taro" {
}

variable "no" {
  type = abc
}
