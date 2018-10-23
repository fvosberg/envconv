# Envconv

Envconv loads environment variables into a standard library FlagSet as defaults, if set.

## Why another env package

I wanted a package, which meets the following criteria:

- Enhance the existing standard library flag.FlagSet,
- to move a concise Go app with a flag set to a 12 factor app
- and which has the following fallback hierarchy: Flag > ENV > default
