# `profile`

Take profile data for emulation core using `profile` package.

## Usage

```sh
> make profile
> ./build/profile/profile -s=30 ROM_PATH
> go tool pprof -http=":8081" ./build/profile/profile ./build/profile/cpu.pprof
# or: go tool pprof -png ./build/profile/profile ./build/profile/cpu.pprof > ./build/profile/cpu.pprof.png
```
