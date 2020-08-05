import re

with open('../krakend/krakend.json', 'r') as config:
    for line in config:
        x = re.search("\{\s\$.*\s\}", line)
        print(x)