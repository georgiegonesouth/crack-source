# **theHarvester**

## Install

```
curl -LsSf https://astral.sh/uv/install.sh | sh
git clone https://github.com/laramies/theHarvester
cd theHarvester
uv sync
uv run theHarvester
```

## Query for Domain Info

```
uv run theHarvester -d <DOMAIN> -b <SOURCE>
```

### Sources

<!-- token-table:SOURCE -->

| Column 1         | Column 2           | Column 3         |
| ---------------- | ------------------ | ---------------- |
| baidu            | bevigil            | bitbucket        |
| brave            | bufferoverun       | builtwith        |
| censys           | certspotter        | chaos            |
| commoncrawl      | criminalip         | crtsh            |
| dehashed         | dnsdumpster        | duckduckgo       |
| dymo             | fofa               | fullhunt         |
| github-code      | gitlab             | hackertarget     |
| haveibeenpwned   | hudsonrock         | hunter           |
| hunterhow        | intelx             | leakix           |
| leaklookup       | mojeek             | netlas           |
| onyphe           | otx                | pentesttools     |
| projectdiscovery | rapiddns           | robtex           |
| rocketreach      | securityscorecard  | securityTrails   |
| sherlockeye      | shodan             | shodanInternetDB |
| subdomaincenter  | subdomainfinderc99 | thc              |
| threatcrowd      | tomba              | urlscan          |
| venacus          | virustotal         | waybackarchive   |
| whoisxml         | windvane           | yahoo            |
| zoomeye          |                    |                  |
