FROM python:alpine3.17

WORKDIR /app

COPY . .

RUN mkdir data && \
    pip install flask flask-markdown

ENV FLASK_APP=app.py

CMD ["flask", "run", "--host=0.0.0.0"]