# PubCopy
This package provides a reflect-based copying of values where only public fields of structures to be copied

## Usage example

```go
if err := pubcopy.Copy(src, dst, pubcopy.PublicOnly); err != nil {
    …
}
```
