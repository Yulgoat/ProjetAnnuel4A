# GUIDE PROJET MYCÉLIUM
![](https://lh7-us.googleusercontent.com/LRDQu3z4KFLsC_jDBDLNSd-HGQXYu8vzk9OZ9kp4AS5crYtpVo3KnwkXeMf1pCfibq7vwgwuK0bm1NXSBHlvac5GLYh30Br9X21tTeUCMQRZi4qBE0RPnoagNe8sehpeGnmCF9_p38g6v1_l6dc9FD8)
# I. Guide VM
**Ce guide explique comment mettre en place une ou plusieurs VM (Virtual Machine qui permettent d’émuler le comportement d’un cluster) et comment lancer ces VMs depuis un terminal.**

 ## 1.  Installations : 

- Installer QEMU (Virtualiseur de machine) :
`sudo apt-get update && sudo apt install qemu`

- Installer Virtual Manager (Interface graphique pour QEMU) :
`sudo apt install virt-manager`

- Redémarrer la machine.

## 2. Créer une VM :
Pour émuler le cluster, il faut créer une VM principale (main), et potentiellement plusieurs VMs qui agiront comme des agents.

  

Ouvrir Virtual Manager et suivre ces instructions :

  

-   Sélectionner QEMU/KVM.
    
-   Cliquer sur  ![](https://lh7-us.googleusercontent.com/MFFR6zugulGFwakyZqtaAZfpUZ4IcQ9pU5raZsSm6FMv5qiPfKMD9kRq7UdV6D23yRGW8WkwOWUAvHL12_grRi0tL35ZG4xzQTjMZMfpDbUdhFWCBhhlru7V-MLlzHkJwNFP2hXgSrbZPJa8zreOk0I) pour créer une nouvelle machine.
    
-   Sélectionner Média d’installation local (ISO).
    
-   Chercher le fichier ISO (ubuntu server 22.04 à télécharger sur le lien : [https://ubuntu.com/download/server](https://ubuntu.com/download/server)).
    
-   Ne pas faire de détection automatique, prendre Debian 11.
    
-   Choisir la mémoire et CPU souhaitée.
    
-   Choisir le volume de stockage souhaité en créant un espace de stockage personnalisé :
    

-   Pour le main : prendre 15 Go au moins (prendre le + possible).
    
-   Pour les agents : prendre 10 Go minimum.
    

  

Valider puis suivre la suite des étapes d’installation dans la VM (il faut tout prendre par défaut).

⚠️ Faire bien attention à utiliser tout l’espace alloué à la VM lors de la configuration. ⚠️

 ## 3. Se connecter en SSH à la VM :
    
Dans la VM, regarder quelle est l’adresse IP de la VM (192. …) :

     hostname -I

Faire cette commande sur son pc :

    ssh -p 22 -L 8082:localhost:8082 -L 8081:localhost:8080 -L 8086:localhost:8086 user@<Adresse IP VM>
(`ssh -p 22 user@<Adresse IP VM>` est suffisant, mais **-L \<portHôte\>:localhost:\<portVM\>** permet de lier le port de l'hôte au port de la VM (utile pour les futures applications que l’on va utiliser)

# II. Guide Kubernetes (K3S)
**Ce guide explique comment installer K3S sur les différentes VM pour leur permettre de communiquer et d’agir comme un cluster.**

 ## 1. Installation de K3S sur la main :
 - Téléchargement et installation de K3S :
`sudo apt install curl`
`sudo curl -sfL https://get.k3s.io | sh -`

- Obtention du token du main, nécessaire pour l’installation des agents (donc à sauvegarder dans un coin) :
    `sudo cat /var/lib/rancher/k3s/server/node-token`

- Vérifier que K3S est bien lancé :
 `sudo systemctl status k3s`

- Vérifier l’état des noeuds :
`sudo kubectl get nodes`

## 2. Installation de K3S sur les agents :
    
`sudo curl -sfL https://get.k3s.io | K3S_URL=https://<IP du main>:6443 K3S_TOKEN=<Token du main> sh -s - --with-node-id <Nom unique de l’agent>`

- Vérifier que K3S est bien lancé :
`sudo systemctl status k3s-agent`

- Pour regarder si l’agent est bien connecté, se connecter sur le main et lancer cette commande :
`sudo kubectl get nodes`

## 3. Commandes optionnelles :
Si l’agent n’arrive pas à se connecter, ces commandes peuvent être utiles.

- Enlever les pares-feu :
`sudo apt-get install ufw`
`sudo ufw disable`

- Ouvrir le port 6443 :
` sudo iptables -A INPUT -p tcp --port 6443 -j ACCEPT`

# III. Guide Docker
**Ce guide permet d’installer Docker et de le rendre utilisable.**

## 1. Installation de Docker :
- Installer snap :
`sudo apt install snapd`

- Éteindre la VM et la relancer.

- Installer la bonne version de Docker :
`sudo snap install --channel=core18 docker`

- Se donner les permissions pour utiliser docker :
` sudo groupadd docker`
`sudo usermod -aG docker \<user\>`

- Éteindre la VM et la relancer.

 ## 2. Commandes utiles :
- Donne la liste des images Docker :
`docker images`

- Supprimer les conteneurs arrêtés :
`docker system prune --all`

- Supprimer une image docker :
`docker rmi <id de l’image>`
