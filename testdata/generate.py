import os
import random
import string
import sys

def random_line(length=50):
    return ''.join(random.choices(string.ascii_letters + string.digits, k=length))

def generate_file(filename, target_size):
    written = 0
    unique_lines = set()

    with open(filename, 'w') as f:
        while written < target_size:
            line = random_line(50)

            unique_lines.add(line.strip())

            if written + len(line) > target_size:
                remaining = target_size - written
                f.write(line[:remaining])
                written += remaining
                break

            f.write(line)
            written += len(line)

    print(f"Generated {filename} with size {os.path.getsize(filename)} bytes")


generate_file("test1.txt", 20 * 1024)
