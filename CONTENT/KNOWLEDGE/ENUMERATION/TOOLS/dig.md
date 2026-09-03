# **DIG**

## Common Commands

| Command                       | Description                                                                                                                                                     |
| ----------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `dig <DOMAIN>`                | Performs a default A record lookup for the domain.                                                                                                              |
| `dig <DOMAIN> A`              | Retrieves the IPv4 address (A record) associated with the domain.                                                                                               |
| `dig <DOMAIN> AAAA`           | Retrieves the IPv6 address (AAAA record) associated with the domain.                                                                                            |
| `dig <DOMAIN> MX`             | Finds the mail servers (MX records) responsible for the domain.                                                                                                 |
| `dig <DOMAIN> NS`             | Identifies the authoritative name servers for the domain.                                                                                                       |
| `dig <DOMAIN> TXT`            | Retrieves any TXT records associated with the domain.                                                                                                           |
| `dig <DOMAIN> CNAME`          | Retrieves the canonical name (CNAME) record for the domain.                                                                                                     |
| `dig <DOMAIN> SOA`            | Retrieves the Start of Authority (SOA) record for the domain.                                                                                                   |
| `dig @<NAMESERVER> <DOMAIN>`       | Queries a specific DNS server; in this example, `1.1.1.1`.                                                                                                      |
| `dig +trace <DOMAIN>`         | Shows the full DNS resolution path from the root servers to the authoritative server.                                                                           |
| `dig -x <TARGET_IP>`          | Performs a reverse DNS lookup (PTR record) on the IP address to find the associated hostname. You may need to specify a name server.                            |
| `dig +short <DOMAIN>`         | Provides a short, concise answer to the query.                                                                                                                  |
| `dig +noall +answer <DOMAIN>` | Displays only the answer section of the query output.                                                                                                           |
| `dig <DOMAIN> ANY`            | Retrieves all available DNS records for the domain. **Note:** Many DNS servers ignore `ANY` queries to reduce load and prevent abuse, as described in RFC 8482. |
