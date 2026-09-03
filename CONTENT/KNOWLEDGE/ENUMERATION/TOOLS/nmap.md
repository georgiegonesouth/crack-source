# COPY PASTE COMMANDS

## all ports:
    nmap -p- --min-rate 10000 -oA nmap/allports <TARGET_IP>

## detailed:
    nmap -p 22,80,443 -sC -sV -oA nmap/detailed <TARGET_IP>

## udp:
    sudo nmap -sU --top-ports 20 -oA nmap/udp <TARGET_IP>


# FLAG COLLECTION:

## scan types

    -sS: SYN SCAN, considered more stealthy -> no full tcp connection
    -sT: TCP SCAN, full three way handshake
    -sU: UDP scan
    -sA: TCP ACK scan, harder for firewalls to filter because it only sends an ACK flag

## ports

    -p <PORTS: e.g individual 22,80,443 or range 22-443>
    --top-ports=int: most used ports from nmap database
    -p- all ports
    -F: fast scan of top 100 ports

## packets

    --packet-trace: shows all packets sent and recieved
    -n: disable DNS
    --disable-arp-ping: disables arp ping
    -Pn: disable icmp echo request

## options

    --stats-every=5s: show stats every {int}(s)econds or {int}(m)inutes
    -A: performs service detection, OS detection, traceroute and uses defaults scripts to scan the target
    -O: OS detection
    --traceroute: traceroute

## output
    -oN <scan_name>: save with .nmap extension 
    -oG <scan_name>: save with .gnamp extension (greppable)
    -oX <scan_name>: save with .xml extension
    -oA <scan_name>: save all output formats 


## scripts

    -sC: use standard scripts
    --script: <category> or <script-name>

## performance

    --initial-rtt-timeout 50ms: sets time value as initial rtt timeout
    --max-rtt-timeout 100ms: sets time value as maximum rtt timeout
    --max-retries 0: sets the number of retries that will be performed during the scan
    --min-rate 300: sets the minimum number of packets to be sent per second
    -T 0-5: sets the timing template from paranoid (0) to insane (9) 

## evasion

    -D RND:5: generates 5 random decoy IP addresses to disguise the real source
    -S: scans target by using different source IP address
    -e tun0: sends all requests through the specified interface
    --dns-server: specifies dns server
    --source-port: specifies source port

# MISC

## NSE

### NSE scripts and the corresponding categories
    
    https://nmap.org/nsedoc/index.html

### script categories

    auth:	  Determination of authentication credentials.
    
    broadcast:	Scripts, which are used for host discovery by broadcasting and the discovered hosts, can be automatically added to the remaining scans.
    
    brute:	Executes scripts that try to log in to the respective service by brute-forcing with credentials.
    
    default:	Default scripts executed by using the -sC option.
    
    discovery:	Evaluation of accessible services.
    
    dos:	These scripts are used to check services for denial of service vulnerabilities and are used less as it harms the services.
    
    exploit:	This category of scripts tries to exploit known vulnerabilities for the scanned port.
    
    external: 	Scripts that use external services for further processing.
    
    fuzzer:	This uses scripts to identify vulnerabilities and unexpected packet handling by sending different fields, which can take much time.
    
    intrusive:	Intrusive scripts that could negatively affect the target system.
    
    malware: 	Checks if some malware infects the target system.
    
    safe:	Defensive scripts that do not perform intrusive and destructive access.
    
    version:	Extension for service detection.
    
    vuln: 	Identification of specific vulnerabilities.

## Visualize 
    Convert XML to HTML for visualization: xsltproc target.xml -o target.html
