import os
import uuid
import markdown
from pygments import highlight
from pygments.lexers import get_lexer_by_name
from pygments.formatters import html
from flask import Flask, render_template, request, redirect, url_for

app = Flask(__name__)
app.jinja_env.filters["markdown"] = markdown.markdown

@app.route("/")
def home():
    filenames = []
    entries = []
    for filename in os.listdir("data"):
        with open(os.path.join("data", filename)) as f:
            text = f.readlines()
        # if len(text) > 4:
        #     text = '\n'.join(text[:4])
        # else:
        #     text = '\n'.join(text)
        text = text[0]
        filenames.append(filename)
        entries.append((filename, text))
    return render_template("index.html", entries=entries, filenames=filenames, markdown=markdown)

@app.route("/print")
def printit():
    return render_template("print.html", markdown=markdown)

@app.route("/print_text", methods=["POST"])
def print_text():
    text = request.form["demo-message"]
    filepath = os.path.join("data", "temporary-print")
    with open(filepath, "w") as f:
        f.write(text)
    return redirect(url_for("render_text", filename="temporary-print"))

@app.route("/submit_text", methods=["POST"])
def submit_text():
    text = request.form["demo-message"]
    filename = str(uuid.uuid4())
    filepath = os.path.join("data", filename)
    with open(filepath, "w") as f:
        f.write(text)
    return redirect(url_for("home"))

@app.route("/render/<filename>")
def render_text(filename):
    with open(os.path.join("data", filename)) as f:
        md = f.read()
        htmlfmd = markdown.markdown(md, extensions=['fenced_code', 'tables'])
    return render_template('render.html', html=htmlfmd)

@app.route("/render_dark/<filename>")
def render_dark_text(filename):
    with open(os.path.join("data", filename)) as f:
        md = f.read()
        htmlfmd = markdown.markdown(md, extensions=['fenced_code'])
    return render_template('render-dark.html', html=htmlfmd)

@app.route("/view/<filename>")
def view_text(filename):
    with open(os.path.join("data", filename)) as f:
        md = f.read()
    return "<pre>" + md + "</pre>"

@app.route("/delete/<filename>")
def delete_text(filename):
    os.remove(os.path.join("data", filename))
    return redirect(url_for("home"))
