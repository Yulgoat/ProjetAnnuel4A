Fonctions-Extraction :
	- Fonctions OpenFaas en Go 
	- Permet de récupérer les JSON décodés des capteurs (envoyé par Chirpstack sur MQTT). 
	- Récupère les données importantes, les stockent dans InfluxDB et va envoyer ces données au fonctions de traitement à lancer

Fonctions-Traitement : 
        - Fonctions OpenFaas en Go
	- Fonctions scénario
	- Envoie le résultat du scénario sur le topic MQTT notification

Fonctions-Standard :
	- Fonctions Go ou Python
	- Ce ne sont pas des fonctions Openfaas, elle se lance à la main
	- Sont là pour l'exemple et pour tester sans avoir à modifier sans cesse des fonctions openfaas	
