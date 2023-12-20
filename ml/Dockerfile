FROM python:3.10

WORKDIR /ml

RUN pip install poetry

COPY . .

RUN poetry install --no-root

CMD ["poetry", "run", "python", "main.py"]