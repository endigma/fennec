# watcher
watcher is a tool to execute <something> from webhook POST to a particular path

it expects a field as part of the JSON POST called "secret" like:

```json
{
    "secret": "secret-string"
}
```

which must match with the secret in the config in order for it to be executed.

it requires a file titled config.json in the working directory
- todo: make this a parameter

example config in repo, you have to remove the comments.

most errors are self-explanatory, ask me on #go-nuts or create a github issue if something is really broken.