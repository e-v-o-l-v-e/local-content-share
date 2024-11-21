import os
import uuid
import markdown
from flask import Flask, render_template, request, redirect, url_for, make_response, send_file

app = Flask(__name__)
app.jinja_env.filters["markdown"] = markdown.markdown

@app.route("/")
def home():
    entries = []
    links = []
    files = []
    for filename in os.listdir("data"):
        if filename == "temporary-print":
            continue
        if filename.startswith("link-"):
            with open(os.path.join("data", filename)) as f:
                link = f.read()
            links.append((filename,link))
            continue
        if filename.startswith("file-"):
            files.append(filename)
            continue
        with open(os.path.join("data", filename)) as f:
            text = f.readlines()
        # if len(text) > 4:
        #     text = '\n'.join(text[:4])
        # else:
        #     text = '\n'.join(text)
        text = text[0]
        entries.append((filename, text))
    return render_template("index.html", entries=entries, files=files, links=links, markdown=markdown)

@app.route("/print")
def printit():
    return render_template("print.html", markdown=markdown)

@app.route("/darkprint")
def darkprintit():
    return render_template("darkprint.html", markdown=markdown)

@app.route("/print_text", methods=["POST"])
def print_text():
    text = request.form["demo-message"]
    filepath = os.path.join("data", "temporary-print")
    with open(filepath, "w") as f:
        f.write(text)
    return redirect(url_for("render_text", filename="temporary-print"))

@app.route("/darkprint_text", methods=["POST"])
def darkprint_text():
    text = request.form["demo-message"]
    filepath = os.path.join("data", "temporary-print")
    with open(filepath, "w") as f:
        f.write(text)
    return redirect(url_for("render_dark_text", filename="temporary-print"))

@app.route("/submit_text", methods=["POST"])
def submit_text():
    text = request.form["demo-message"]
    if len(text) == 0:
        return redirect(url_for("home"))
    filename = str(uuid.uuid4())
    filepath = os.path.join("data", filename)
    with open(filepath, "w") as f:
        f.write(text)
    return redirect(url_for("home"))

@app.route("/submit_link", methods=["POST"])
def submit_link():
    text = request.form["demo-message-2"]
    filename = str(uuid.uuid4())
    filepath = os.path.join("data", "link-" + filename)
    with open(filepath, "w") as f:
        f.write(text)
    return redirect(url_for("home"))

@app.route("/submit_file", methods=["POST"])
def submit_file():
    if "file" not in request.files:
        return redirect(url_for("home"))
    upfile = request.files['file']
    filename = str(uuid.uuid4())
    filepath = os.path.join("data", "file-" + filename[:10] + upfile.filename)
    upfile.save(filepath)
    return redirect(url_for("home"))

@app.route("/download/<filename>")
def download_file(filename):
    filepath = os.path.join("data", filename)
    return send_file(filepath, as_attachment=True)

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
    resp = make_response(md, 200)
    resp.mimetype = "text/plain"
    return resp

@app.route("/delete/<filename>")
def delete_text(filename):
    os.remove(os.path.join("data", filename))
    return redirect(url_for("home"))

@app.route("/count_files")
def count_files():
    num_files = len([f for f in os.listdir("data") if not f.startswith("temporary-print")])
    return str(num_files)
