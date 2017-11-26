# mtoi

[![Build Status](https://travis-ci.org/WindomZ/mtoi.svg?branch=master)](https://travis-ci.org/WindomZ/mtoi)
[![Coverage Status](https://coveralls.io/repos/github/WindomZ/mtoi/badge.svg?branch=master)](https://coveralls.io/github/WindomZ/mtoi?branch=master)

> Key-Value structures: KV, Slice, Cache...

## Feature

- **Key-Value** base storage structure.
- **Non-blocking** to write.
- **sync.Map** if >= go1.9.

### KV([kv.go](./kv.go))
- Value is **single** instance.
```
Put(k string, v interface{})
Get(k string) (v interface{}, ok bool)
```

### Slice([slice.go](./slice.go))
- Value is go **slice** structure.
```
Put(k string, v interface{})
Get(k string) (v []interface{}, ok bool)
```

### Cache([cache.go](./cache.go))
- Value is **single** instance.
- Support **expire** time.
```
Put(k string, v interface{}, expire time.Duration)
Get(k string) (v interface{}, ok bool)
```

## Contributing

Welcome to pull requests, report bugs, suggest ideas and discuss 
**mtoi** on [issues page](https://github.com/WindomZ/mtoi/issues).

If you like it then you can put a :star: on it.

## License

[MIT](https://github.com/WindomZ/mtoi/blob/master/LICENSE)
