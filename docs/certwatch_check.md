## certwatch check

The check command will validate the host from the command. (Defaults to port 443.)

### Synopsis

Check the host specified in the command Defaults to port 443.

```
certwatch check [host] [flags]
```

### Options

```
  -h, --help          help for check
  -p, --port string   The port to check the host on. (default "443")
```

### Options inherited from parent commands

```
  -c, --config string   path to your file containing environment variables (default "/etc/env/.certwatch.env")
```

### SEE ALSO

* [certwatch](certwatch.md)	 - certwatch is a simple recurring task to validate ssl certificates via a systemd service

###### Auto generated by spf13/cobra on 13-Jan-2024