# **NIKTO**

## Install Nikto
```
sudo apt update && sudo apt install -y perl
git clone https://github.com/sullo/nikto
cd nikto/program
chmod +x ./nikto.pl
```

## Scan Host

    nikto -h <DOMAIN> -Tuning b
