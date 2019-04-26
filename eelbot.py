'''
    EelBot
'''

import argparse
import configparser
import time

CONFIG = {}


def poll():
    run = True

    while run:
        print("Polled")
        time.sleep(CONFIG["POLL_TIME"])


parser = argparse.ArgumentParser(description="Discord bot, listener server")
parser.add_argument("-c", "--config", default="config.ini")

config = configparser.ConfigParser()

if __name__ == '__main__':
    print("Starting EelBot Discord bot")

    args = parser.parse_args()

    config.read(args.config)

    CONFIG["POLL_TIME"] = config["GENERAL"].get("POLL_TIME", 1)

    poll()
