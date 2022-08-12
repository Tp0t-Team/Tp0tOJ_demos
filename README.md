# Tp0tOJ_demos
this is a repository for some ctf game examples


```yaml
name: dynamic_flag              # challenge name
category: PWN                   # support [MISC|RE|PWN|WEB|CRYPTO|HARDWARE|RW]
score:
  baseScore: 1000               # base score
  dynamic: true                 # if the score change with solved number
flag: 
  value: xxxxxxxxxx             # means that random dynamic flag length is 10
  type: 3                       # 0-Single 1-Multiple 2-Regexp 3-Dynamic
description: "description" 
externalLink: ["http://link"]
singleton: false                # false means this challenge will give every ctfer a container

# below is no need for singleton challenge
nodeConfig:
  - name: "pwn1"                # give this name same as your uploaded image
    image: "pwn1"               # give this name same as your uploaded image
    servicePorts:               # default & DON'T CHANGE IT
      - name: http              # default & DON'T CHANGE IT
        protocol: TCP           # default & DON'T CHANGE IT
        external: 8888          # default & DON'T CHANGE IT
        internal: 8888          # default & DON'T CHANGE IT
        pod: 0                  # default & DON'T CHANGE IT

```