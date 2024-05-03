# GUIDE MYCÉLIUM

![](https://lh7-us.googleusercontent.com/LRDQu3z4KFLsC_jDBDLNSd-HGQXYu8vzk9OZ9kp4AS5crYtpVo3KnwkXeMf1pCfibq7vwgwuK0bm1NXSBHlvac5GLYh30Br9X21tTeUCMQRZi4qBE0RPnoagNe8sehpeGnmCF9_p38g6v1_l6dc9FD8)

**Table of Contents**

[TOCM]

[TOC]

# Guide VM

**Ce guide explique comment mettre en place une ou plusieurs VM (Virtual Machine qui permettent d'émuler le comportement d'un cluster et comment lancer ces VMs depuis un terminal.**

## Installations:

Installer QEMU (Virtualiseur de machine) :
```sh 
sudo apt-get update && sudo apt install qemu
```
Installer Virtual Manager (Interface graphique pour QEMU) :
```sh 
sudo apt install virt-manager
```
Redémarrer la machine.

## Créer une VM

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

## Se connecter en SSH à la VM

Dans la VM, regarder quelle est l’adresse IP de la VM (192. …) :
```sh 
hostname -I 
```

Faire cette commande sur son pc :
```sh 
ssh -p 22 -L 8082:localhost:8082 -L 8081:localhost:8080 -L 8086:localhost:8086 user@<Adresse IP VM>
```

(`ssh -p 22 user@<Adresse IP VM>` est suffisant, mais **-L \<portHôte\>:localhost:\<portVM\>** permet de lier le port de l'hôte au port de la VM (utile pour les futures applications que l'on va utiliser)

# Guide Kubernetes (K3S)

Ce guide explique comment installer K3S sur les différentes VM pour leur permettre de communiquer et d’agir comme un cluster.

## Installation de K3S sur la main :

Téléchargement et installation de K3S : 
```sh
sudo apt install curl
sudo curl -sfL https://get.k3s.io | sh -
```
Obtention du token du main, nécessaire pour l’installation des agents (donc à sauvegarder dans un coin) : 
```sh
sudo cat /var/lib/rancher/k3s/server/node-token
```
Vérifier que K3S est bien lancé : 
```sh
sudo systemctl status k3s
```
Vérifier létat des noeuds :
```sh
sudo kubectl get nodes
```

## Installation de K3S sur les agents :

```sh
sudo curl -sfL https://get.k3s.io | K3S_URL=https://<IP du main>:6443 K3S_TOKEN=<Token du main> sh -s - --with-node-id <Nom unique de l’agent>
```

Vérifier que K3S est bien lancé : 
```sh
sudo systemctl status k3s-agent	
```

Pour regarder si l’agent est bien connecté, se connecter sur le main et lancer cette commande :
```sh
sudo kubectl get nodes
```

## Commandes optionnelles :

Si l'agent n'arrive pas à se connecter, ces commandes peuvent être utiles.

Enlever les pares-feu :
```sh
sudo apt-get install ufw
sudo ufw disable
```

Ouvrir le port 6443 :
```sh
sudo iptables -A INPUT -p tcp --port 6443 -j ACCEPT
```

# Guide Docker

**Ce guide permet d'installer Docker et de le rendre utilisable.**

## Installation de Docker :
Installer snap :
```sh
sudo apt install snapd
```
Éteindre la VM et la relancer.
Installer la bonne version de Docker :
```sh
sudo snap install --channel=core18 docker
```
Se donner les permissions pour utiliser docker : 
```sh
sudo groupadd docker
sudo usermod -aG docker <user>
```
Éteindre la VM et la relancer.

## Commandes utiles :

Donne la liste des images Docker :
```sh
docker images
```
Supprimer les conteneurs arrêtés : 
```sh
docker system prune --all
```
Supprimer une image docker :
```sh
docker rmi <id de l’image>
```

# Guide Open Faas

Ce guide permet d'installer Open Faas ainsi que de créer et modifier des fonctions.
Pour réaliser ce guide, Docker doit être installé sur la VM. Pour vérifier, lancez la commande :
```sh
docker -v
```

## Installations :

Installer Arkade :
```sh
curl -sLS https://get.arkade.dev | sudo -E sh
```

Installer faas-cli :
```sh
arkade get faas-cli
sudo mv /home/<user>/.arkade/bin/faas-cli /usr/local/bin/ 
```
(./root/arkade en cas d’installation sur le VPS)

Installer helm :
```sh
arkade get helm
sudo mv /home/<user>/.arkade/bin/helm /usr/local/bin/
```
ou

```sh
curl -sSLf https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3 | bash
```

Cloner le repository git faas-netes : 
```sh
sudo apt install git
git clone https://github.com/openfaas/faas-netes.git
```

Créer un fichier openfaas.yaml :
```sh
cd faas-netes
helm template \
  openfaas chart/openfaas/ \
  --namespace openfaas \
  --set basic_auth=true \
  --set functionNamespace=openfaas-fn > openfaas.yaml
```

Installer Open Faas avec kubectl :
```sh
sudo chmod 644 /etc/rancher/k3s/k3s.yaml
```

Pour effectuer ça à chaque démarrage :
```sh
sudo nano ~/.bashrc
```

Ajouter la ligne “sudo chmod 644 /etc/rancher/k3s/k3s.yaml” à la fin du fichier puis entrer la ligne de commande :
```sh
kubectl apply -f namespaces.yml,openfaas.yaml
```

Obtenir mot de passe admin :
```sh
PASSWORD=$(k3s kubectl -n openfaas get secret basic-auth -o jsonpath="{.data.basic-auth-password}" | base64 --decode) && \
echo "OpenFaaS admin password: $PASSWORD"
```
On peut maintenant quitter le dossier faas-netes.

## Préparation du cluster :

Redirige les requêtes vers http://localhost:8080 :

```sh
kubectl port-forward -n openfaas svc/gateway 8080:8080 &
```

Pour se connecter sur la page web, l’username est admin et le mot de passe est le mot de passe obtenu dans la partie précédente.

Se connecter à faas-cli :
```sh
faas-cli login --password $PASSWORD
```

En cas d’erreur : Unable to read /etc/rancher/k3s/k3s.yaml : 
```sh
export KUBECONFIG=~/.kube/config
mkdir ~/.kube 2> /dev/null
sudo k3s kubectl config view --raw > "$KUBECONFIG"
chmod 600 "$KUBECONFIG"
sudo nano ~/.bashrc
```

Et ajouter à la fin la ligne : export KUBECONFIG=~/.kube/config

## Création et modification de fonctions OpenFaas :

- **Création de la fonction :**

(⚠️ le code de la fonction apparaîtra dans le répertoire actuel du terminal)

Télécharger template de création de fonction en go (pull python3 pour Python) :
```sh
faas-cli template store pull golang-middleware
```

Créer une nouvelle fonction :
```sh
faas-cli new <nom fonction> --lang golang-middleware
```

Il faut maintenant créer une image Docker de la fonction. Tout d’abord, il faut avoir un compte Docker Hub (comme Github mais pour Docker).

⇒ Compte docker créé spécialement pour le projet Mycélium :
 user : myceliumir   MDP : MyceliumIR

Ensuite, il faut ouvrir le fichier test-function.yml et modifier le champ image: de cette façon :
```sh
image: <ton_pseudo_docker_hub>/test-function:latest
```

Se connecter à Docker Hub :
```sh
docker login --username <pseudo docker hub>
```

Build, Push sur Docker Hub et Deploy la fonction :
```sh
faas-cli up -f test-function.yml
```

Démarrer la fonction : 
```sh
faas-cli invoke -f test-function.yml test-function
```

- **Modification de la fonction :**

Modifiez les fichiers comme vous voulez. Le code de la fonction en Go se situe dans le fichier test-function/handler.go.

Supprime la fonction précédente : 
```sh
faas-cli remove test-function
```

Build, Push and Deploy la fonction : 
```sh
faas-cli up -f test-function.yml
```

Parfois le code est pas bien formaté et il faut utiliser une commande pour régler le problème :
```sh
gofmt -s -w
```

Pour pouvoir lancer cette commande, il faut que go soit installé sur la machine :
```sh
sudo apt-get -y install golang-go
```

Redémarrer la fonction : 
```sh
faas-cli invoke -f test-function.yml test-function
```

## Commandes utiles :

- **Open Faas :**

| Action | Commande |
| :-------------: | :-------------: |
| Voir si le cluster a bien démarré :  |  ` kubectl get deployments -n openfaas -l "release=openfaas, app=openfaas" `  |
| Vérifier si la gateway est bien déployée :  | `kubectl rollout status -n openfaas deploy/gateway`  |
| Voir les que la redirection est bien activée (une ligne avec Running devrait s'afficher) :  | `jobs`  |
| Détails sur un service: | `kubectl get service <name-of-service> -n openfaas` |
| Supprimer les pods inutiles :   |
| Affiche tous les pods du cluster :  | `kubectl get pods -A`  |
| Informations détaillées sur les noeuds :  | `kubectl describe nodes`  |
| Affiche des détails spécifique au pod XXX :  | `kubectl describe pod XXX -n openfaas` |
| Liste des déploiements du cluster :   | `kubectl get deployments -A` |
| Supprime le déploiement XXX :  | `kubectl delete deployment XXX -n openfaas` |

- **Fonctions Open Faas :**

| Action | Commande |
| :-------------: | :-------------: |
| Affiche tout ce qui a été print dans le programme (par fmt.Println()), utile pour débug :    |  `faas-cli list`  |
| Obtenir les détails sur une fonction : | `faas-cli describe <nom-fonction>` |
| Supprimer une fonction : | ` faas-cli remove <nom-fonction>` |

- **Go :**

| Action | Commande |
| :-------------: | :-------------: |
| Formate le code (pratique si ça veut pas build) : | `gofmt -s -w <nom-fonction>` |
| Teste la fonction | `go test <nom-fonction>`  |
| Pour ajouter une librairies externe : | `go get <nom-librairie>` |

## Exemples de fonctions :

[Des exemple peuvent être trouvés ici.](https://gitlab.insa-rennes.fr/mycelium-3.0/mycelium-3.0)

# Guide MQTT
**Ce guide permet de créer un serveur MQTT et de le faire fonctionner.**

## Création du serveur MQTT :

Se placer dans le dossier `faas-netes/chart` :
```sh
cd ~/faas-netes/chart
```

Regarder quelle est l'adresse IP de la VM (192. …) :
```sh
hostname -I  
```

Dans `mqtt-connectors/values.yaml`, on change l'adresse du broker avec l'adresse IP de la machine : `broker: tcp://<IP_MACHINE>:1883`

Vous pouvez aussi changer le nom du topic.

Appliquer les paramètres :
```sh
helm template -n openfaas --namespace openfaas mqtt-connector/ | kubectl apply -f -
```

Installer mosquitto :
```sh
sudo apt install mosquitto
systemctl start mosquitto.service
```

Voir si le serveur MQTT est déployé :
```sh
systemctl status mosquitto
```

Commande à faire pour que cela fonctionne (ça rajoute deux lignes importantes dans un fichier) :
```sh
echo -e "allow_anonymous true\nlistener 1883" | sudo tee -a /etc/mosquitto/mosquitto.conf
```

## Test du serveur MQTT :

Pour tester, on peut créer une fonction Open Faas simple, affichant simplement "Hello" par exemple (voir tutoriel OpenFaas). Il faut ensuite modifier le fichier .yml de la fonction en rajoutant le topic.

Exemple :
```
version: 1.0
provider:
  name: openfaas
  gateway: http://127.0.0.1:8080
functions:
  <name_function>:
    lang: golang-middleware
    handler: ./<path>
    image: <pseudo_docker>/<name-function>:latest
    annotations:
      topic: <name-topic>
```

Il faut ensuite déployer la fonction. On va publier sur le topic pour vérifier que la fonction se lance bien.

Installer des clients pour tester :
```sh
sudo apt-get install mosquitto-clients
```

Faire le test du serveur MQTT :
```sh
mosquitto_pub -h localhost -t <name-topic> -m "Hello World!"
```

On peut ensuite regarder les logs des fonctions et du mqtt connector pour voir que cela fonctionne (pour cela, on peut utiliser les commandes de la partie openfaas ou bien les commandes de la partie k9s).

Pour voir ce qu’on reçoit sur un topic :
```sh
mosquitto_sub -h localhost -t <name-topic>
```

## Méthode pour ajouter plusieurs topics :

Chaque topic nécessite la création d'un nouveau mqtt-connector. Pour cela, il faut créer un autre fichier values.yaml (et renommer le fichier). Ce fichier sera identique au fichier que l'on vient de créer, excepté le topic et le clientID qui doit être unique. Il faut se placer dans le dossier `faas-netes/chart`. 

On utilise ensuite la commande (dans faas/chart), le nom du déploiement peut être ce que l’on souhaite :
```sh
helm template -n openfaas --namespace openfaas --values <chemin vers le values.yaml> <nom du déploiement>  mqtt-connector/ | kubectl apply -f -
```

# Guide K9s

**Ce guide explique comment utiliser K9s, une interface utilisateur permettant de visualiser facilement les logs pour Open Faas.**

K9s sert à regarder les log des différents deployments et pods facilement depuis le terminal sans passer par d'autres commandes.

## Installations :
Téléchargement de Brew :
```sh
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

Ajout dans le path :
```sh
(echo; echo 'eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"') >> /home/user/.bashrc eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"
```

Installer K9S :
```sh
brew install k9s
```
## Commandes utiles :

| Action | Commande |
| :-------------: | :-------------: |
| Démarrer k9s pour openfaas : | `k9s -n openfaas` |
| Démarrer k9s pour voir les fonctions openfaas : | `k9s -n openfaas-fn`  |
| Plus globalement, la commande est : | `k9s -n <namespace>` |

# Guide Influx DB

**Ce guide permet de mettre en place une base de donnée avec Influx DB et de l’utiliser avec des fonctions Open Faas**

Télécharger Influx Data :
```sh
wget https://repos.influxdata.com/influxdata-archive_compat.key
echo '393e8779c89ac8d958f81f942f9ad7fb82a25e133faddaf92e15b16e6ac9ce4c influxdata-archive_compat.key' | sha256sum -c && cat influxdata-archive_compat.key | gpg --dearmor | sudo tee /etc/apt/trusted.gpg.d/influxdata-archive_compat.gpg > /dev/null
echo 'deb [signed-by=/etc/apt/trusted.gpg.d/influxdata-archive_compat.gpg] https://repos.influxdata.com/debian stable main' | sudo tee /etc/apt/sources.list.d/influxdata.list
```

Installer Influx DB :
```sh
sudo apt-get update
sudo apt-get install influxdb2
```

Démarrer Influx DB :
```sh
sudo service influxdb start
```

Vérifier si InfluxDB est bien en mode “active (running)” :
```sh
sushdo service influxdb status
```

Installer Influx CLI :
```sh
wget https://dl.influxdata.com/influxdb/releases/influxdb2-client-2.7.3-linux-amd64.tar.gz
tar xvzf ./influxdb2-client-2.7.3-linux-amd64.tar.gz
sudo mv ./influx /usr/local/bin/
```

Supprimer les fichiers inutiles téléchargés pendant l’installation :
```sh
rm influxdata-archive_compat.key && rm influxdb2-client-2.7.3-linux-amd64.tar.gz && rm LICENSE && rm README.md
```

Sur le PC, se connecter à localhost:8086, et remplir les différents champs. On obtient alors un token (⚠️Bien conserver le login et le mot de passe).

Créer la configuration InfluxDB :
```sh
influx config create --config-name <CONFIG_NAME_A_CHOISIR> \
  --host-url http://localhost:8086 \
  --org <ORGANISATION> \
  --token <TOKEN_API> \
  --active
```

Pour utiliser Influx dans des fonctions en Go, il faut installer le paquet dans le même répertoire que le fichier handler.go avec la commande :
```sh
go get github.com/influxdata/influxdb-client-go/v2
```

Si Go n’est pas installé :
```
sudo apt-get -y install golang-go
```

Une fonction utilisant cela peut être trouvé sur le git : https://gitlab.insa-rennes.fr/mycelium-3.0/mycelium-3.0

Dans la fonction, il faut modifier le token Influx DB.
À la ligne :   url :=   il faut rajouter l’adresse de la VM. Normalement, il n’y a pas besoin de la changer.

Pour trouver l’adresse, suivre cette manipulation :
```sh
sudo apt install net-tools
ifconfig
```
(l’adresse se trouve dans cni0, c’est l’adresse inet)

⚠️ Dans Influx DB, il faut créer le bucket utilisé dans la fonction, en utilisant le même nom.


# Guide ChirpStack

[Chirpstack Getting Started](https://www.chirpstack.io/docs/getting-started/debian-ubuntu.html)

Changer la ligne dans Setup software repository.

```sh
sudo apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys 1CE2AFD36DBCCA00
```

## Configurer Chirpstack gateway bridge :

```sh
sudo chmod go+rwx -R /etc/chirpstack-gateway-bridge/
nano /etc/chirpstack-gateway-bridge/chirpstack-gateway-bridge.toml
```

Trouver ces lignes et recopier eu868 qui manque : 

```
[integration.mqtt]
event_topic_template="eu868/gateway/{{ .GatewayID }}/event/{{ .EventType }}"
command_topic_template="eu868/gateway/{{ .GatewayID }}/command/#"
```

Redémarrer chirpstack-gateway-bridge :
```sh
sudo systemctl restart chirpstack-gateway-bridge
```

Modifier le port de l’API (sinon conflit avec openfaas car 8080 aussi, 8082 est parfait), ouvrir avec :
```sh
sudo nano /etc/chirpstack/chirpstack.toml
```
Puis mettre le contenu suivant : 
```
[api]

   # interface:port to bind the API interface to.
   bind=“0.0.0.0:8082”

   # Secret.
```

Pour accéder à l’interface dans un navigateur mettre localhost:8082 (ou le port choisi), les identifiants de base sont admin admin.

[ChirpStack open-source LoRaWAN® Network Server documentation](https://www.chirpstack.io/docs/getting-started/debian-ubuntu.html)

## Ajouter une gateway :

Si la gateway n’est pas encore installée et configurée allez au guide Gateway puis revenir après.

Dans Chirpstack → Gateways → Add gateway → Entrez un nom et le Gateway ID de votre gateway (pour notre gateway mycelium : 00800000a000368d) puis submit.

Si la gateway n’est pas marquée comme active, allez au guide Gateway et vérifiez que l’adresse du serveur est bien l’adresse du nœud où Chirpstack tourne. 

## Ajouter un capteur sur Chirpstack :

Dans device profil → Add device profil → General : Mettre un nom dans Name.

Puis dans l’onglet codec → Choisir comme Payload codec, JavaScript functions et remplacer le code par celui spécifique du capteur à ajouter.

Cliquer sur submit.

Ensuite dans Application → Add application : Mettre un nom puis submit.

Puis pour chaque capteur à ajouter dans l’application créée (il peut y en avoir plusieurs dans une même application).
Cliquer sur Add device, mettre un nom, entrer le Device EUI du capteur, choisir le device profil créé précédemment puis submit. 
Retourner dans les paramètre du capteurs (sur Chirpstack), aller dans OTAA keys et entrer l’app key (pour trouver la valeurs il faut regarder la notice constructeur, pour capteur Milesight il faut télécharger l’application, lire en nfc les données du capteurs, là vous cherchez et trouvez “Application Key”.
Soit vous en refaites un, soit de base c’est 5572404c696e6b4c6f52613230313823. 
Cette app key doit être la même pour le capteur et Chirpstack.

Site pour les codec Milesight Sensor : https://github.com/Milesight-IoT/SensorDecoders

# Guide Gateway
## Connexion Gateway-Cluster :

Si besoin, pour se connecter sur la gateway directement à partir d’un pc branché en ethernet sur la gateway (pour ubuntu) : 

Allez dans Paramètre → Réseau → Filaire → Cliquer pour paramétrer le réseau → IPV4 → Mettre en manuel → Dans adresse mettre par exemple 192.168.2.42 →dans masque 255.255.255.0 → Desactiver DNS et Routes → Appliquer

Maintenant accéder à l’interface de la gateway Multitech dans un navigateur mettre 192.168.2.1 (si ip mise dans réseau 192.168.2.42 sinon ?.?.?.1). Sur la page de connection mettre les identifiants de la Gateway (Pour mycélium c’est dans la boîte).

## Configuration de la gateway Multitech :

- **Si la Gateway a été réinitialisée, suivre ces étapes :**

(https://gitlab.insa-rennes.fr/suivi-environnemental/croix-verte/infrastructure/-/blob/master/gateway/README.md?ref_type=heads)

Firewall :
Add filter in input filter rules :
```
Name: Chirpstack input
Destination Port: 1700
Source interface: 1700
Protocol: UDP
```

Add filter in output filter rules :
```
Name: Chirpstack output
Destination Port: 1700
Source interface: 1700
Protocol: UDP
```

Save and restart.

Suivre la suite du guide.

- **Dans tous les cas :**

Changer l’adresse serveur dans LoRaWAN → Network Settings pour mettre celle où est installée Chirpstack (pour l’instant 192.168.2.11).

⇒ Pour moi, ça a marché uniquement avec 192.168.2.42, soit l’adresse que l’on a donné à l’arrivé de la connection ethernet → J’ai changé l’adresse par 192.168.2.44 et même constat, temps que je met pas cette IP ça marche pas et dès que je la met c’est bon.

## Connection Gateway-VM :

L’installation de Chirpstack reste la même, idem pour la configuration de la gateway. Cependant, pour la connexion à la gateway, il va falloir bricoler un peu plus.

Le problème est que la VM ne détecte pas le port Ethernet, elle obtient la connexion par le biais du PC hôte (Ethernet  → PC Hôte → VM). On va donc devoir faire en sorte que la VM détecte le port ethernet comme si c’était le sien.

Pour cela, sur Virt-Manager (Gestionnaire de machine virtuelle en français), on va cliquer sur la VM sans la lancer. Puis dans le menu en haut aller dans afficher > Détails. Ensuite ajouter un matériel > Périphérique Hôte PCI et choisissez le PCI qui correspond à votre port ethernet (faire lspci dans un terminal et chercher le ethernet → ex : 3b:00.0 Ethernet controller: Realtek Semiconductor Co., Ltd. Device 2502 (rev 1f), dans ce cas là c’est le port 3b:00.0 qu’il faut ajouter à notre VM).

Là, deux situations, soit le relais de périphérique PCI est activé sur votre PC, soit non. Si ce n’est pas le cas, votre VM ne se lancera pas. Pour régler ce problème :

```sh
sudo nano /etc/default/grub
```

Mettre dans GRUB_CMDLINE_LINUX intel_iommu=on si vous avez un processeur intel, ou amd_iommu=on si vous avez un AMD
GRUB_CMDLINE_LINUX="intel_iommu=on,igfx_off"
			ou
GRUB_CMDLINE_LINUX="amd_iommu=on,igfx=off" 
Pensez à vérifier que votre processeur est compatible avec VT-d (normalement ça touche que le processeur, mais certains bios semble pouvoir poser problème (ASUS en particulier).
```sh
sudo update-grub
```
Reboot la machine.

Normalement maintenant, la VM va pouvoir se lancer sans problème, faites lspci dans la VM et vérifier que maintenant il y a bien un ethernet.

## Configurer la connexion sur la VM :

À présent, il faut configurer la connexion, faites ifconfig et regarder les enp[0-9]*s0, normalement vous en avez un. Ce n’est pas celui qu’on veut, on va donc devoir ajouter. 
Faites :

```sh
ip addr
```

Chercher l’autre enp[0-9]*s0.
```sh
sudo ip link set enp7s0 up
```

(moi c’était 7, mettez le votre)
```sh
sudo ip addr add 192.168.2.44/24 dev enp7s0
```
	
Et voilà, maintenant la connexion venant de la gateway est configurée et ChirpStack va détecter votre gateway (sauf si vous avez foiré une étape).

La méthode ci-dessus doit être refait à chaque lancement de VM (normalement le enp reste le même).

# Guide Cluster

Pour pouvoir accéder au cluster : 
Chaque membre du groupe doit récupérer son adresse matériel (https://blog.shevarezo.fr/post/2019/01/08/comment-obtenir-adresse-ip-linux-ligne-de-commande).
Faire un ticket à la DSI de l’INSA pour avoir accès au réseau insaIOT, il faut dire que vous en avez besoin pour le projet Mycélium et préciser vos adresses MAC.
Se connecter à insaIOT (le mdp est unique et connu par Nikos) (Login : nikos, Mdp : nikos)

Pour se mettre en azerty :
```sh
sudo loadkeys fr
```

Fichier fait pour connexion wifi :
```sh
network:
  version: 2
  wifis:
    toto:
    access-points:
      mycelium
        password:
   
    dhcp4: true
```

## Scan adresses cluster :

Pour récupérer les adresses des autres clusters :
```sh
nmap -sn 192.168.2.0/24
```

Normalement on devrait récupérer les adresses:
- 192.168.2.10 
- 192.168.2.11 
- 192.168.2.12 
- 192.168.2.13 
- 192.168.2.14

## Manips tailscale (possiblement pas nécessaires) :

Tailscale : Service qui permet de créer un réseau sécurisé entre vos serveurs, VPN pour contourner insaIOT.

Installer tailscale sur machine perso et nikos0 :
```sh
sudo apt-get install tailscale
```

Si jamais cette commande ne fonctionne pas :
```sh
curl -fsSL https://tailscale.com/install.sh | sh
```

- **Sur nikos :**
```sh
sudo tailscale up
```

- **Sur machine perso :**
```sh
ssh nikos0
```

-**Utiliser l’agent :**

Si pas encore fait :
```sh
run ssh-agent
```

Faire :
```sh
eval "$(ssh-agent)"
```

Ajoutez la clé que vous souhaitez transférer à l'agent ssh :
```sh
ssh-add [chemin clé]/[nom clé].pem
```

Log sur nikos@nikos0 :
```sh
ssh -A [user]@[hostname]
```

Enregistrer les clés sous le nom id_rsa et id_rsa.pub.

Donner les autorisations pour les clés :
```sh
chmod 644 .id_rsa.pub && chmod 600 .id_rsa
```

Se connecter à la raspberry “principale” : 
```sh
ssh nikos@nikos0
```

Naviguer entre les machines avec ssh nikos@nikosX, se déconnecter après chaque ssh avec ctrl + D.

Installer les librairies précisées dans la rubrique Chirpstack si besoin pour chaque raspberry.

# Guide VPS

Aller sur [le vpn de l'INSA Rennes](https://vpn.insa-rennes.fr/gate/cloud/).

Installer GateClient s’il n’est pas installé. La connexion se fait avec les identifiants INSA.

pull le server RSS si il n’est pas présent : docker pull thomasderrien/rss-app
lancer le server RSS : docker run -p 80:80 rss-app 

# Guide Git avec SSH

Ce guide permet de mettre en place Git sur la VM avec une clé SSH pour pouvoir utiliser Git (c’est pas vraiment spécifique au projet mais ça peut être pratique)

```sh
ssh-keygen -t ed25519 -C "<adresse mail>"
sudo nano /home/user/.ssh/id_ed25519.pub
```

Copier la clé et la coller dans Git (section Ajouter clé SSH).

```sh
git config --global user.name "<pseudo Git>"
git config --global user.email <adresse mail>
```

# Commandes à relancer si besoin à chaque lancement de la VM

- **Se connecter en ssh :**
```sh
ssh -p 22 user@<Adresse IP VM>
```

On peut aussi se connecter tout en redirigeant les ports pour pouvoir ouvrir les interfaces web directement depuis le pc hôte :
```sh
ssh -p 22 -L 8082:localhost:8082 -L 8081:localhost:8080 -L 8086:localhost:8086  user@<Adresse IP VM>
```

Ces redirections sont effectués respectivement pour:
Chirpstack (8082),
OpenFaas (8081),
InfluxDB (8086)

- **Redirection pour Open Faas :**
```sh
kubectl port-forward -n openfaas svc/gateway 8080:8080 &
```

- **Port Ethernet pour Chirpstack :** 
```sh
sudo ip link set enp7s0 up  && sudo ip addr add 192.168.2.44/24 dev enp7s0 
```
(moi c’était 7, mettez le votre (voir Configurer la connexion sur la VM))