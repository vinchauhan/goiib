FROM vinchauhan/my-iib:rc

RUN mkdir /home/iibuser/go && \
    sudo apt-get update && \
    sudo apt-get install -y git golang-1.10-go && \
    sudo rm -rf /var/lib/apt/lists/*

ENV GOPATH /home/iibuser/go

ENV PATH "$PATH:/usr/lib/go-1.10/bin:/opt/ibm/iib-10.0.0.10/server/bin"

USER iibuser

#VOLUME /Users/vchauhan/go/src/github.com/vinchauhan/goiib:/home/iibuser/go/src/github.com/vinchauhan/goiib
#RUN sudo ln -s usr/local/bin/iib_env.sh / # backwards compat

#ENTRYPOINT ["iib_env.sh"]

# docker run -d -e LICENCE=accept -p 4414:4414 -p 7800:7800 \
# -v /Users/vchauhan/go/src/github.com/vinchauhan/goiib:/home/iibuser/go/src/github.com/vinchauhan/goiib \
# vinchauhan/goiib:dev