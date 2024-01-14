#  certwatch
Routinely checks ssl certificate validity on specified hosts. Will notify the specified
individuals by email of expiration reminders up to 14 days prior, and immmediately alert if the certificate
is no longer valid. 

### Installation and Configuration.

1. Install certwatch:

```
$ go install github.com/masonictemple4/certwatch@latest
```

2. Configure your [environment variables](#Environment Variables):
See the section below.


3. (optional) Create a hosts.json file using the format below in your preferred location:

```json
[
    {"hostname": "example.com", "port": "443"},
    {"hostname": "test.example.com", "port": "443"},
]
```

4. (Optional) Create a log file for the service output. I prefer `/var/log/certwatch.log`. All logging,
defaults to os.Stdout if you decide not to log to a file.


### Usage:


**One off checks**:

    $ certwatch check example.com

The default port is 443 however if necessary you can use the `-p` flag to specify an
alternative one.

This will output a preview of the Cert and a list of any errors that may be associated to that.
certificate.

**Manually running the scheduled task**:

    $ certwatch run www.google.com 443


#### Setup a systemd service:

1. Create a new service file: `/etc/systemd/system/certwatch.service`

    ```ini
    [Unit]
    Description=certwatch service

    [Service]
    User=certwatchcertwtch
    Group=certwtch
    Type=simple
    Restart=always
    ExecStart=/usr/local/go/bin/certwatch run -H /etc/certwatch/hosts.json -l /var/log/certwatch.log

    [Install]
    WantedBy=multi-user.target
    ```

Above is a very basic example of what one might look like!


2. Run this whenever you create a new service or make changes to an exisisting one: 

    ```
    $ sudo systemctl daemon-reload
    ```

3. Next, let's start the service:  
    
    ```
    $ sudo systemctl start certwatch.service
    ```

4. (Optional) If you would like this service to start automatically after bootup.

    ```
    $ sudo systemctl enable certwatch.service
    ```

You can also **disable** startup after boot with by calling `disable` instead


**Stopping your service**:

    $ sudo systemctl stop certwatch.service

**Restarting your service** (Run this if you have updated/re-installed a service)  

    $ sudo systemctl restart certwatch.service

**Checking the status of your service**: 

    $ sudo systemctl status certwatch.service

This will also show some of the stdout associated
with the particular service. You can use the `-n` flag to increase how many lines it shows. 
You can also attach yourself to the output.



### Environment Variables (REQUIRED):
This application depends on environment variables, 
you can specify your own `.env` file using the `-c` flag.

Or you can create a directory `/etc/env/` and and place a file
named `.certwatch.env` in it this is the default location.  

You may also set environment variables inside of your 

`/etc/systemd/system/certwatch.service` definition.

```
# You can choose any provider you'd like
# I went with Gmail because I already had it setup.
SMTP_HOST="smtp.gmail.com"
SMTP_PORT="587"
EMAIL="your-sender@gmail.com"

# NOTE: If you have 2FA enabled domain wide or on the individual
# account you wish to send mail from you will need to grab an app token
# in order to properly authenticate.
#
# Otherwise this can be your normal password.
EMAIL_PW="secret-here"

# You may specify one or many emails to receive the 
# notifications for expiring and invalid certs.
# NOTE: This value must be compliant with the RFC 822 To: header.
REPORT_RECIPIENTS="target1@yahoo.com, target2@hotmail.com, target3@gmail.com"
```


#### Visit [CLI Docs](./docs/certwatch.md) for additional information.


**Quick note** You can use specify a [host] command line argument while also specifying
a hosts.json file with the `-H` flag, this will just append the one to the list after
it's read from the file.
