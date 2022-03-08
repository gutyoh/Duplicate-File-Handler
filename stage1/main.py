import os
import sys

arg = sys.argv

print("The ARGUMENTS ARE: ", arg)

print(os.path.join('module', 'root_folder'))

if len(arg) == 1:
    print("Directory is not specified")
else:
    for root, dirs, files in os.walk(arg[1], topdown=True):
        for name in files:
            print(os.path.join(root, name))