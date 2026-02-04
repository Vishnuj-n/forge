# Template: my-first-template

## Overview

This is a Forge template for bootstrapping my-first-template projects.

## What This Template Does

1. Initializes a git repository
2. Copies template files
3. Applies custom patches

## Getting Started

### Edit the template

1. template.yaml — Define commands and file operations
2. files/ — Add files/directories to copy into projects
3. patches/ — Add content to append to existing files

### Test the template

Run: forge test templates/my-first-template

### Use the template

Run: forge init templates/my-first-template ./my-project

## Prerequisites

- Git (for git init commands)
- Any other tools used in the commands section

## Customization

### Add files to copy
1. Add files to the files/ directory
2. In template.yaml, add to files.copy:
   copy:
     - files/README.md
     - files/config.json

### Add append patches
1. Create patch files in patches/
2. In template.yaml, add to files.append:
   append:
     - target: ".gitignore"
       source: "patches/gitignore.append"

### Add commands
In template.yaml, add to the commands section:
commands:
  - cmd: ["git", "init"]
  - cmd: ["echo", "Hello from template!"]

## Tips

- Commands are executed in token-array format (no shell strings)
- Target files for append operations must be created by commands or copy operations
- Use forge test to debug templates without committing
- Keep commands simple and deterministic

## For More Information

See the main Forge documentation:
- README.md — Project overview
- TEMPLATE-GUIDE.md — Complete template guide
