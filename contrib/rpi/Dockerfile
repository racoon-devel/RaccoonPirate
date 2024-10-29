FROM ghcr.io/tttapa/docker-arm-cross-toolchain:__TARGET_TRIPLET__-0.4.0

RUN echo "develop" | sudo -S sh -c "echo '%sudo ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers"

RUN sudo -S apt-get update && \
    sudo apt-get -y upgrade && \
    sudo apt-get -y install golang-go libfuse-dev && \
    sudo apt-get clean autoclean
