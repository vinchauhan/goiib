FROM vinchauhan/my-iib:rc

MAINTAINER Chauhan.Vineet@gmail.com

COPY $GOPATH/bin/goiib /build

#Add the build tool to the path
ENV PATH="/build/:${PATH}"



