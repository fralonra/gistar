# gistar

A golang tool to fetch a list of your starred repositories on Github, and output it to a markdown document. Inspired by [starred](https://github.com/maguowei/starred) and [repogen](https://github.com/skyjia/repogen).

## Usage

```bash
go get github.com/fralonra/gistar
```

A basic example:
```bash
gistar username > output.md
```
This will generate an `output.md` with all your starred repositories inside.

### Flags
| Flag | Description | Type | Default |
| --- | --- | --- | --- |
| -d | Hide description. | bool | false |
| -f | Hide forks. | bool | false |
| -l | Hide language. | bool | false |
| -s | Hide stars. | bool | false |
| -w | Hide watches. | bool | false |
| -sort | How to sort the repository list. Available values: 'created', 'updated', 'pushed', 'full_name' and 'lang'. The first four values are options for [go-github](https://godoc.org/github.com/google/go-github/github#RepositoryListOptions). By default, the value is `lang`, and it will sort repositories by their top language. | string | 'lang' |
| -stl | Badget styles. Available values: 'flat', 'flat-square', 'for-the-badget', 'plastic' and 'social'. See [here](https://shields.io/) for more infomation. | string | 'flat' |

```bash
gistar -f -l -s -w username > output.md
```
The above command will generate a file with only description for the repository.

## Example

```bash
gistar fralonra > README.md
```

You can check the output `README.md` [here](https://github.com/fralonra/gistar/tree/master/examples/README.md).
