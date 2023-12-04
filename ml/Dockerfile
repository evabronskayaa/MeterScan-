FROM python:3.10
 
WORKDIR /consumer

ENV PYTHONDONTWRITEBYTECODE 1
ENV PYTHONUNBUFFERED 1

RUN pip install poetry

COPY ./consumer /consumer

RUN poetry install
 
CMD ["poetry", "run", "python", "main.py"]