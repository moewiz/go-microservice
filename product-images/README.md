# Product Images

## Uploading

*Note*: need to use `--data-bianry` to ensure file is not converted to text

```zsh
curl -vv localhost:9090/1/go.mod -X PUT --data-binary @test.png
```
