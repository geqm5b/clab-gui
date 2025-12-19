# Instalación NetosLab

## Readme inicial para test

## Docker

### Eliminar versiones anteriores

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

### Instalación

```bash
sudo apt install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
```

### Status

```bash
sudo systemctl status Docker
```

### Start

```bash
sudo systemctl start Docker
```

### Non-root

```bash
sudo groupadd docker
sudo usermod -aG docker $USER
newgrp Docker
```

```bash
sudo chown "$USER":"$USER" /home/"$USER"/.docker -R
sudo chmod g+rwx "$HOME/.docker" -R
```

### Prueba

```bash
docker run hello-world
```

---

## Containerlab

```bash
bash -c "$(curl -sL https://get.containerlab.dev)"
```

---

## Crear las imágenes de prueba

```bash
docker build -f Dockerfile.dns -t netoslab-dns:latest .
```

```bash
docker build -f Dockerfile.dns -t netoslab-dhcp:latest .
```
