package main

import (
	"fmt"

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

func (p *authProvider) Finduser(username string, password string) (user *UserModel, error error) {
	result := &UserModel{}
	err := p.connection.Collection(p.cfg.CollectionName).FindOne(bson.M{p.cfg.UsernameField: user, p.cfg.PasswordField: password}, result)
	if err != nil {
		fmt.Println(err.Error())
		error = err
		return result, error
	} else {
		return result, nil
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
	user, err := p.Finduser(username, password)
	if err != nil && len(user.Username) > 0 {
		return vlauth.StatusAllow
	}
	return vlauth.StatusDeny
}

func (p *authProvider) ACL(clientID, user, topic string, access vlauth.AccessType) error {
	return vlauth.StatusDeny
}

func (p *authProvider) Shutdown() error {
	return nil
}
