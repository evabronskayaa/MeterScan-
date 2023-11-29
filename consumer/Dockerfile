FROM python:3.10
 
WORKDIR /consumer

RUN pip install poetry

COPY ./consumer /consumer

RUN poetry install
 
CMD ["poetry", "run", "python", "main.py"]