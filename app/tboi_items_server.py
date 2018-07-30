#!/usr/bin/python

from flask import Flask, jsonify, request, render_template
#from lxml import html
import requests
import socket
import random
from tboi_items_ref import tboi_items

app = Flask(__name__)

version=0.1

def get_item_from_wiki():
    item_box = {}
    base_url = "https://bindingofisaacrebirth.gamepedia.com"
    item = random.choice(tboi_items)
    url = "{}/{}".format(base_url, item.replace(' ','_'))
    #r = requests.get(url)
    #tree = html.fromstring(r.content)
    item_box['ItemName'] = item
    item_box['ItemURL'] = url
    item_box['ItemType'] = "Afterbirth+"
    item_box['ItemId'] = False
    
    return item_box

@app.route('/version', methods=['GET'])
def show_version():
    """Endpoint to return app version"""
    response = "App Version: " + str(version)
    return response

@app.route('/', methods=['GET'])
def root():
    # Default. This route returns a TBOI item 
    item = get_item_from_wiki()
    return render_template("index.html.j2", item=item, container=str(socket.gethostname()))

if __name__ == "__main__":
    app.run(host='0.0.0.0')
