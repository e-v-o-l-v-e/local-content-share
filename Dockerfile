FROM python

WORKDIR /app

COPY . .

RUN mkdir data && \
    pip install flask flask-markdown uuid pygments pygments-style-github

ENV FLASK_APP=app.py
EXPOSE 5000

CMD ["flask", "run", "--host=0.0.0.0"]
