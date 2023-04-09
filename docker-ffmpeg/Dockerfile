FROM ubuntu:bionic

ENV FFMPEG_VERSION  4.2.2
ENV FDKAAC_VERSION  0.1.6
ENV SRC             /usr/local
ENV LD_LIBRARY_PATH ${SRC}/lib
ENV PKG_CONFIG_PATH ${SRC}/lib/pkgconfig

COPY build-ffmpeg.sh /build-ffmpeg.sh
RUN /build-ffmpeg.sh && /usr/local/bin/ffmpeg -buildconf && wget http://zebulon.bok.net/Bento4/binaries/Bento4-SDK-1-5-1-624.x86_64-unknown-linux.zip -O bento.zip && unzip -x bento.zip && mv Bento4*/bin/* /usr/bin && chmod +x /usr/bin/mp4* && rm -rf Bento4* bento.zip