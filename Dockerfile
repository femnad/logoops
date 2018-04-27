from fedora

run dnf -y update
run dnf -y install rsyslog

copy debug.conf /etc/rsyslog.d/debug.conf
copy tcp-input.conf /etc/rsyslog.d/tcp-input.conf
copy rsyslog.conf /etc/rsyslog.conf

entrypoint rsyslogd -nd

expose 514
