from tempfile import mkstemp
from shutil import move, copymode
from os import fdopen
import re, os

regex = "\{\s(\$.*?)\s\}"

fh, abs_path = mkstemp()
with fdopen(fh,'w') as new_file:
    with open('./krakend.json') as old_file:
        for line in old_file:
            exists = re.search(regex, line)
            if not exists:
                new_file.write(line)
                continue

            pattern = exists.group()
            env_var_name = pattern[3:len(pattern)-2]
            env_var_value = os.getenv(env_var_name)

            if env_var_value is None:
                raise Exception(f"the {env_var_name} secret is not set")

            new_file.write(line.replace(pattern, env_var_value))

    copymode('krakend.json', abs_path)
    move(abs_path, './replaced_with_env_values.json')

