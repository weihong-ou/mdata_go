# Commands run to package golang as a service

1. Move environment files to correct location /etc/default<br>
    `sudo cp packaging/systemd/etc/default/sawtooth-mdata-tp-go /etc/default/sawtooth-mdata-tp-go`
    `sudo chmod 644 /etc/default/sawtooth-mdata-tp-go`

2. Move service files to correct location /lib/systemd/system<br>
    `sudo cp packaging/systemd/lib/systemd/system/sawtooth-mdata-tp-go.service /lib/systemd/system/`
    `sudo chmod 644 /lib/systemd/system/sawtooth-mdata-go.service`

3. Move binaries to correct location /usr/bin<br>
    `sudo cp bin/mdata-tp-go /usr/bin/`
    `sudo chmod 755 /usr/bin/mdata-tp-go`

4. Enable services <br>
    `sudo systemctl enable sawtooth-mdata-tp-go.service`

# Enable logs to journalctl
1. Uncomment the following in `/etc/rsyslog.conf`
    ```
    module(load="imtcp")
    input(type="imtcp" port="514")
    ```

2. Create /etc/rsyslog.d/30-mdata-tp-go.conf
    `sudo touch /etc/rsyslog.d/30-mdata-tp-go.conf` <br>   
     ```echo -e "if $programname == 'sawtooth-mdata-tp-go' or $syslogtag == 'mdata_tp' then /var/log/mdata_tp/mdata_tp.log & stop" > /etc/rsyslog.d/30-mdata-tp-go.conf```<br>    
    `sudo systemctl restart rsyslog`