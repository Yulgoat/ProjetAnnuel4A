package function

import (
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"io"
	"log"
	"net/http"
	"time"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	var input []byte

	if r.Body != nil {
		defer r.Body.Close()

		body, _ := io.ReadAll(r.Body)

		input = body
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://192.168.122.21:1883") // Remplacer <IP_MACHINE>
	opts.SetClientID("fonction-mqtt")           // Remplacer <CLIENT_ID>

	// Création du client MQTT
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	// Publication du message sur le topic
	topic := "topic1" // Remplacer <NOM_TOPIC>
	message := string(input)

	token := client.Publish(topic, 0, false, message)
	token.Wait()

	fmt.Printf("Message publié sur le topic %s: %s\n", topic, message)

	// Attente avant de se déconnecter (facultatif)
	time.Sleep(2 * time.Second)

	// Déconnexion du client MQTT client.Disconnect (250)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Message MQTT Envoyé")))
}
