FROM fedora:32
ENV D=/Generator

WORKDIR $D
COPY . $D/

RUN dnf install -y gcc make qpid-proton-c-devel
RUN make && mv gen /gen


USER 1001

ENTRYPOINT ["/gen"]


