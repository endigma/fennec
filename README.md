# fennec
fennec is a tool to execute <something> from webhook POST to a particular URL

it expects a field as part of the JSON POST called "secret" like:

```json
{
    "secret": "secret-string",
    "data": ["b", "f", "21"]
}
```

things you should know:
- there is also a field called "data" that can accept any data type, this will be appended as JSON to your script, if enabled.
- which must match with the secret in the config in order for it to be executed.
- it requires a parameter for a path to a config file
- example config in repo, you have to remove the comments.

## docker support

Fennec can be built into a container.
