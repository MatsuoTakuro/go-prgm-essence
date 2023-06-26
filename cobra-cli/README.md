# cobra-cli

```bash
go install github.com/spf13/cobra-cli@latest
# create .cobra.yaml in your project
ln -s /path/to/your/cobra-cli/.cobra.yaml ~/.cobra.yaml
cobra-cli init
go run main.go
cobra-cli add update
# add flags and the corresponding processes to update command
go run main.go update -h
go run main.go update --foo hello -t
```
