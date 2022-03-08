import os
import sys
import hashlib

root_folder = "../../Duplicate File Handler"
same_size = dict()
sort_reverse = False

if not len(sys.argv) > 1:
    print("Directory is not specified")
else:

    print("Enter file format:")
    file_type = input()

    while True:
        print("Size sorting options:")
        print("1. Descending")
        print("2. Ascending")
        sort_option = input()

        if sort_option in ("1", "2"):
            if sort_option == "1":
                sort_reverse = True
            break
        else:
            print("Wrong option")

    root_folder = sys.argv[1]
    for root, dirs, files in os.walk(root_folder, topdown=True):
        for file in files:
            full_filename = f"{root}\\{file}"

            file_root, ext = os.path.splitext(full_filename)

            if file_type == '' or ext == file_type:
                if os.path.getsize(full_filename) in same_size:
                    same_size[os.path.getsize(full_filename)].append(full_filename)
                else:
                    same_size[os.path.getsize(full_filename)] = [full_filename]

    for size, li in sorted(same_size.items(), reverse=sort_reverse):
        print(f"{size} bytes")
        for item in li:
            print(item)
        print()

    while True:
        print("Check for duplicates?")
        answer = input()
        same_hash = {}

        if answer in ("yes", "no"):
            if answer == "yes":
                for size, li in sorted(same_size.items(), reverse=sort_reverse):
                    same_hash.setdefault(size, dict())

                    for item in li:
                        with open(item, "rb") as f:
                            h = hashlib.md5(f.read())
                        if h.hexdigest() in same_hash[size]:
                            same_hash[size][h.hexdigest()].append(item)
                        else:
                            same_hash[size][h.hexdigest()] = [item]
                counter = 1
                for size, hash_dict in sorted(same_hash.items(), reverse=sort_reverse):
                    print(f"{size} bytes")
                    for h, li in hash_dict.items():
                        if len(li) > 1:
                            print(f"Hash: {h}")
                            for el in li:
                                print(f"{counter}. {el}")
                                counter += 1
                    print()
            break
        else:
            print("Wrong option")