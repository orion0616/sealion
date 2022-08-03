# Sealion
Sealion is a CLI tool for Todoist.

## How to use
1. Install
```
go install github.com/orion0616/sealion@latest
```

2. Setting

Set your API token to your environment variable named $TODOIST_TOKEN

3. Example
For example, you can add labels to items in a project
```
$ sealion add labels <label1> <label2> ... -p <projectName>
```
