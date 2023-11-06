FROM python:3.10-slim

WORKDIR /app

ENV PYTHONDONTWRITEBYTECODE 1
ENV PYTHONUNBUFFERED 1

COPY requirements.txt requirements.txt
RUN pip install --upgrade pip &&\
    pip install -r requirements.txt

COPY . .

RUN apt-get update && apt-get install ffmpeg libsm6 libxext6  -y

CMD ["python", "./main.py"]