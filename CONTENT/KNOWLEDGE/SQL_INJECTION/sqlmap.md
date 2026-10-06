# SQLMAP

## With Burpsuite

> Steps:
> - determine testing paratmeter e.g. id=1
> - intercept a request with burp, copy request text, paste into editor
> - add asterisk to testing parameter e.g. id=1*
> - save as req.txt

### Command

    sqlmap -r req.txt --batch --dump