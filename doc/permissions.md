Source: https://peternixon.net/news/2012/01/28/configure-tcpdump-work-non-root-user-opensuse-using-file-system-capabilities/


groupadd pcap
usermod -a -G pcap ${USER}

chgrp pcap /usr/sbin/tcpdump
chmod 750 /usr/sbin/tcpdump

setcap cap_net_raw,cap_net_admin=eip /usr/sbin/tcpdump