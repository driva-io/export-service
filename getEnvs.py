import base64

FILE = ".env"
DOKKU_APP = "export-service-only"

f = open(FILE, "r")

ls = f.readlines()

envs = []
for l in ls:
    if l[0] != "#" and "=" in l:
        key = l.split("=")[0]
        value = l.split("=")[1]
        enc = (
            key
            + "="
            + base64.b64encode(value.replace("\n", "").encode("ascii")).decode("ascii")
        )
        envs.append(enc)

f.close()

print(f"dokku config:set {DOKKU_APP} --encoded " + " ".join(envs))