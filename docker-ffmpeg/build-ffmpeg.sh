#!/usr/bin/env bash
# Install required dependencies
apt-get update
apt-get install -y build-essential autoconf \
    automake cmake libtool git \
    nasm yasm libass-dev libfreetype6-dev libsdl2-dev p11-kit \
    libva-dev libvdpau-dev libvorbis-dev libxcb1-dev libxcb-shm0-dev \
    libxcb-xfixes0-dev pkg-config texinfo wget zlib1g-dev libchromaprint-dev \
    frei0r-plugins-dev gnutls-dev ladspa-sdk libcaca-dev libcdio-paranoia-dev \
    libcodec2-dev libfontconfig1-dev libfreetype6-dev libfribidi-dev libgme-dev \
    libgsm1-dev libjack-dev libmodplug-dev libmp3lame-dev libopencore-amrnb-dev \
    libopencore-amrwb-dev libopenjp2-7-dev libopenmpt-dev libopus-dev \
    libpulse-dev librsvg2-dev librubberband-dev librtmp-dev libshine-dev \
    libsmbclient-dev libsnappy-dev libsoxr-dev libspeex-dev libssh-dev \
    libtesseract-dev libtheora-dev libtwolame-dev libv4l-dev libvo-amrwbenc-dev \
    libvorbis-dev libvpx-dev libwavpack-dev libwebp-dev libx265-dev \
    libxvidcore-dev libxml2-dev libzmq3-dev libzvbi-dev liblilv-dev \
    libopenal-dev opencl-dev libjack-dev unzip imagemagick curl \
    libbluray-dev

# x264
DIR=$(mktemp -d) && cd ${DIR} && \
              git -C x264 pull 2> /dev/null || git clone --depth 1 https://code.videolan.org/videolan/x264.git && \
              cd x264 && \
              ./configure --prefix="/usr/local" --bindir="/usr/local/bin" --enable-static --enable-pic && \
              make -j && \
              make install && \
              make distclean && \
              rm -rf ${DIR}

# fdk-aac
DIR=$(mktemp -d) && cd ${DIR} && \
              curl -s https://codeload.github.com/mstorsjo/fdk-aac/tar.gz/v${FDKAAC_VERSION} | tar zxvf - && \
              cd fdk-aac-${FDKAAC_VERSION} && \
              autoreconf -fiv && \
              ./configure --prefix="/usr/local" --enable-static --enable-pic && \
              make -j && \
              make install && \
              make distclean && \
              rm -rf ${DIR}

# ffmpeg
DIR=$(mktemp -d) && cd ${DIR} && \
              curl -Os http://ffmpeg.org/releases/ffmpeg-${FFMPEG_VERSION}.tar.gz && \
              tar xzvf ffmpeg-${FFMPEG_VERSION}.tar.gz && \
              cd ffmpeg-${FFMPEG_VERSION} && \
              ./configure --prefix="/usr/local" --extra-cflags="-I/usr/local/include" --extra-ldflags="-L/usr/local/lib" --bindir="/usr/local/bin" \
              --enable-gpl --enable-version3 \
			  --enable-shared --enable-small --enable-avisynth --enable-chromaprint \
			  --enable-frei0r --enable-gmp --enable-gnutls --enable-ladspa \
			  --enable-libass --enable-libcaca --enable-libcdio \
			  --enable-libcodec2 --enable-libfontconfig --enable-libfreetype \
			  --enable-libfribidi --enable-libgme --enable-libgsm --enable-libjack \
			  --enable-libmodplug --enable-libmp3lame --enable-libopencore-amrnb \
			  --enable-libopencore-amrwb --enable-libopencore-amrwb \
			  --enable-libopenjpeg --enable-libopenmpt --enable-libopus --enable-libpulse \
			  --enable-librsvg --enable-librubberband --enable-librtmp --enable-libshine \
			  --enable-libsnappy --enable-libsoxr --enable-libspeex \
			  --enable-libssh --enable-libtesseract --enable-libtheora \
			  --enable-libtwolame --enable-libv4l2 --enable-libvo-amrwbenc \
			  --enable-libvorbis --enable-libvpx --enable-libwavpack --enable-libwebp \
			  --enable-libx264 --enable-libx265 --enable-libxvid --enable-libxml2 \
			  --enable-libzmq --enable-libzvbi --enable-lv2 \
			  --enable-openal --enable-opencl --enable-opengl --enable-libdrm \
			  --enable-nonfree --enable-libfdk-aac --enable-libbluray \
			  --enable-encoder=png --enable-decoder=png --enable-libx264 && \
              make -j && \
              make install && \
              make distclean && \
              hash -r && \
              cd tools && \
              make qt-faststart && \
              cp qt-faststart /usr/local/bin && \
              rm -rf ${DIR}

# Update library
ldconfig

apt-get remove -y --purge git build-essential yasm nasm autoconf libtool pkg-config libspeex-dev curl zlib1g-dev mercurial cmake cmake-curses-gui
apt-get clean