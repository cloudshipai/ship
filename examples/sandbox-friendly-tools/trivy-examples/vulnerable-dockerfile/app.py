#!/usr/bin/env python3

import os
import subprocess
import pickle

# Example vulnerable Python application for security scanning
def vulnerable_function():
    # SQL Injection vulnerability
    user_input = input("Enter username: ")
    query = f"SELECT * FROM users WHERE username = '{user_input}'"
    
    # Command injection vulnerability
    filename = input("Enter filename: ")
    subprocess.call(f"ls {filename}", shell=True)
    
    # Deserialization vulnerability
    data = input("Enter serialized data: ")
    pickle.loads(data.encode())
    
    # Path traversal vulnerability
    file_path = input("Enter file path: ")
    with open(file_path, 'r') as f:
        content = f.read()
    
    return content

if __name__ == "__main__":
    print("Vulnerable Application")
    vulnerable_function()