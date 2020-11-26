FROM gitpod/workspace-full

RUN brew install zsh
RUN sudo apt-key adv --keyserver keyserver.ubuntu.com --recv-key C99B11DEB97541F0
RUN sudo apt-add-repository https://cli.github.com/packages
RUN sudo apt update
RUN sudo apt install gh
RUN sh -c "$(curl -fsSL https://raw.github.com/ohmyzsh/ohmyzsh/master/tools/install.sh)"
RUN npm install -g auto-changelog

# Install custom tools, runtimes, etc.
# For example "bastet", a command-line tetris clone:
# RUN brew install bastet
#
# More information: https://www.gitpod.io/docs/config-docker/
