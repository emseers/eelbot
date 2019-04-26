import json
import random

def printJoke(param, param2 = None):
    try:
        if param == "me":
            with open('jokes.json') as jsonFile:
                data = json.load(jsonFile)
                jokes = data['jokes']
                i = random.randint(0, len(jokes))
                if jokes[i]['type'] == "badjokeeel":
                    return jokes[i]['top'] + "\n" + jokes[i]['bottom']
    except:
        return -1

def printEel(param, param2 = None):
    try:
        if param == "me":
            return "An eel pic..."
        elif param == "bomb":
            return param2 + " eel pics."
    except:
        return -1

options = {
    "badjoke" : printJoke,
    "eel"     : printEel
}

def processCmd(cmd):
    words = cmd.strip().split(' ')
    if words[0][0] == '/':
        words[0] = words[0][1:]
        try:
            return options[words[0]](words[1], words[2] if (len(words) > 2) else None)
        except:
            return -1