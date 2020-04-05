package main

import (
	"github.com/VolantMQ/vlapi/vlauth"
	"github.com/go-bongo/bongo"
	"gopkg.in/mgo.v2/bson"
)

type authProvider struct {
	cfg        config
	connection *bongo.Connection
}

type UserModel struct {
	bongo.DocumentBase `bson:",inline"`
	Username           string `json:"username" bson:"username"`
	Password           string `json:"password" bson:"password"`
	SubscriptionList   bool   `json:"subscription_list" bson:"subscription_list"`
	PublishList        bool   `json:"publish_list" bson:"publish_list"`
}

func (p *authProvider) Finduser(username string, password string) (user UserModel, error error) {

	err := p.connection.Collection(p.cfg.CollectionName).FindOne(bson.M{p.cfg.UsernameField: username, p.cfg.PasswordField: password}, &user)
	if err != nil {
		error = err
		return user, error
	} else {
		return user, nil
	}

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

func (p *authProvider) ACL(clientID, user, topic string, access vlauth.AccessType) error {
	return vlauth.StatusAllow
}

func (p *authProvider) Shutdown() error {
	return nil
}
