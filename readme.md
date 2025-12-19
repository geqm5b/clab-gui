# Instalation NetosLab

## redme for testing purposes.

## Docker

### Delete older versions

```bash
sudo apt remove $(dpkg --get-selections docker.io docker-compose docker-doc podman-docker containerd runc | cut -f1)
```

### Add Docker's official GPG key

```bash
sudo apt update
sudo apt install ca-certificates curl
sudo install -m 0755 -d /etc/apt/keyrings
sudo curl -fsSL https://download.docker.com/linux/debian/gpg -o /etc/apt/keyrings/docker.asc
sudo chmod a+r /etc/apt/keyrings/docker.asc
```

### Add the repository to Apt sources

```bash
sudo tee /etc/apt/sources.list.d/docker.sources <<EOF
Types: deb
URIs: https://download.docker.com/linux/debian
Suites: $(. /etc/os-release && echo "$VERSION_CODENAME")
Components: stable
Signed-By: /etc/apt/keyrings/docker.asc
EOF
```

```bash
sudo apt update
```

### Instalation

```bash
sudo apt install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
```

### docker Status

```bash
sudo systemctl status Docker
```

### Start docker

```bash
sudo systemctl start Docker
```

### Non-root config

```bash
sudo groupadd docker
sudo usermod -aG docker $USER
newgrp Docker
```

```bash
sudo chown "$USER":"$USER" /home/"$USER"/.docker -R
sudo chmod g+rwx "$HOME/.docker" -R
```

### test

```bash
docker run hello-world
```

---

## Containerlab

```bash
bash -c "$(curl -sL https://get.containerlab.dev)"
```

---

## build docker images 

```bash
docker build -f Dockerfile.dns -t netoslab-dns .
```

```bash
docker build -f Dockerfile.dns -t netoslab-dhcp .
```

```bash
docker build -f Dockerfile.dns -t netoslab-cliente .
```

## start the app ()

```bash
sudo docker compose up --build -d
```