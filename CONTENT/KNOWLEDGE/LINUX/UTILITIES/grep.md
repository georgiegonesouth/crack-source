## Usage

### BASIC COMMANDS

# Basic search
grep "text" file.txt

# Case-insensitive search
grep -i "text" file.txt

# Recursive search in current directory
grep -r "text" .

# Recursive search showing filenames and line numbers
grep -rn "text" .

# Search only certain file types
grep -r --include="*.py" "text" .

# Search multiple file types
grep -r --include="*.{js,ts}" "text" .

# Exclude file types
grep -r --exclude="*.log" "text" .

# Exclude directories
grep -r --exclude-dir=node_modules "text" .

# Show line numbers
grep -n "text" file.txt

# Show only matching filenames
grep -l "text" *

# Show filenames that do NOT contain the match
grep -L "text" *

# Count matches
grep -c "text" file.txt

# Show only the matched text
grep -o "pattern" file.txt

# Invert match (everything except matching lines)
grep -v "text" file.txt

# Match whole words only
grep -w "word" file.txt

# Match entire lines
grep -x "exact line" file.txt

# Extended regular expressions
grep -E "cat|dog|bird" file.txt

# Search multiple patterns
grep -e "foo" -e "bar" file.txt

# Read patterns from a file
grep -f patterns.txt file.txt

# Ignore binary files
grep -I "text" *

# Treat binary files as text
grep -a "text" file.bin

# Quiet mode (exit code only)
grep -q "text" file.txt

# Show 3 lines after each match
grep -A 3 "text" file.txt

# Show 3 lines before each match
grep -B 3 "text" file.txt

# Show 3 lines before and after each match
grep -C 3 "text" file.txt

# Highlight matches (usually enabled automatically)
grep --color=auto "text" file.txt

# Search for lines starting with a pattern
grep "^start" file.txt

# Search for lines ending with a pattern
grep "end$" file.txt

# Match empty lines
grep "^$" file.txt

# Match non-empty lines
grep "." file.txt

# Find TODO comments recursively
grep -rn "TODO" .

# Find FIXME comments
grep -rn "FIXME" .

# Find function definitions (simple example)
grep -E "^def |^function " file.py

# Find IP addresses (simple regex)
grep -E '[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+' file.txt

# Find email addresses (simple regex)
grep -E '[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}' file.txt

# Search compressed gzip files
zgrep "text" file.gz

# Search command history
history | grep "docker"

# Search running processes
ps aux | grep "nginx"

# Search listening ports
ss -tulpn | grep ":80"

# Filter ls output
ls -la | grep ".log"

# Search environment variables
env | grep "PATH"

# Find failed logins
grep "Failed password" /var/log/auth.log

# Multiple alternatives
grep -E "error|warning|critical" logfile.log

# Search recursively and ignore case
grep -ri "password" .

# Search recursively with line numbers and color
grep -rin --color=auto "config" .


### REGEX PATTERNS

^text          # starts with "text"
text$          # ends with "text"
^$             # empty line
.              # any single character
.*             # anything
[a-z]          # lowercase letter
[A-Z]          # uppercase letter
[0-9]          # digit
[a-zA-Z0-9]    # alphanumeric
[^0-9]         # not a digit
foo|bar        # foo OR bar (grep -E)
colou?r        # color/colour (grep -E)
ab*c           # ac, abc, abbc...
ab+c           # abc, abbc... (grep -E)
ab?c           # ac or abc (grep -E)


### COMMON COMBOS

# Find all Python TODOs
grep -rn --include="*.py" "TODO" .

# Find API keys
grep -rE "(api[_-]?key|secret|token)" .

# Find passwords
grep -ri "password" .

# Find imports
grep -r "^import " .

# Find all HTTP URLs
grep -rhoE "https?://[^ ]+" .

# Find UUIDs
grep -E "[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}" file.txt

# Count occurrences recursively
grep -roh "ERROR" . | wc -ld