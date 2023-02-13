variable "foo" {
  type = string
}

variable "bar" {
  type    = number
  default = 123
}

variable "baz" {
  type = map(number)
}

variable "lorem" {
  type        = list(string)
  description = "lorem description"
  default     = ["lorem1", "lorem2"]
}

variable "ipsum" {
  type    = object({ a = string, b = { b1 = number, b2 = list(string) } })
  default = { a = "ipsum", b = { b1 = 1, b2 = ["ipsum1", "ipsum2"] } }
}

variable "dolor" {
  type    = bool
  default = false
}
