## certwatch completion powershell

Generate the autocompletion script for powershell

### Synopsis

Generate the autocompletion script for powershell.

To load completions in your current shell session:

	certwatch completion powershell | Out-String | Invoke-Expression

To load completions for every new session, add the output of the above command
to your powershell profile.


```
certwatch completion powershell [flags]
```

### Options

```
  -h, --help              help for powershell
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
  -c, --config string   path to your file containing environment variables (default "/etc/env/.certwatch.env")
```

### SEE ALSO

* [certwatch completion](certwatch_completion.md)	 - Generate the autocompletion script for the specified shell

###### Auto generated by spf13/cobra on 13-Jan-2024
