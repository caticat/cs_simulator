package pnet

type PMessageServer struct {
	*PMessage

	client *PServerClient
}

func newPMessageServer(client *PServerClient) *PMessageServer {
	return &PMessageServer{
		PMessage: newPMessage(),
		client:   client,
	}
}

func (this *PMessageServer) GetClient() *PServerClient {
	return this.client
}

func (this *PMessageServer) GetClientAddr() string {
	return this.client.getClientAddr()
}
