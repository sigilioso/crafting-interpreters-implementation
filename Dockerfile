FROM dart:2

ADD . /book
WORKDIR /book
RUN cd tool && dart pub get
