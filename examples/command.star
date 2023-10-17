""" 
Run an arbitrary command in the shell
"""

name = "command_playbook"

command.run("pwd")
command.run("ls", ["-lah"])


