# **LINUX PRIVESC CHECKLIST**

## Try in Order

> ### Sudo Permissions
>
>     sudo -l

> ### Find Suid Binaries
>
>     find / -perm -4000 2>/dev/null

> ### Capabilities
>
>     getcap -r / 2>/dev/null

> ### List Listening Network Services
>
>     ss -tulnp

> ### Find all Files with .service Extension
>
>     find / -type f -name "*.service" 2>/dev/null

> ### Check Writable Files
>
>     find / -type f -writable 2>/dev/null

> ### Check Writable Directories
>
>     find / -type d -writable 2>/dev/null

> ### Run LinPeas
>
>     /usr/share/peass/linpeas/linpeas.sh
