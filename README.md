# tlscheck
A TLS certificate checking and alerting tool

## Usage

```shell
tlscheck config.yaml
```
## Configuration

```yaml
domains:
  - example.com
validation:
  #debug: true
  days_advance: 7
notify:
  wxwork_bot: "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
```

## Credits

GUO YANKE, MIT License