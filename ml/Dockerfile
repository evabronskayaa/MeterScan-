FROM python:3.10

WORKDIR /ml

RUN apt-get update && apt-get install ffmpeg libsm6 libxext6  -y

RUN pip install poetry

COPY . .

RUN poetry install --no-root

CMD ["poetry", "run", "python", "main.py"]