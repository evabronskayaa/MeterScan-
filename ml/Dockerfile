FROM python:3.10
 
WORKDIR /app

ENV PYTHONDONTWRITEBYTECODE 1
ENV PYTHONUNBUFFERED 1

RUN pip install poetry

COPY . .

RUN poetry install
 
CMD ["poetry", "run", "python", "main.py"]