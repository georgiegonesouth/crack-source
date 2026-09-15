# HASHCAT 

## Dictionary Attacks

> ### Basic
>
>     hashcat -a 0 -m <MODE> <HASH> <WORDLIST>

> ### With Rule
>
>     hashcat -a 0 -m <MODE> <HASH> <WORDLIST> -r <RULE>

### good wordlists

> <!-- token-section:WORDLIST -->
>
>```
> - /usr/share/wordlists/rockyou.txt
>```

### good rules

> <!-- token-section:RULE -->
>
>```
> - /usr/share/hashcat/rules/best64.rule
>```

## Building Rules

> ### Create a List
>
>     hashcat --force <WORDS> -r <CUSTOM_RULE> --stdout | sort -u > <OUTPUT_LIST>

### Operators

| Name | Function | Description | Example Rule | Input Word | Output Word | Note |
|---|---|---|---|---|---|---|
| Nothing | : | Do nothing (passthrough) | : | p@ssW0rd | p@ssW0rd | |
| Lowercase | l | Lowercase all letters | l | p@ssW0rd | p@ssw0rd | |
| Uppercase | u | Uppercase all letters | u | p@ssW0rd | P@SSW0RD | |
| Capitalize | c | Capitalize the first character, lowercase the rest | c | p@ssW0rd | P@ssw0rd | |
| Invert Capitalize | C | Lowercase the first character, uppercase the rest | C | p@ssW0rd | p@SSW0RD | |
| Toggle Case | t | Toggle the case of all characters in word | t | p@ssW0rd | P@SSw0RD | |
| Toggle @ | TN | Toggle the case of characters at position N | T3 | p@ssW0rd | p@sSW0rd | * |
| Reverse | r | Reverse the entire word | r | p@ssW0rd | dr0Wss@p | |
| Duplicate | d | Duplicate entire word | d | p@ssW0rd | p@ssW0rdp@ssW0rd | |
| Duplicate N | pN | Append duplicated word N times | p2 | p@ssW0rd | p@ssW0rdp@ssW0rdp@ssW0rd | * |
| Reflect | f | Duplicate word reversed | f | p@ssW0rd | p@ssW0rddr0Wss@p | |
| Rotate Left | { | Rotate the word left | { | p@ssW0rd | @ssW0rdp | |
| Rotate Right | } | Rotate the word right | } | p@ssW0rd | dp@ssW0r | |
| Append Character | $X | Append character X to end | $1$2 | p@ssW0rd | p@ssW0rd12 | |
| Prepend Character | ^X | Prepend character X to front | ^2^1 | p@ssW0rd | 12p@ssW0rd | |
| Truncate left | [ | Delete first character | [ | p@ssW0rd | @ssW0rd | |
| Truncate right | ] | Delete last character | ] | p@ssW0rd | p@ssW0r | |
| Delete @ N | DN | Delete character at position N | D3 | p@ssW0rd | p@sW0rd | * |
| Extract range | xNM | Extract M characters, starting at position N | x04 | p@ssW0rd | p@ss | * # |
| Omit range | ONM | Delete M characters, starting at position N | O12 | p@ssW0rd | psW0rd | * |
| Insert @ N | iNX | Insert character X at position N | i4! | p@ssW0rd | p@ss!W0rd | * |
| Overwrite @ N | oNX | Overwrite character at position N with X | o3$ | p@ssW0rd | p@s$W0rd | * |
| Truncate @ N | 'N | Truncate word at position N | '6 | p@ssW0rd | p@ssW0 | * |
| Replace | sXY | Replace all instances of X with Y | ss$ | p@ssW0rd | p@$$W0rd | |
| Purge | @X | Purge all instances of X | @s | p@ssW0rd | p@W0rd | |
| Duplicate first N | zN | Duplicate first character N times | z2 | p@ssW0rd | ppp@ssW0rd | * |
| Duplicate last N | ZN | Duplicate last character N times | Z2 | p@ssW0rd | p@ssW0rddd | * |
| Duplicate all | q | Duplicate every character | q | p@ssW0rd | pp@@ssssWW00rrdd | |
| Extract memory | XNMI | Insert substring of length M starting from position N of word saved to memory at position I | lMX428 | p@ssW0rd | p@ssw0rdw0 | + |
| Append memory | 4 | Append the word saved to memory to current word | uMl4 | p@ssW0rd | p@ssw0rdP@SSW0RD | + |
| Prepend memory | 6 | Prepend the word saved to memory to current word | rMr6 | p@ssW0rd | dr0Wss@pp@ssW0rd | + |
| Memorize | M | Memorize current word | lMuX084 | p@ssW0rd | P@SSp@ssw0rdW0RD | + |

## Mask Attacks

> ### Basic
>
>     hashcat -a 3 -m <MODE> <HASH> <MASK>

## Mask Operators

| Symbol | Charset |
|---|---|
| ?l | abcdefghijklmnopqrstuvwxyz |
| ?u | ABCDEFGHIJKLMNOPQRSTUVWXYZ |
| ?d | 0123456789 |
| ?h | 0123456789abcdef |
| ?H | 0123456789ABCDEF |
| ?s | «space»!"#$%&'()*+,-./:;<=>?@[]^_`{ |
| ?a | ?l?u?d?s |
| ?b | 0x00 - 0xff |

## Filter Modes
    hashcat --help | grep -A 475 'Hash modes' | grep -i <KEYWORD>