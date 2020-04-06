package main

import (
	"github.com/VolantMQ/vlapi/vlauth"
	"github.com/go-bongo/bongo"
	"gopkg.in/mgo.v2/bson" // TODO: Remove this, auto imported in vscode
)

type authProvider struct {
	cfg        config
	connection *bongo.Connection
}

type UserModel struct {
	bongo.DocumentBase `bson:",inline"`
	Username           string   `json:"username" bson:"username"`
	Password           string   `json:"password" bson:"password"`
	SubscriptionList   []string `json:"subscribe" bson:"subscribe"`
	PublishList        []string `json:"publish" bson:"publish"`
}

func (p *authProvider) Finduser(username string, password string) (user UserModel, error error) {
	err := p.connection.Collection(p.cfg.CollectionName).FindOne(bson.M{username: username, password: password}, &user)
	if err != nil {
		error = err
		return user, error
	}
	return user, nil
}

// We cannot use function FindUser with blank password because that could lead to unathorized access by sending blank password.
func (p *authProvider) FindUserByUsername(username string) (user UserModel, err error) {
	err = p.connection.Collection(p.cfg.CollectionName).FindOne(bson.M{username: username}, &user)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (p *authProvider) Connect() error {
	config := &bongo.Config{
		ConnectionString: p.cfg.MongodbBaseURI,
		Database:         p.cfg.DatabaseName,
	}
	connection, err := bongo.Connect(config)
	p.connection = connection
	if err != nil {
		return err
	}
	return nil
}

func (p *authProvider) Init() error {
	err := p.Connect()
	return err
}

func (p *authProvider) Password(clientID, username, password string) error {
	_, err := p.Finduser(username, password)
	if err != nil {
		return vlauth.StatusDeny
	}
	return vlauth.StatusAllow
}

func (p *authProvider) ACL(clientID, username, topic string, access vlauth.AccessType) error {
	// fmt.Printf("Allowing permission for %s, %s, %s %d \n", clientID, user, topic, access)
	user, err := p.FindUserByUsername(username)
	if err != nil {
		return vlauth.StatusDeny
	}
	permission := access.Type()
	allowed := false
	if permission == "write" {
		allowed = IsTopicAllowed(topic, user.PublishList)
	} else {
		allowed = IsTopicAllowed(topic, user.SubscriptionList)
	}
	if err == nil && allowed {
		return vlauth.StatusAllow
	}
	return vlauth.StatusDeny
}

func (p *authProvider) Shutdown() error {
	return nil
}
