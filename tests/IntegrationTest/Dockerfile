FROM python:3.6.5-alpine

WORKDIR /test

COPY ./requirements.txt ./requirements.txt

RUN pip install -r requirements.txt

COPY . .

CMD [ "python", "-m", "unittest", "discover", "-v" ]