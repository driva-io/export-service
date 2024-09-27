import base64

FILE = ".env"
DOKKU_APP = "export-service-only"

f = open(FILE, "r")

lines = f.readlines()

envs = []
for line in lines:
    if line[0] == "#":
        continue

    key = line.split("=")[0]
    value = line.split("=")[1]

    enc = key + "=" + base64.b64encode(value.replace("\n", "").encode()).decode()
    envs.append(enc)

f.close()

print(f"dokku config:set {DOKKU_APP} --encoded " + " ".join(envs))